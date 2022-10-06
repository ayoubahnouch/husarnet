package main

type dumbObserver[T comparable] struct {
	lastValue  T
	changeHook func(oldValue, NewValue T)
}

func NewDumbObserver[T comparable](changeHook func(oldValue, NewValue T)) *dumbObserver[T] {
	observer := dumbObserver[T]{
		changeHook: changeHook,
	}

	return &observer
}

func (o *dumbObserver[T]) Update(newValue T) {
	if newValue != o.lastValue {
		o.changeHook(o.lastValue, newValue)
		o.lastValue = newValue
	}
}
