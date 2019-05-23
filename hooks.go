package sam

import "fmt"

type Hook func(from, to State) error

type HookList struct {
	before      []Hook
	after       []Hook
	rollback    map[State]Hook
	beforeState map[State]Hook
	afterState  map[State]Hook
}

func (HookList) New() HookList {
	return HookList{
		before:      []Hook{},
		after:       []Hook{},
		beforeState: map[State]Hook{},
		afterState:  map[State]Hook{},
		rollback:    map[State]Hook{},
	}
}

func (hl *HookList) ExecuteRollback(from, to State) (err error) {
	hook, ok := hl.rollback[from]
	if ok {
		if err = hook(from, to); err != nil {
			return fmt.Errorf("rollaback hook for [%s] failed; err: %v", from, err)
		}
	}
	return
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
