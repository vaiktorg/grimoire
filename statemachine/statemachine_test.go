package grimoire

import (
	"errors"
	"fmt"
	"testing"
)

func TestStateMachine(t *testing.T) {
	sm := NewStateMachine("Idle")
	sm.OnStateChanged(func(prev, curr State) {
		fmt.Printf("Changing states from: %s to: %s...\n", prev, curr)
	})
	configure(sm)

	sm.Fire("Start")
	sm.Fire("Update")
	sm.Fire("Stop")
	sm.Fire("Dispose")
}
func Entrying(state State) error {
	fmt.Println("Entrying State ", state)
	return errors.New(string(state))
}
func Exiting(state State) error {
	fmt.Println("Exiting State ", state)
	return nil
}
func configure(sm *StateMachine) {
	const (
		Idle      = State("Idle")
		Running   = State("Idle")
		Updating  = State("Idle")
		Disposing = State("Idle")

		Start   = Trigger("Start")
		Update  = Trigger("Update")
		Stop    = Trigger("Stop")
		Dispose = Trigger("Dispose")
	)

	sm.State(Idle).
		OnEntry(Entrying).
		Allow(Start, Running).
		Allow(Dispose, Disposing).
		OnExit(Exiting)
	sm.State(Running).
		OnEntry(Entrying).
		Allow(Stop, Idle).
		Allow(Update, Updating).
		OnExit(Exiting)
	sm.State(Updating).
		OnEntry(Entrying).
		Allow(Stop, Idle).
		Allow(Dispose, Disposing).
		OnExit(Exiting)
	sm.State(Disposing).
		OnEntry(Entrying)
}
