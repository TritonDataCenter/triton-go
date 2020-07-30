//
// Copyright 2020 Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package account_test

import (
	"context"
	"fmt"
	"testing"

	triton "github.com/joyent/triton-go/v2"
	"github.com/joyent/triton-go/v2/account"
	"github.com/joyent/triton-go/v2/testutils"
)

// Placeholder for the generated AccessKey.AccessKeyID
var accessKeyId = ""
var accessKey *account.AccessKey

func TestAccAccessKey_Create(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{

			&testutils.StepClient{
				StateBagKey: "accesskey",
				CallFunc: func(config *triton.ClientConfig) (interface{}, error) {
					return account.NewClient(config)
				},
			},

			&testutils.StepAPICall{
				StateBagKey: "accesskey",
				CallFunc: func(client interface{}) (interface{}, error) {
					c := client.(*account.AccountClient)
					ctx := context.Background()
					input := &account.CreateAccessKeyInput{}
					created, err := c.AccessKeys().CreateAccessKey(ctx, input)
					if err != nil {
						return nil, err
					}
					fmt.Printf("[DEBUG] Created access key %+v\n", created)
					accessKeyId = created.AccessKeyID
					return created, nil
				},
				CleanupFunc: func(client interface{}, callState interface{}) {
					c := client.(*account.AccountClient)
					ctx := context.Background()
					input := &account.DeleteAccessKeyInput{
						AccessKeyID: accessKeyId,
					}
					c.AccessKeys().DeleteAccessKey(ctx, input)
				},
			},
			&testutils.StepAssertFunc{
				AssertFunc: func(state testutils.TritonStateBag) error {
					accessKey := state.Get("accesskey").(*account.AccessKey)
					if accessKey.AccessKeyID == "" {
						t.Fatalf("Expected to have some AccessKeyID value, have: \"%v\"", accessKey.AccessKeyID)
					}

					if accessKey.SecretAccessKey == "" {
						t.Fatalf("Expected to have some SecretAccessKey value, have: \"%v\"", accessKey.SecretAccessKey)
					}

					if accessKey.CreateDate.IsZero() {
						t.Fatalf("Expected access key CreatedDate to be non zero time, got: \"%v\"", accessKey.CreateDate)
					}
					return nil
				},
			},
		},
	})
}

func TestAccAccessKey_GetAndList(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{

			&testutils.StepClient{
				StateBagKey: "accesskey",
				CallFunc: func(config *triton.ClientConfig) (interface{}, error) {
					return account.NewClient(config)
				},
			},

			&testutils.StepAPICall{
				StateBagKey: "accesskey",
				CallFunc: func(client interface{}) (interface{}, error) {
					c := client.(*account.AccountClient)
					ctx := context.Background()
					input := &account.CreateAccessKeyInput{}
					created, err := c.AccessKeys().CreateAccessKey(ctx, input)
					if err != nil {
						return nil, err
					}
					fmt.Printf("[DEBUG] Created access key %+v\n", created)
					accessKeyId = created.AccessKeyID
					return created, nil
				},
				CleanupFunc: func(client interface{}, callState interface{}) {
					c := client.(*account.AccountClient)
					ctx := context.Background()
					input := &account.DeleteAccessKeyInput{
						AccessKeyID: accessKeyId,
					}
					c.AccessKeys().DeleteAccessKey(ctx, input)
				},
			},

			&testutils.StepAPICall{
				StateBagKey: "getAccessKey",
				CallFunc: func(client interface{}) (interface{}, error) {
					c := client.(*account.AccountClient)
					ctx := context.Background()
					input := &account.GetAccessKeyInput{
						AccessKeyID: accessKeyId,
					}
					retrieved, err := c.AccessKeys().GetAccessKey(ctx, input)
					if err != nil {
						return nil, err
					}
					if retrieved.AccessKeyID != accessKeyId {
						t.Fatalf("Expected to retrieve a key with AccessKeyID \"%s\", but got \"%s\"", retrieved.AccessKeyID, accessKeyId)
					}
					return retrieved, nil
				},
			},

			&testutils.StepAPICall{
				StateBagKey: "listAccessKeys",
				CallFunc: func(client interface{}) (interface{}, error) {
					c := client.(*account.AccountClient)
					ctx := context.Background()
					input := &account.ListAccessKeysInput{}
					listed, err := c.AccessKeys().ListAccessKeys(ctx, input)
					if err != nil {
						return nil, err
					}
					if len(listed) == 0 {
						t.Fatalf("Expected to retrieve a list of access keys, but got empty list")
					}
					if listed[0].AccessKeyID != accessKeyId {
						t.Fatalf("Expected to retrieve a key with AccessKeyID \"%s\", but got \"%s\"", listed[0].AccessKeyID, accessKeyId)
					}
					return listed, nil
				},
			},
		},
	})
}

func TestAccAccessKey_Delete(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{

			&testutils.StepClient{
				StateBagKey: "accesskey",
				CallFunc: func(config *triton.ClientConfig) (interface{}, error) {
					return account.NewClient(config)
				},
			},

			&testutils.StepAPICall{
				StateBagKey: "accesskey",
				CallFunc: func(client interface{}) (interface{}, error) {
					c := client.(*account.AccountClient)
					ctx := context.Background()
					input := &account.CreateAccessKeyInput{}
					created, err := c.AccessKeys().CreateAccessKey(ctx, input)
					if err != nil {
						return nil, err
					}
					fmt.Printf("[DEBUG] Created access key %+v\n", created)
					accessKeyId = created.AccessKeyID
					return created, nil
				},
				CleanupFunc: func(client interface{}, callState interface{}) {
					c := client.(*account.AccountClient)
					ctx := context.Background()
					input := &account.DeleteAccessKeyInput{
						AccessKeyID: accessKeyId,
					}
					c.AccessKeys().DeleteAccessKey(ctx, input)
				},
			},

			&testutils.StepAPICall{
				StateBagKey: "noop",
				CallFunc: func(client interface{}) (interface{}, error) {
					c := client.(*account.AccountClient)
					ctx := context.Background()
					input := &account.DeleteAccessKeyInput{
						AccessKeyID: accessKeyId,
					}
					return nil, c.AccessKeys().DeleteAccessKey(ctx, input)
				},
			},

			&testutils.StepAPICall{
				ErrorKey: "getAccessKeyError",
				CallFunc: func(client interface{}) (interface{}, error) {
					c := client.(*account.AccountClient)
					ctx := context.Background()
					input := &account.GetAccessKeyInput{
						AccessKeyID: accessKeyId,
					}
					return c.AccessKeys().GetAccessKey(ctx, input)
				},
			},

			&testutils.StepAssertTritonError{
				ErrorKey: "getAccessKeyError",
				Code:     "ResourceNotFound",
			},
		},
	})
}
