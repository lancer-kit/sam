package examples

import "github.com/sheb-gregor/sam"

const (
	WStateDisabled    sam.State = "Disabled"
	WStateEnabled     sam.State = "Enabled"
	WStateInitialized sam.State = "Initialized"
	WStateRun         sam.State = "Run"
	WStateStopped     sam.State = "Stopped"
	WStateFailed      sam.State = "Failed"
)

// newWorkerSM returns filled state machine of worker lifecycle
//
// (*) -> [Disabled] -> [Enabled] -> [Initialized] -> [Run] <-> [Stopped]
//          ↑ ↑____________|  |          |  |  ↑         |
//          |_________________|__________|  |  |------|  ↓
//                            |-------------|-----> [Failed]

func NewWorkerSM() sam.StateMachine {
	workerSM := sam.NewStateMachine()
	_ = workerSM.AddTransitions(WStateDisabled, WStateEnabled)
	_ = workerSM.AddTransitions(WStateEnabled, WStateInitialized, WStateFailed, WStateDisabled)
	_ = workerSM.AddTransitions(WStateInitialized, WStateRun, WStateFailed, WStateDisabled)
	_ = workerSM.AddTransitions(WStateRun, WStateStopped, WStateFailed)
	_ = workerSM.AddTransitions(WStateStopped, WStateRun)
	_ = workerSM.AddTransitions(WStateFailed, WStateInitialized, WStateDisabled)
	workerSM.SetState(WStateDisabled)
	return workerSM
}
