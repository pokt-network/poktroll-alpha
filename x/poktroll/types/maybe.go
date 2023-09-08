package types

import (
	"reflect"
)

type Maybe[T any] struct {
	value T
	error error
}

func NewMaybe[T any](value T, error error) Maybe[T] {
	return Maybe[T]{value: value, error: error}
}

func Just[T any](value T) Maybe[T] {
	return Maybe[T]{value: value, error: nil}
}

func JustError[T any](error error) Maybe[T] {
	// see: https://stackoverflow.com/questions/73864711/get-type-parameter-from-a-generic-struct-using-reflection
	var zeroT [0]T
	typeT := reflect.TypeOf(zeroT).Elem()
	zeroValue := reflect.Zero(typeT).Interface().(T)
	return Maybe[T]{value: zeroValue, error: error}
}

func JustErrorChan[T any](err error) <-chan Maybe[T] {
	resultCh := make(chan Maybe[T])
	go func() {
		resultCh <- JustError[T](err)
	}()
	return resultCh
}

func JustValueChan[T any](value T) <-chan Maybe[T] {
	resultCh := make(chan Maybe[T])
	go func() {
		resultCh <- Just[T](value)
	}()
	return resultCh
}

func (m Maybe[T]) ValueOrError() (T, error) {
	return m.value, m.error
}

func (m Maybe[T]) Value() T {
	return m.value
}

func (m Maybe[T]) Error() error {
	return m.error
}
