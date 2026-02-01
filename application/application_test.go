package application

import (
	"errors"
	"testing"
	"time"

	"github.com/242617/core/mocks"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		options []Option
		wantErr bool
		errMsg  string
	}{
		{
			name:    "defaults only",
			options: nil,
			wantErr: false,
		},
		{
			name: "with custom options",
			options: []Option{
				WithName("test-app"),
				WithStartTimeout(10 * time.Second),
				WithStopTimeout(5 * time.Second),
			},
			wantErr: false,
		},
		{
			name: "nil component",
			options: []Option{
				WithComponents(nil),
			},
			wantErr: true,
			errMsg:  "component at index 0 is nil",
		},
		{
			name: "nil logger",
			options: []Option{
				WithLogger(nil),
			},
			wantErr: true,
			errMsg:  "empty logger",
		},
		{
			name: "invalid start timeout",
			options: []Option{
				WithStartTimeout(-1),
			},
			wantErr: true,
			errMsg:  "invalid start timeout",
		},
		{
			name: "zero start timeout",
			options: []Option{
				WithStartTimeout(0),
			},
			wantErr: true,
			errMsg:  "invalid start timeout",
		},
		{
			name: "invalid stop timeout",
			options: []Option{
				WithStopTimeout(-1),
			},
			wantErr: true,
			errMsg:  "invalid stop timeout",
		},
		{
			name: "zero stop timeout",
			options: []Option{
				WithStopTimeout(0),
			},
			wantErr: true,
			errMsg:  "invalid stop timeout",
		},
		{
			name: "empty hostname",
			options: []Option{
				WithHostname(""),
			},
			wantErr: true,
			errMsg:  "empty hostname",
		},
		{
			name: "empty name",
			options: []Option{
				WithName(""),
			},
			wantErr: true,
			errMsg:  "empty name",
		},
		{
			name: "with valid components",
			options: func() []Option {
				mockComp1 := mocks.NewComponent(t)
				mockComp1.On("String").Maybe().Return("comp1")
				mockComp1.On("Start", mock.Anything).Maybe().Return(nil)
				mockComp1.On("Stop", mock.Anything).Maybe().Return(nil)

				mockComp2 := mocks.NewComponent(t)
				mockComp2.On("String").Maybe().Return("comp2")
				mockComp2.On("Start", mock.Anything).Maybe().Return(nil)
				mockComp2.On("Stop", mock.Anything).Maybe().Return(nil)

				return []Option{WithComponents(mockComp1, mockComp2)}
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := New(tt.options...)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errMsg != "" && !errors.Is(err, nil) && !containsSubstring(err.Error(), tt.errMsg) {
					t.Errorf("error message should contain %q, got %q", tt.errMsg, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if app == nil {
				t.Fatal("expected app, got nil")
			}
			if app.stopCh == nil {
				t.Error("expected stopCh to be initialized")
			}
			if app.cancel == nil {
				t.Error("expected cancel to be initialized")
			}
		})
	}
}

func TestExit(t *testing.T) {
	tests := []struct {
		name        string
		exitCount   int
		expectPanic bool
	}{
		{
			name:      "single exit",
			exitCount: 1,
		},
		{
			name:      "multiple exits",
			exitCount: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLog := mocks.NewLogger(t)
			mockLog.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
			mockLog.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
			mockLog.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
			mockLog.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
			app, err := New(WithLogger(mockLog))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for i := 0; i < tt.exitCount; i++ {
				app.Exit()
			}
		})
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstringHelper(s, substr))
}

func containsSubstringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
