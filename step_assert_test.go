package triton

import (
	"errors"
	"fmt"
	"log"

	"github.com/abdullin/seq"
	"github.com/hashicorp/errwrap"
)

type AssertFunc func(TritonStateBag) error

type StepAssertFunc struct {
	AssertFunc AssertFunc
}

func (s *StepAssertFunc) Run(state TritonStateBag) StepAction {
	if s.AssertFunc == nil {
		state.AppendError(errors.New("StepAssertFunc may not have a nil AssertFunc"))
		return Halt
	}

	err := s.AssertFunc(state)
	if err != nil {
		state.AppendError(err)
		return Halt
	}

	return Continue
}

func (s *StepAssertFunc) Cleanup(state TritonStateBag) {
	return
}

type StepAssert struct {
	StateBagKey string
	Assertions  seq.Map
}

func (s *StepAssert) Run(state TritonStateBag) StepAction {
	actual, ok := state.GetOk(s.StateBagKey)
	if !ok {
		state.AppendError(fmt.Errorf("Key %q not found in state", s.StateBagKey))
	}

	for k, v := range s.Assertions {
		path := fmt.Sprintf("%s.%s", s.StateBagKey, k)
		log.Printf("[INFO] Asserting %q has value \"%v\"...", path, v)
	}

	result := s.Assertions.Test(actual)

	if result.Ok() {
		return Continue
	}

	for _, v := range result.Issues {
		err := fmt.Sprintf("Expected %q to be \"%v\" but got %q",
			v.Path,
			v.ExpectedValue,
			v.ActualValue,
		)
		state.AppendError(fmt.Errorf(err))
	}

	return Halt
}

func (s *StepAssert) Cleanup(state TritonStateBag) {
	return
}

type StepAssertTritonError struct {
	ErrorKey string
	Code     string
}

func (s *StepAssertTritonError) Run(state TritonStateBag) StepAction {
	err, ok := state.GetOk(s.ErrorKey)
	if !ok {
		state.AppendError(fmt.Errorf("Expected TritonError %q to be in state", s.Code))
		return Halt
	}

	tritonErrorInterface := errwrap.GetType(err.(error), &TritonError{})
	if tritonErrorInterface == nil {
		state.AppendError(errors.New("Expected a TritonError in wrapped error chain"))
		return Halt
	}

	tritonErr := tritonErrorInterface.(*TritonError)
	if tritonErr.Code == s.Code {
		return Continue
	}

	state.AppendError(fmt.Errorf("Expected TritonError code %q to be in state key %q, was %q", s.Code, s.ErrorKey, tritonErr.Code))
	return Halt
}

func (s *StepAssertTritonError) Cleanup(state TritonStateBag) {
	return
}
