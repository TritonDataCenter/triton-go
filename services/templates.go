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
	"fmt"
	"net/http"
	"path"

	"github.com/joyent/triton-go/client"
	pkgerrors "github.com/pkg/errors"
)

type TemplatesClient struct {
	client *client.Client
}

type InstanceTemplate struct {
	ID                 int64
	TemplateName       string
	AccountId          string
	Package            string
	ImageId            string
	InstanceNamePrefix string
	FirewallEnabled    bool
	Networks           []string
	UserData           string
	MetaData           map[string]string
	Tags               map[string]string
}

type ListTemplatesInput struct{}

func (c *TemplatesClient) List(ctx context.Context, _ *ListTemplatesInput) ([]*InstanceTemplate, error) {
	fullPath := path.Join("/v1/tsg/templates")

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

	var results []*InstanceTemplate
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&results); err != nil {
		return nil, pkgerrors.Wrap(err, "unable to decode list templates response")
	}

	return results, nil
}

type GetTemplateInput struct {
	Name string
}

func (i *GetTemplateInput) Validate() error {
	if i.Name == "" {
		return fmt.Errorf("template name can not be empty")
	}

	return nil
}

func (c *TemplatesClient) Get(ctx context.Context, input *GetTemplateInput) (*InstanceTemplate, error) {
	if err := input.Validate(); err != nil {
		return nil, pkgerrors.Wrap(err, "unable to get instance template")
	}

	fullPath := path.Join("/v1/tsg/templates", input.Name)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequestTSG(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to get template")
	}

	var results *InstanceTemplate
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&results); err != nil {
		return nil, pkgerrors.Wrap(err, "unable to decode get template response")
	}

	return results, nil
}

type CreateTemplateInput struct {
	TemplateName       string
	AccountId          string
	Package            string
	ImageId            string
	InstanceNamePrefix string
	FirewallEnabled    bool
	Networks           []string
	UserData           string
	MetaData           map[string]string
	Tags               map[string]string
}

func (input *CreateTemplateInput) toAPI() map[string]interface{} {
	result := make(map[string]interface{})

	if input.TemplateName != "" {
		result["TemplateName"] = input.TemplateName
	}

	if input.AccountId != "" {
		result["AccountId"] = input.AccountId
	}

	if input.Package != "" {
		result["Package"] = input.Package
	}

	if input.ImageId != "" {
		result["ImageId"] = input.ImageId
	}

	if input.InstanceNamePrefix != "" {
		result["InstanceNamePrefix"] = input.InstanceNamePrefix
	}

	result["FirewallEnabled"] = input.FirewallEnabled

	if len(input.Networks) > 0 {
		result["networks"] = input.Networks
	}

	if input.UserData != "" {
		result["UserData"] = input.UserData
	}

	if len(input.MetaData) > 0 {
		result["MetaData"] = input.MetaData
	}

	if len(input.Tags) > 0 {
		result["Tags"] = input.Tags
	}

	return result
}

func (c *TemplatesClient) Create(ctx context.Context, input *CreateTemplateInput) error {
	fullPath := path.Join("/v1/tsg/templates")
	body := input.toAPI()

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Body:   body,
	}
	respReader, err := c.client.ExecuteRequestTSG(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return pkgerrors.Wrap(err, "unable to create template")
	}

	return nil
}

type DeleteTemplateInput struct {
	Name string
}

func (i *DeleteTemplateInput) Validate() error {
	if i.Name == "" {
		return fmt.Errorf("template name can not be empty")
	}

	return nil
}

func (c *TemplatesClient) Delete(ctx context.Context, input *DeleteTemplateInput) error {
	if err := input.Validate(); err != nil {
		return pkgerrors.Wrap(err, "unable to validate delete template input")
	}

	fullPath := path.Join("/v1/tsg/templates", input.Name)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequestTSG(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return pkgerrors.Wrap(err, "unable to delete template")
	}

	return nil
}
