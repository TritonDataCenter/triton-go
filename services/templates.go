//
// Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package services

import (
	"context"
	"encoding/json"
	"net/http"
	"path"
	"sort"

	"github.com/joyent/triton-go/client"
)

type TemplatesClient struct {
	client *client.Client
}

type Template struct {
	Name string `json:"name"`
}

type ListTemplatesInput struct{}

func (c *TemplatesClient) List(ctx context.Context, _ *ListTemplatesInput) ([]*Template, error) {
	fullPath := path.Join("/", c.client.AccountName, "templates")

	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequestTSG(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to list templates")
	}

	var intermediate map[string]string
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&intermediate); err != nil {
		return nil, pkgerrors.Wrap(err, "unable to decode list templates response")
	}

	keys := make([]string, len(intermediate))
	i := 0
	for k := range intermediate {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	result := make([]*Template, len(intermediate))
	i = 0
	for _, key := range keys {
		result[i] = &Template{
			Name: key,
		}
		i++
	}

	return result, nil
}
