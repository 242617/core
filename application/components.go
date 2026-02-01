package application

import (
	"fmt"

	"github.com/242617/core/protocol"
)

const (
	ComponentPhaseStart = "start"
	ComponentPhaseStop  = "stop"
)

// Component with Start/Stop lifecycle.
type Component interface {
	fmt.Stringer
	protocol.Lifecycle
}

type Components []Component

// ByName finds component by name.
func (cmps *Components) ByName(name string) Component {
	for _, cmp := range *cmps {
		if cmp.String() == name {
			return cmp
		}
	}
	return nil
}

func NewLifecycleComponent(name string, cmp protocol.Lifecycle) *LifecycleComponent {
	return &LifecycleComponent{name, cmp}
}

type LifecycleComponent struct {
	string
	protocol.Lifecycle
}

func (s *LifecycleComponent) String() string { return s.string }
