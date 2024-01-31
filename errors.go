package main

import (
	"errors"
	"golang.org/x/exp/slog"
	"os"
	"runtime/debug"
)

var (
	exitCode   = 1
	printStack = false
)

func RecoverAndExit() {
	// Solidify exitCode
	defer os.Exit(exitCode)
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			if ep := new(ErrPanic); errors.As(err, &ep) {
				slog.Error("error", "error", ep.error)
			} else if errors.Is(err, ErrLogin) && exitCode == 0 {
				// Don't print stack
				return
			}
		} else if r != nil {
			slog.Error("caught panic", "recover", r)
		}
		if printStack {
			debug.PrintStack()
		}
	}
}

type ErrPanic struct{ error }

func (e *ErrPanic) Error() string { return e.error.Error() }

func Check(err error) {
	if err != nil {
		exitCode = 1
		panic(&ErrPanic{err})
	}
}

func Must[T any](t T, err error) T {
	Check(err)
	return t
}

func Defer(fn func() error) {
	if err := fn(); err != nil {
		exitCode = 1
		slog.Error("defer", "error", err)
	}
}

var ErrLogin = errors.New("")

func Login() {
	exitCode = 0
	panic(ErrLogin)
}
