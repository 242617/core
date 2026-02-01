package application

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/242617/core/mocks"
)

func TestDefaults(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if app.name != "application" {
		t.Errorf("default name = %q, want 'application'", app.name)
	}

	if app.startTimeout != 30*time.Second {
		t.Errorf("default start timeout = %v, want 30s", app.startTimeout)
	}

	if app.stopTimeout != 30*time.Second {
		t.Errorf("default stop timeout = %v, want 30s", app.stopTimeout)
	}

	if app.hostname == "" {
		t.Error("default hostname should not be empty")
	}
}

func TestWithName(t *testing.T) {
	tests := []struct {
		name      string
		appName   string
		want      string
		wantEmpty bool
	}{
		{
			name:    "valid name",
			appName: "my-service",
			want:    "my-service",
		},
		{
			name:    "name with spaces",
			appName: "my service",
			want:    "my service",
		},
		{
			name:    "name with special chars",
			appName: "service_v1.0",
			want:    "service_v1.0",
		},
		{
			name:    "default name",
			appName: "",
			want:    "application",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts []Option
			if tt.appName != "" {
				opts = []Option{WithName(tt.appName)}
			}
			app, err := New(opts...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if app.name != tt.want {
				t.Errorf("name = %q, want %q", app.name, tt.want)
			}
		})
	}
}

func TestWithStartTimeout(t *testing.T) {
	tests := []struct {
		name        string
		timeout     time.Duration
		wantErr     bool
		errContains string
		wantDefault bool
	}{
		{
			name:    "valid timeout",
			timeout: 10 * time.Second,
			wantErr: false,
		},
		{
			name:    "minimum timeout",
			timeout: 1 * time.Nanosecond,
			wantErr: false,
		},
		{
			name:        "negative timeout",
			timeout:     -1 * time.Second,
			wantErr:     true,
			errContains: "invalid start timeout",
		},
		{
			name:        "zero timeout",
			timeout:     0,
			wantErr:     true,
			errContains: "invalid start timeout",
		},
		{
			name:        "zero with other options",
			timeout:     0,
			wantErr:     true,
			errContains: "invalid start timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := New(WithStartTimeout(tt.timeout))
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !containsSubstring(err.Error(), tt.errContains) {
					t.Errorf("error should contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if app.startTimeout != tt.timeout {
				t.Errorf("start timeout = %v, want %v", app.startTimeout, tt.timeout)
			}
		})
	}
}

func TestWithStopTimeout(t *testing.T) {
	tests := []struct {
		name        string
		timeout     time.Duration
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid timeout",
			timeout: 5 * time.Second,
			wantErr: false,
		},
		{
			name:    "long timeout",
			timeout: 1 * time.Hour,
			wantErr: false,
		},
		{
			name:        "negative timeout",
			timeout:     -1 * time.Second,
			wantErr:     true,
			errContains: "invalid stop timeout",
		},
		{
			name:        "zero timeout",
			timeout:     0,
			wantErr:     true,
			errContains: "invalid stop timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := New(WithStopTimeout(tt.timeout))
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !containsSubstring(err.Error(), tt.errContains) {
					t.Errorf("error should contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if app.stopTimeout != tt.timeout {
				t.Errorf("stop timeout = %v, want %v", app.stopTimeout, tt.timeout)
			}
		})
	}
}

func TestWithLogger(t *testing.T) {
	t.Run("valid logger", func(t *testing.T) {
		logger := mocks.NewLogger(t)
		logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
		logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
		logger.On("Warn", mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
		logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
		app, err := New(WithLogger(logger))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if app.log == nil {
			t.Error("logger should not be nil")
		}
	})

	t.Run("nil logger", func(t *testing.T) {
		_, err := New(WithLogger(nil))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestWithComponents(t *testing.T) {
	tests := []struct {
		name       string
		components []Component
		wantErr    bool
		wantCount  int
	}{
		{
			name: "single component",
			components: func() []Component {
				comp := mocks.NewComponent(t)
				comp.On("String").Maybe().Return("db")
				comp.On("Start", mock.Anything).Maybe().Return(nil)
				comp.On("Stop", mock.Anything).Maybe().Return(nil)
				return []Component{comp}
			}(),
			wantErr:   false,
			wantCount: 1,
		},
		{
			name: "multiple components",
			components: func() []Component {
				db := mocks.NewComponent(t)
				db.On("String").Maybe().Return("db")
				db.On("Start", mock.Anything).Maybe().Return(nil)
				db.On("Stop", mock.Anything).Maybe().Return(nil)

				cache := mocks.NewComponent(t)
				cache.On("String").Maybe().Return("cache")
				cache.On("Start", mock.Anything).Maybe().Return(nil)
				cache.On("Stop", mock.Anything).Maybe().Return(nil)

				server := mocks.NewComponent(t)
				server.On("String").Maybe().Return("server")
				server.On("Start", mock.Anything).Maybe().Return(nil)
				server.On("Stop", mock.Anything).Maybe().Return(nil)

				return []Component{db, cache, server}
			}(),
			wantErr:   false,
			wantCount: 3,
		},
		{
			name:       "nil components",
			components: nil,
			wantErr:    false,
			wantCount:  0,
		},
		{
			name:       "empty components",
			components: []Component{},
			wantErr:    false,
			wantCount:  0,
		},
		{
			name: "nil component in slice",
			components: func() []Component {
				comp := mocks.NewComponent(t)
				comp.On("String").Maybe().Return("db")
				comp.On("Start", mock.Anything).Maybe().Return(nil)
				comp.On("Stop", mock.Anything).Maybe().Return(nil)
				return []Component{comp, nil}
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := New(WithComponents(tt.components...))
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(app.components) != tt.wantCount {
				t.Errorf("components count = %d, want %d", len(app.components), tt.wantCount)
			}
		})
	}
}

func TestWithHostname(t *testing.T) {
	tests := []struct {
		name        string
		hostname    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid hostname",
			hostname: "localhost",
			wantErr:  false,
		},
		{
			name:     "hostname with domain",
			hostname: "my-host.example.com",
			wantErr:  false,
		},
		{
			name:     "IP address",
			hostname: "192.168.1.1",
			wantErr:  false,
		},
		{
			name:        "empty hostname",
			hostname:    "",
			wantErr:     true,
			errContains: "empty hostname",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := New(WithHostname(tt.hostname))
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !containsSubstring(err.Error(), tt.errContains) {
					t.Errorf("error should contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if app.hostname != tt.hostname {
				t.Errorf("hostname = %q, want %q", app.hostname, tt.hostname)
			}
		})
	}
}

func TestMultipleOptions(t *testing.T) {
	db := mocks.NewComponent(t)
	db.On("String").Maybe().Return("db")
	db.On("Start", mock.Anything).Maybe().Return(nil)
	db.On("Stop", mock.Anything).Maybe().Return(nil)

	logger := mocks.NewLogger(t)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	app, err := New(
		WithName("test-app"),
		WithStartTimeout(10*time.Second),
		WithStopTimeout(5*time.Second),
		WithHostname("test-host"),
		WithLogger(logger),
		WithComponents(db),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if app.name != "test-app" {
		t.Errorf("name = %q, want 'test-app'", app.name)
	}

	if app.startTimeout != 10*time.Second {
		t.Errorf("start timeout = %v, want 10s", app.startTimeout)
	}

	if app.stopTimeout != 5*time.Second {
		t.Errorf("stop timeout = %v, want 5s", app.stopTimeout)
	}

	if app.hostname != "test-host" {
		t.Errorf("hostname = %q, want 'test-host'", app.hostname)
	}

	if len(app.components) != 1 {
		t.Errorf("components count = %d, want 1", len(app.components))
	}
}

func TestOptionOverride(t *testing.T) {
	app, err := New(
		WithName("custom-name"),
		WithName("overridden-name"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if app.name != "overridden-name" {
		t.Errorf("name = %q, want 'overridden-name'", app.name)
	}
}
