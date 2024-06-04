package utils

import "github.com/diamondburned/gotk4/pkg/glib/v2"

type TaskSequence[W any] struct {
	sequenceEnded func()
	functions     []func() (W, error)
	callbacks     []func(W, error) bool
}

func NewTaskSequence[W any](endCallback func()) *TaskSequence[W] {
	ts := &TaskSequence[W]{}
	ts.functions = make([]func() (W, error), 0)
	ts.callbacks = make([]func(W, error) bool, 0)
	ts.sequenceEnded = endCallback

	return ts
}

func (ts *TaskSequence[W]) Add(function func() (W, error), callback func(W, error) bool) {
	ts.functions = append(ts.functions, function)
	ts.callbacks = append(ts.callbacks, callback)
}

func (ts *TaskSequence[W]) Execute() {
	ts.executeNext(true)
}

func (ts *TaskSequence[W]) executeNext(keepWorking bool) {
	if len(ts.functions) == 0 || !keepWorking {
		ts.sequenceEnded()
		return
	}

	go func() {
		function := ts.functions[0]
		callback := ts.callbacks[0]

		ts.functions[0] = ts.functions[len(ts.functions)-1]
		ts.functions = ts.functions[:len(ts.functions)-1]
		ts.callbacks[0] = ts.callbacks[len(ts.callbacks)-1]
		ts.callbacks = ts.callbacks[:len(ts.callbacks)-1]

		data, err := function()
		callbackInMainThreadContinue(data, err, callback, ts.executeNext)
	}()
}

func callbackInMainThreadContinue[T any](data T, err error, function func(T, error) bool, callback func(bool)) {
	glib.IdleAdd(func() bool {
		response := function(data, err)
		callback(response)
		return false
	})
}
