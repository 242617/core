package application

import "fmt"

var ErrApplicationAlreadyStarted = fmt.Errorf("application already started")

// ComponentError for start/stop failures.
type ComponentError struct {
	Component string
	Phase     string // ComponentPhaseStart or ComponentPhaseStop
	Err       error
}

func (e *ComponentError) Error() string {
	return fmt.Sprintf("%s component %q: %v", e.Phase, e.Component, e.Err)
}
func (e *ComponentError) Unwrap() error { return e.Err }

// ErrInvalidSignals returned when invalid signals are provided.
type ErrInvalidSignals string

func (e ErrInvalidSignals) Error() string { return string(e) }
