// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/glestaris/ice-clique/dispatcher"
)

type FakeInterruptible struct {
	InterruptStub        func()
	interruptMutex       sync.RWMutex
	interruptArgsForCall []struct{}
	ResumeStub           func()
	resumeMutex          sync.RWMutex
	resumeArgsForCall    []struct{}
}

func (fake *FakeInterruptible) Interrupt() {
	fake.interruptMutex.Lock()
	fake.interruptArgsForCall = append(fake.interruptArgsForCall, struct{}{})
	fake.interruptMutex.Unlock()
	if fake.InterruptStub != nil {
		fake.InterruptStub()
	}
}

func (fake *FakeInterruptible) InterruptCallCount() int {
	fake.interruptMutex.RLock()
	defer fake.interruptMutex.RUnlock()
	return len(fake.interruptArgsForCall)
}

func (fake *FakeInterruptible) Resume() {
	fake.resumeMutex.Lock()
	fake.resumeArgsForCall = append(fake.resumeArgsForCall, struct{}{})
	fake.resumeMutex.Unlock()
	if fake.ResumeStub != nil {
		fake.ResumeStub()
	}
}

func (fake *FakeInterruptible) ResumeCallCount() int {
	fake.resumeMutex.RLock()
	defer fake.resumeMutex.RUnlock()
	return len(fake.resumeArgsForCall)
}

var _ dispatcher.Interruptible = new(FakeInterruptible)