# Sam

**Sam** is the finite state machine engine.

You can define state diagram by adding *Transitions* and add transition hooks.

## Usage

Add to project using `go get`

```shell

go get github.com/lancer-kit/sam

```

## Examples

- Without hooks

```go
package main

import (
	"fmt"

	"github.com/lancer-kit/sam"
)

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
func main() {
	workerSM, err := sam.NewStateMachine().
		AddTransitions(WStateDisabled, WStateEnabled).
		AddTransitions(WStateEnabled, WStateInitialized, WStateFailed, WStateDisabled).
		AddTransitions(WStateInitialized, WStateRun, WStateFailed, WStateDisabled).
		AddTransitions(WStateRun, WStateStopped, WStateFailed).
		AddTransitions(WStateStopped, WStateRun).
		AddTransitions(WStateFailed, WStateInitialized, WStateDisabled).
		Finalize(WStateDisabled)
	if err != nil || workerSM == nil {
		fmt.Println("init failed: ", err)
	}
	clone := workerSM.Clone()
	
	err = clone.GoTo(WStateEnabled)
	if err != nil {
		fmt.Println("error: ", err)
	}

}

```
