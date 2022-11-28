package application

import (
	"github.com/looplab/fsm"
	"github.com/pkg/errors"
)

var ErrUnexpectedState = errors.New("unexpected transition")

const (
	ActionStart = "start"
	ActionStop  = "stop"

	StateStopped = "stopped"
	StateStarted = "started"
)

func NewFSM() *FSM {
	return &FSM{
		fsm: fsm.NewFSM(
			StateStopped,
			fsm.Events{
				{Name: ActionStart, Src: []string{StateStopped}, Dst: StateStarted},
				{Name: ActionStop, Src: []string{StateStarted}, Dst: StateStopped},
			},
			fsm.Callbacks{},
		),
	}
}

type FSM struct {
	fsm *fsm.FSM
}

func (fsm *FSM) Start() error             { return fsm.fsm.Event(ActionStart) }
func (fsm *FSM) Stop() error              { return fsm.fsm.Event(ActionStop) }
func (fsm *FSM) Event(event string) error { return fsm.fsm.Event(event) }

func (fsm *FSM) CanStart() bool    { return fsm.fsm.Can(ActionStart) }
func (fsm *FSM) CannotStart() bool { return !fsm.CanStart() }
func (fsm *FSM) CanStop() bool     { return fsm.fsm.Can(ActionStop) }
func (fsm *FSM) CannotStop() bool  { return !fsm.CanStop() }
