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
	workerSM, err := sam.NewStateMachine().
		AddTransitions(WStateDisabled, WStateEnabled).
		AddTransitions(WStateEnabled, WStateInitialized, WStateFailed, WStateDisabled).
		AddTransitions(WStateInitialized, WStateRun, WStateFailed, WStateDisabled).
		AddTransitions(WStateRun, WStateStopped, WStateFailed).
		AddTransitions(WStateStopped, WStateRun).
		AddTransitions(WStateFailed, WStateInitialized, WStateDisabled).
		Finalize(WStateDisabled)
	if err != nil || workerSM == nil {
		panic("Init failed:" + err.Error())
	}
	return workerSM.Clone()
}
