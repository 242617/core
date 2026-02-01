package application

import (
	"errors"
	"fmt"
	"testing"
)

func TestComponentError_Error(t *testing.T) {
	tests := []struct {
		name      string
		component string
		phase     string
		err       error
		want      string
	}{
		{
			name:      "start error",
			component: "database",
			phase:     ComponentPhaseStart,
			err:       errors.New("connection failed"),
			want:      "start component \"database\": connection failed",
		},
		{
			name:      "stop error",
			component: "cache",
			phase:     ComponentPhaseStop,
			err:       errors.New("timeout"),
			want:      "stop component \"cache\": timeout",
		},
		{
			name:      "wrapped error",
			component: "server",
			phase:     ComponentPhaseStart,
			err:       fmt.Errorf("bind failed: %w", errors.New("address in use")),
			want:      "start component \"server\": bind failed: address in use",
		},
		{
			name:      "nil error",
			component: "component",
			phase:     ComponentPhaseStart,
			err:       nil,
			want:      "start component \"component\": <nil>",
		},
		{
			name:      "empty component name",
			component: "",
			phase:     ComponentPhaseStart,
			err:       errors.New("error"),
			want:      "start component \"\": error",
		},
		{
			name:      "empty phase",
			component: "comp",
			phase:     "",
			err:       errors.New("error"),
			want:      " component \"comp\": error",
		},
		{
			name:      "special characters in name",
			component: "my-service_v1.0",
			phase:     ComponentPhaseStop,
			err:       errors.New("error"),
			want:      "stop component \"my-service_v1.0\": error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce := &ComponentError{
				Component: tt.component,
				Phase:     tt.phase,
				Err:       tt.err,
			}
			if got := ce.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestComponentError_Unwrap(t *testing.T) {
	tests := []struct {
		name string
		ce   *ComponentError
		want error
	}{
		{
			name: "unwrapped error",
			ce: &ComponentError{
				Component: "db",
				Phase:     ComponentPhaseStart,
				Err:       errors.New("failed"),
			},
			want: errors.New("failed"),
		},
		{
			name: "wrapped error",
			ce: &ComponentError{
				Component: "cache",
				Phase:     ComponentPhaseStart,
				Err:       fmt.Errorf("failed: %w", errors.New("connection")),
			},
			want: fmt.Errorf("failed: %w", errors.New("connection")),
		},
		{
			name: "nil error",
			ce: &ComponentError{
				Component: "comp",
				Phase:     ComponentPhaseStart,
				Err:       nil,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ce.Unwrap()
			if tt.want == nil {
				if got != nil {
					t.Errorf("Unwrap() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Error("Unwrap() = nil, want error")
				return
			}
			if got.Error() != tt.want.Error() {
				t.Errorf("Unwrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComponentError_ErrorsIs(t *testing.T) {
	baseErr := errors.New("base error")
	ce := &ComponentError{
		Component: "comp",
		Phase:     ComponentPhaseStart,
		Err:       fmt.Errorf("wrapped: %w", baseErr),
	}

	if !errors.Is(ce, baseErr) {
		t.Error("errors.Is should find base error")
	}
}

func TestErrApplicationAlreadyStarted(t *testing.T) {
	if ErrApplicationAlreadyStarted == nil {
		t.Fatal("ErrApplicationAlreadyStarted should not be nil")
	}
	if ErrApplicationAlreadyStarted.Error() != "application already started" {
		t.Errorf("unexpected error message: %s", ErrApplicationAlreadyStarted.Error())
	}
}

func TestErrInvalidSignals(t *testing.T) {
	tests := []struct {
		name string
		err  ErrInvalidSignals
		want string
	}{
		{
			name: "simple error",
			err:  "invalid signal",
			want: "invalid signal",
		},
		{
			name: "error with details",
			err:  "invalid signal SIGKILL",
			want: "invalid signal SIGKILL",
		},
		{
			name: "empty error",
			err:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestComponentPhaseConstants(t *testing.T) {
	tests := []struct {
		name  string
		phase string
	}{
		{
			name:  "start phase",
			phase: ComponentPhaseStart,
		},
		{
			name:  "stop phase",
			phase: ComponentPhaseStop,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.phase == "" {
				t.Error("phase should not be empty")
			}
		})
	}
}
