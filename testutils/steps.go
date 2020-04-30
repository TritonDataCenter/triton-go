//
// Copyright 2020 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package testutils

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/abdullin/seq"
	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/compute"
	"github.com/joyent/triton-go/errors"
	"github.com/joyent/triton-go/network"
	pkgerrors "github.com/pkg/errors"
)

type StepClient struct {
	StateBagKey string
	ErrorKey    string
	CallFunc    func(config *triton.ClientConfig) (interface{}, error)
	CleanupFunc func(client interface{}, callState interface{})
}

func (s *StepClient) Run(state TritonStateBag) StepAction {
	client, err := s.CallFunc(state.Config())
	if err != nil {
		if s.ErrorKey == "" {
			state.AppendError(err)
			return Halt
		}

		state.Put(s.ErrorKey, err)
		return Continue
	}

	state.PutClient(client)
	return Continue
}

func (s *StepClient) Cleanup(state TritonStateBag) {
	return
}

type StepComputeClient struct {
	StateBagKey string
	ErrorKey    string
	CallFunc    func(state TritonStateBag, config *compute.ComputeClient) (interface{}, error)
	CleanupFunc func(state TritonStateBag, config *compute.ComputeClient)
}

func (s *StepComputeClient) Run(state TritonStateBag) StepAction {
	computeClient, err := compute.NewClient(state.Config())
	if err != nil {
		state.AppendError(err)
		return Halt
	}
	state.Put("computeClient", computeClient)

	result, err := s.CallFunc(state, computeClient)
	if err != nil {
		if s.ErrorKey == "" {
			state.AppendError(err)
			return Halt
		}

		state.Put(s.ErrorKey, err)
		return Continue
	}

	state.Put(s.StateBagKey, result)
	return Continue
}

func (s *StepComputeClient) Cleanup(state TritonStateBag) {
	if s.CleanupFunc == nil {
		return
	}

	computeClient := state.Get("computeClient").(*compute.ComputeClient)
	s.CleanupFunc(state, computeClient)
}

type StepNetworkClient struct {
	StateBagKey string
	ErrorKey    string
	CallFunc    func(state TritonStateBag, config *network.NetworkClient) (interface{}, error)
	CleanupFunc func(client interface{}, callState interface{})
}

func (s *StepNetworkClient) Run(state TritonStateBag) StepAction {
	networkClient, err := network.NewClient(state.Config())
	if err != nil {
		state.AppendError(err)
		return Halt
	}
	state.Put("networkClient", networkClient)

	result, err := s.CallFunc(state, networkClient)
	if err != nil {
		if s.ErrorKey == "" {
			state.AppendError(err)
			return Halt
		}

		state.Put(s.ErrorKey, err)
		return Continue
	}

	state.Put(s.StateBagKey, result)
	return Continue
}

func (s *StepNetworkClient) Cleanup(state TritonStateBag) {
	if s.CleanupFunc == nil {
		return
	}

	if callState, ok := state.GetOk(s.StateBagKey); ok {
		s.CleanupFunc(state.Client(), callState)
	} else {
		log.Print("[INFO] No state for API call, calling cleanup with nil call state")
		s.CleanupFunc(state.Client(), nil)
	}
}

type StepGetImage struct {
	StateBagKey string
}

func (s *StepGetImage) Run(state TritonStateBag) StepAction {
	const imageName = "ubuntu-16.04"

	computeClient, err := compute.NewClient(state.Config())
	if err != nil {
		state.AppendError(err)
		return Halt
	}

	images, err := computeClient.Images().List(context.Background(), &compute.ListImagesInput{
		Name: imageName,
	})
	if err != nil {
		state.AppendError(err)
		return Halt
	}

	if len(images) == 0 {
		state.AppendError(pkgerrors.Errorf("No images matching image name %s",
			imageName))
		return Halt
	}

	// Images will be sorted by creation date - so we want the
	// most recent (the last) one.
	img := images[len(images)-1]

	// log.Printf("[DEBUG] Test img %+v", img)

	state.Put(s.StateBagKey, img)
	return Continue
}

func (s *StepGetImage) Cleanup(state TritonStateBag) {
	return
}

type StepGetExternalNetwork struct {
	StateBagKey string
}

func (s *StepGetExternalNetwork) Run(state TritonStateBag) StepAction {
	networkClient, err := network.NewClient(state.Config())
	if err != nil {
		state.AppendError(err)
		return Halt
	}

	nets, err := networkClient.List(context.Background(), &network.ListInput{})
	if err != nil {
		state.AppendError(err)
		return Halt
	}

	var net *network.Network
	// Take the first public network.
	for _, found := range nets {
		if found.Public == true {
			net = found
			break
		}
	}

	if net == nil {
		state.AppendError(pkgerrors.New("Unable to find external network"))
		return Halt
	}

	state.Put(s.StateBagKey, net)
	return Continue
}

func (s *StepGetExternalNetwork) Cleanup(state TritonStateBag) {
	return
}

type StepGetPackage struct {
	StateBagKey string
}

func (s *StepGetPackage) Run(state TritonStateBag) StepAction {
	computeClient, err := compute.NewClient(state.Config())
	if err != nil {
		state.AppendError(err)
		return Halt
	}

	pkgs, err := computeClient.Packages().List(context.Background(), &compute.ListPackagesInput{})
	if err != nil {
		state.AppendError(err)
		return Halt
	}

	var pkg *compute.Package
	var minMem int64 = 128
	var maxMem int64 = 1024
	// Take the smallest generic package that is in the range 128MB-1024MB memory.
	for _, found := range pkgs {
		// log.Printf("[INFO] Pkg %+v", found)
		if minMem <= found.Memory && found.Memory <= maxMem {
			if pkg == nil || pkg.Memory > found.Memory {
				pkg = found
			}
		}
	}

	if pkg == nil {
		state.AppendError(pkgerrors.Errorf(
			"Unable to find a package with %d to %d MB memory", minMem, maxMem))
		return Halt
	}

	// log.Printf("[DEBUG] Test pkg %+v", pkg)

	state.Put(s.StateBagKey, pkg)
	return Continue
}

func (s *StepGetPackage) Cleanup(state TritonStateBag) {
	return
}

type StepAPICall struct {
	StateBagKey string
	ErrorKey    string
	CallFunc    func(client interface{}) (interface{}, error)
	CleanupFunc func(client interface{}, callState interface{})
}

func (s *StepAPICall) Run(state TritonStateBag) StepAction {
	result, err := s.CallFunc(state.Client())
	if err != nil {
		if s.ErrorKey == "" {
			state.AppendError(err)
			return Halt
		}

		state.Put(s.ErrorKey, err)
		return Continue
	}

	state.Put(s.StateBagKey, result)
	return Continue
}

func (s *StepAPICall) Cleanup(state TritonStateBag) {
	if s.CleanupFunc == nil {
		return
	}

	if callState, ok := state.GetOk(s.StateBagKey); ok {
		s.CleanupFunc(state.Client(), callState)
	} else {
		log.Print("[INFO] No state for API call, calling cleanup with nil call state")
		s.CleanupFunc(state.Client(), nil)
	}
}

type AssertFunc func(TritonStateBag) error

type StepAssertFunc struct {
	AssertFunc AssertFunc
}

func (s *StepAssertFunc) Run(state TritonStateBag) StepAction {
	if s.AssertFunc == nil {
		state.AppendError(pkgerrors.New("StepAssertFunc may not have a nil AssertFunc"))
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
		log.Printf("[INFO] Asserting %q has value \"%v\"", path, v)
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

type StepAssertSet struct {
	StateBagKey string
	Keys        []string
}

func (s *StepAssertSet) Run(state TritonStateBag) StepAction {
	actual, ok := state.GetOk(s.StateBagKey)
	if !ok {
		state.AppendError(fmt.Errorf("Key %q not found in state", s.StateBagKey))
	}

	var pass = true
	for _, key := range s.Keys {
		r := reflect.ValueOf(actual)
		f := reflect.Indirect(r).FieldByName(key)

		log.Printf("[INFO] Asserting %q has a non-zero value", key)
		if f.Interface() == reflect.Zero(reflect.TypeOf(f)).Interface() {
			err := fmt.Sprintf("Expected %q to have a non-zero value", key)
			state.AppendError(fmt.Errorf(err))
			pass = false
		}
	}

	if !pass {
		return Halt
	}

	return Continue
}

func (s *StepAssertSet) Cleanup(state TritonStateBag) {
	return
}

type StepAssertTritonError struct {
	ErrorKey string
	Code     string
}

func (s *StepAssertTritonError) Run(state TritonStateBag) StepAction {
	err, ok := state.GetOk(s.ErrorKey)
	if !ok {
		state.AppendError(fmt.Errorf("Expected APIError %q to be in state", s.Code))
		return Halt
	}

	switch err := pkgerrors.Cause(err.(error)).(type) {
	case *errors.APIError:
		if err.Code == s.Code {
			return Continue
		}
		state.AppendError(fmt.Errorf("Expected APIError code %q to be in state key %q, was %q", s.Code, s.ErrorKey, err.Code))
		return Halt
	default:
		state.AppendError(pkgerrors.New("Expected a APIError in wrapped error chain"))
	}

	return Halt
}

func (s *StepAssertTritonError) Cleanup(state TritonStateBag) {
	return
}
