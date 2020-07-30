//
// Copyright 2020 Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package account

import (
	"context"
	"encoding/json"
	"net/http"
	"path"
	"time"

	"github.com/joyent/triton-go/v2/client"
	"github.com/pkg/errors"
)

// AccessKeysClient holds a pointer to triton client
type AccessKeysClient struct {
	client *client.Client
}

// AccessKey represents an access key
type AccessKey struct {
	// AccessKeyId id of the key
	AccessKeyID string `json:"accesskeyid"`

	// SecretAccessKey the secret used for signing requests
	SecretAccessKey string `json:"accesskeysecret"`

	// CreateDate
	CreateDate time.Time `json:"created"`

	// Status either "Active" or "Inactive"
	Status string `json:"status"`

	// UserName the uuid of the user the access key is associated with
	UserName string
}

// ListAccessKeysInput empty struct payload for ListAccessKeys
type ListAccessKeysInput struct{}

// ListAccessKeys lists all access keys we have on record for the specified
// account/user.
func (c *AccessKeysClient) ListAccessKeys(ctx context.Context, _ *ListAccessKeysInput) ([]*AccessKey, error) {
	fullPath := path.Join("/", c.client.AccountName, "accesskeys")
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to list accesskeys")
	}

	var result []*AccessKey
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode list accesskeys response")
	}
	for _, elm := range result {
		elm.UserName = c.client.AccountName
	}
	return result, nil
}

// GetAccessKeyInput payload for GetAccessKey
type GetAccessKeyInput struct {
	AccessKeyID string
}

// GetAccessKey returns an access key with the provided AccessKeyID
func (c *AccessKeysClient) GetAccessKey(ctx context.Context, input *GetAccessKeyInput) (*AccessKey, error) {
	fullPath := path.Join("/", c.client.AccountName, "accesskeys", input.AccessKeyID)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to get key")
	}

	var result *AccessKey
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode get key response")
	}
	result.UserName = c.client.AccountName

	return result, nil
}

// DeleteAccessKeyInput payload for DeleteAccessKey
type DeleteAccessKeyInput struct {
	AccessKeyID string
}

// DeleteAccessKey with the provided AccessKeyID
func (c *AccessKeysClient) DeleteAccessKey(ctx context.Context, input *DeleteAccessKeyInput) error {
	fullPath := path.Join("/", c.client.AccountName, "accesskeys", input.AccessKeyID)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errors.Wrap(err, "unable to delete key")
	}

	return nil
}

// CreateAccessKeyInput is the empty payload used when creating a new access key.
type CreateAccessKeyInput struct{}

// CreateAccessKey generates a new AccessKey with a new secret
func (c *AccessKeysClient) CreateAccessKey(ctx context.Context, input *CreateAccessKeyInput) (*AccessKey, error) {
	fullPath := path.Join("/", c.client.AccountName, "accesskeys")
	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Body:   input,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to create access key")
	}

	var result *AccessKey
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errors.Wrap(err, "unable to decode create access key response")
	}

	result.UserName = c.client.AccountName

	return result, nil
}
