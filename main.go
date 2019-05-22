package sam

import (
	"fmt"
)

type State string

type stateObj struct {
	name  State
	prev  State
	to    map[State]struct{}
	from  map[State]struct{}
	exist bool
}

var invalidTransition = func(from, to State) error {
	return fmt.Errorf("invalid transition: %v --> %v", from, to)
}

var stateNotFound = func(name State) error {
	return fmt.Errorf("state not found: %v", name)
}

type Hook func(from, to State) error

type HookList struct {
	before      []Hook
	after       []Hook
	beforeState map[State]Hook
	afterState  map[State]Hook
}

func (HookList) New() HookList {
	return HookList{
		before:      []Hook{},
		after:       []Hook{},
		beforeState: map[State]Hook{},
		afterState:  map[State]Hook{},
	}
}

func (hl *HookList) Execute(from, to State) (err error) {
	for i, hook := range hl.after {
		if err = hook(from, to); err != nil {
			return fmt.Errorf("after hook #%d failed; err: %v", i, err)
		}
	}

	hook, ok := hl.afterState[from]
	if ok {
		if err = hook(from, to); err != nil {
			return fmt.Errorf("after hook for [%s] failed; err: %v", from, err)
		}
	}

	for i, hook := range hl.before {
		if err = hook(from, to); err != nil {
			return fmt.Errorf("before hook #%d failed; err: %v", i, err)
		}
	}

	hook, ok = hl.beforeState[to]
	if ok {
		if err = hook(from, to); err != nil {
			return fmt.Errorf("before hook for [%s] failed; err: %v", to, err)
		}
	}
	return
}

type StateMachine struct {
	current State
	error   error
	states  map[State]stateObj
	hooks   HookList
}

func NewStateMachine() StateMachine {
	return StateMachine{}.New()
}

func (StateMachine) New() StateMachine {
	return StateMachine{
		current: "",
		states:  map[State]stateObj{},
		hooks:   HookList.New(HookList{}),
	}
}

func (sm *StateMachine) Clone() StateMachine {
	return *sm
}

func (sm *StateMachine) State() State {
	return sm.current
}

func (sm *StateMachine) SetState(state State) error {
	if sm.error != nil {
		return sm.error
	}

	stateObj := sm.getState(state)
	stateObj.prev = sm.current

	sm.setState(stateObj)
	sm.current = stateObj.name

	return sm.error
}

func (sm *StateMachine) Error() error {
	return sm.error
}

func (sm *StateMachine) Finalize(state State) (*StateMachine, error) {
	err := sm.SetState(state)
	return sm, err
}

func (sm *StateMachine) getState(name State) stateObj {
	state, ok := sm.states[name]
	if !ok {
		state = stateObj{
			name:  name,
			exist: false,
			to:    map[State]struct{}{},
			from:  map[State]struct{}{},
		}
	}
	return state
}

func (sm *StateMachine) setState(state stateObj) {
	state.exist = true
	sm.states[state.name] = state
}

func (sm *StateMachine) AddAfterAllHook(hook Hook) *StateMachine {
	sm.hooks.after = append(sm.hooks.after, hook)
	return sm
}

func (sm *StateMachine) SetAfterHook(state State, hook Hook) *StateMachine {
	sm.hooks.afterState[state] = hook
	return sm
}

func (sm *StateMachine) AddBeforeAllHook(hook Hook) *StateMachine {
	sm.hooks.before = append(sm.hooks.before, hook)
	return sm
}

func (sm *StateMachine) SetBeforeHook(state State, hook Hook) *StateMachine {
	sm.hooks.beforeState[state] = hook
	return sm
}

func (sm *StateMachine) RegisterState(state State, before, after Hook) *StateMachine {
	stateObj := sm.getState(state)
	sm.setState(stateObj)

	sm.hooks.beforeState[state] = before
	sm.hooks.afterState[state] = after

	return sm
}

func (sm *StateMachine) AddTransitions(from State, to ...State) *StateMachine {
	for _, name := range to {
		sm.AddTransition(from, name)
		if sm.Error() != nil {
			return sm
		}
	}
	return sm
}

func (sm *StateMachine) AddTransition(from, to State) *StateMachine {
	if sm.error != nil {
		return sm
	}

	if from == to {
		sm.error = invalidTransition(from, to)
		return sm
	}

	fromState := sm.getState(from)
	fromState.to[to] = struct{}{}
	sm.setState(fromState)

	toState := sm.getState(to)
	toState.from[from] = struct{}{}
	sm.setState(toState)

	return sm
}

func (sm *StateMachine) GoTo(toState State) error {
	newState, ok := sm.states[toState]
	if !ok {
		return stateNotFound(toState)
	}
	if sm.current == toState {
		return nil
	}

	state := sm.states[sm.current]
	_, ok = state.to[toState]
	if !ok {
		return invalidTransition(sm.current, toState)
	}

	err := sm.hooks.Execute(sm.current, toState)
	if err != nil {
		return err
	}

	newState.prev = sm.current
	sm.states[toState] = newState
	sm.current = toState

	return nil
}

func (sm *StateMachine) GoBack() error {
	current := sm.getState(sm.current)
	if !current.exist {
		return stateNotFound(sm.current)
	}

	_, ok := current.from[current.prev]
	if !ok {
		return invalidTransition(sm.current, current.prev)
	}

	prev := sm.getState(current.prev)
	if !prev.exist {
		return stateNotFound(current.prev)
	}

	sm.current = prev.name
	return nil
}
