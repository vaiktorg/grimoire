package grimoire

import (
	"fmt"

	"github.com/vaiktorg/grimoire/errs"
)

type State string
type Trigger string

type StateMachine struct {
	configs        map[State]*StateConfig
	onStateChanged func(prev, curr State)
	prevState      State
	currState      State
}

func NewStateMachine(initialState State) *StateMachine {
	return &StateMachine{
		configs:   make(map[State]*StateConfig),
		currState: initialState,
	}
}

func (sm *StateMachine) State(newState State) *StateConfig {
	if _, ok := sm.configs[newState]; !ok {
		sm.configs[newState] = &StateConfig{allowable: make(map[Trigger]State)}
	}
	return sm.configs[newState]
}

func (sm *StateMachine) OnStateChanged(onstatechanged func(prev, curr State)) {
	if onstatechanged == nil {
		sm.onStateChanged = func(prev, curr State) {
			fmt.Printf("State transitions from %s to %s\n", prev, curr)
		}
	}
	sm.onStateChanged = onstatechanged
}

func (sm *StateMachine) Fire(trigger Trigger) error {
	if currConf, ok := sm.configs[sm.currState]; ok {
		if nextState, ok := currConf.allowable[trigger]; ok {
			if err := currConf.onExit(sm.currState); err != nil {
				return err
			}

			if newConf, ok := sm.configs[nextState]; ok {
				var fallbacked bool

				if newConf.onEntry == nil {
					return errs.Error(fmt.Sprintf("No state behavior set for triggered state %s.\n", nextState))
				}

				if err := newConf.onEntry(nextState); err != nil {
					if newConf.onFail != nil {
						newConf.onFail(err)
					}

					if len(newConf.fallback) == 0 {
						return err
					}

					newConf, ok = sm.configs[newConf.fallback]
					if !ok {
						return errs.Error(fmt.Sprintf("No configured state %s\n", trigger))
					}

					fallbacked = true
				}

				sm.onStateChanged(sm.currState, nextState)
				sm.prevState = sm.currState
				sm.currState = nextState

				if fallbacked {
					return nil
				}

				if newConf.onSuccess != nil {
					newConf.onSuccess()
				}

				return nil
			}
			return errs.Error(fmt.Sprintf("No configured state %s\n", trigger))
		}
		return errs.Error(fmt.Sprintf("Trigger %s not allowed...\n", trigger))
	}
	return errs.Error(fmt.Sprintf("SM has no current state %v\n", sm.currState))
}

func (sm *StateMachine) GetCurrentState() State {
	return sm.currState
}
