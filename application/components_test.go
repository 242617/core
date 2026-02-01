package application

import (
	"context"
	"testing"

	"github.com/242617/core/mocks"
	"github.com/stretchr/testify/mock"
)

func TestComponentsByName(t *testing.T) {
	tests := []struct {
		name       string
		components Components
		searchName string
		want       string
		found      bool
	}{
		{
			name:       "find existing component",
			searchName: "db",
			want:       "db",
			found:      true,
		},
		{
			name:       "find second component",
			searchName: "cache",
			want:       "cache",
			found:      true,
		},
		{
			name:       "find last component",
			searchName: "cache",
			want:       "cache",
			found:      true,
		},
		{
			name:       "component not found",
			searchName: "nonexistent",
			found:      false,
		},
		{
			name:       "empty components",
			components: Components{},
			searchName: "any",
			found:      false,
		},
		{
			name:       "nil components",
			components: nil,
			searchName: "any",
			found:      false,
		},
		{
			name:       "empty search name",
			searchName: "",
			found:      false,
		},
		{
			name:       "duplicate names - first match",
			searchName: "db",
			want:       "db",
			found:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var comps Components
			if tt.components != nil {
				comps = tt.components
			} else {
				dbComp := mocks.NewComponent(t)
				dbComp.On("String").Maybe().Return("db")
				dbComp.On("Start", mock.Anything).Maybe().Return(nil)
				dbComp.On("Stop", mock.Anything).Maybe().Return(nil)

				cacheComp := mocks.NewComponent(t)
				cacheComp.On("String").Maybe().Return("cache")
				cacheComp.On("Start", mock.Anything).Maybe().Return(nil)
				cacheComp.On("Stop", mock.Anything).Maybe().Return(nil)

				serverComp := mocks.NewComponent(t)
				serverComp.On("String").Maybe().Return("server")
				serverComp.On("Start", mock.Anything).Maybe().Return(nil)
				serverComp.On("Stop", mock.Anything).Maybe().Return(nil)

				comps = Components{dbComp, cacheComp, serverComp}
			}

			got := comps.ByName(tt.searchName)
			if tt.found {
				if got == nil {
					t.Error("expected component, got nil")
					return
				}
				if got.String() != tt.want {
					t.Errorf("got %q, want %q", got.String(), tt.want)
				}
			} else {
				if got != nil {
					t.Errorf("expected nil, got %q", got.String())
				}
			}
		})
	}
}

func TestNewLifecycleComponent(t *testing.T) {
	tests := []struct {
		name          string
		componentName string
	}{
		{
			name:          "valid component",
			componentName: "test-component",
		},
		{
			name:          "empty name",
			componentName: "",
		},
		{
			name:          "name with spaces",
			componentName: "my component",
		},
		{
			name:          "name with special chars",
			componentName: "test-component_v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := mocks.NewLifecycle(t)
			lc := NewLifecycleComponent(tt.componentName, ml)
			if lc == nil {
				t.Fatal("expected LifecycleComponent, got nil")
			}
			if lc.String() != tt.componentName {
				t.Errorf("String() = %q, want %q", lc.String(), tt.componentName)
			}
			if lc.Lifecycle != ml {
				t.Error("Lifecycle field not set correctly")
			}
		})
	}
}

func TestLifecycleComponentString(t *testing.T) {
	tests := []struct {
		name string
		lc   *LifecycleComponent
		want string
	}{
		{
			name: "simple name",
			lc:   NewLifecycleComponent("test", mocks.NewLifecycle(t)),
			want: "test",
		},
		{
			name: "empty name",
			lc:   NewLifecycleComponent("", mocks.NewLifecycle(t)),
			want: "",
		},
		{
			name: "name with special characters",
			lc:   NewLifecycleComponent("my-service_v1.0", mocks.NewLifecycle(t)),
			want: "my-service_v1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lc.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLifecycleComponentLifecycle(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "start and stop",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := mocks.NewLifecycle(t)
			ml.On("Start", mock.Anything).Return(nil)
			ml.On("Stop", mock.Anything).Return(nil)

			lc := NewLifecycleComponent("test", ml)

			ctx := context.Background()

			if err := lc.Start(ctx); err != nil {
				t.Fatalf("Start() error = %v", err)
			}

			if err := lc.Stop(ctx); err != nil {
				t.Fatalf("Stop() error = %v", err)
			}
		})
	}
}
