//
// Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package identity_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/joyent/triton-go/identity"
	"github.com/joyent/triton-go/testutils"
)

var (
	fakeRoleId          = "1234562335"
	listRolesErrorType  = errors.New("unable to list roles")
	getRoleErrorType    = errors.New("unable to get role")
	deleteRoleErrorType = errors.New("unable to delete role")
	createRoleErrorType = errors.New("unable to create role")
	updateRoleErrorType = errors.New("unable to update role")
)

func TestDeleteRole(t *testing.T) {
	identityClient := MockIdentityClient()

	do := func(ctx context.Context, ic *identity.IdentityClient) error {
		defer testutils.DeactivateClient()

		return ic.Roles().Delete(ctx, &identity.DeleteRoleInput{
			RoleID: fakeRoleId,
		})
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("DELETE", path.Join("/", accountUrl, "roles", fakeRoleId), deleteRoleSuccess)

		err := do(context.Background(), identityClient)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("DELETE", path.Join("/", accountUrl, "roles"), deleteRoleError)

		err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "unable to delete role") {
			t.Errorf("expected error to equal testError: found %s", err)
		}
	})
}

func TestCreateRole(t *testing.T) {
	identityClient := MockIdentityClient()

	do := func(ctx context.Context, ic *identity.IdentityClient) (*identity.Role, error) {
		defer testutils.DeactivateClient()

		role, err := ic.Roles().Create(ctx, &identity.CreateRoleInput{
			Name: "readable",
		})

		if err != nil {
			return nil, err
		}
		return role, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("POST", path.Join("/", accountUrl, "roles"), createRoleSuccess)

		_, err := do(context.Background(), identityClient)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("POST", path.Join("/", accountUrl, "roles"), createRoleError)

		_, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "unable to create role") {
			t.Errorf("expected error to equal testError: found %s", err)
		}
	})
}

func TestUpdateRole(t *testing.T) {
	identityClient := MockIdentityClient()

	do := func(ctx context.Context, ic *identity.IdentityClient) (*identity.Role, error) {
		defer testutils.DeactivateClient()

		user, err := ic.Roles().Update(ctx, &identity.UpdateRoleInput{
			RoleID: "e53b8fec-e661-4ded-a21e-959c9ba08cb2",
			Name:   "updated-role-name",
		})
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("POST", path.Join("/", accountUrl, "roles", "e53b8fec-e661-4ded-a21e-959c9ba08cb2"), updateRoleSuccess)

		_, err := do(context.Background(), identityClient)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("POST", path.Join("/", accountUrl, "roles", "e53b8fec-e661-4ded-a21e-959c9ba08cb2"), updateRoleError)

		_, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "unable to update role") {
			t.Errorf("expected error to equal testError: found %s", err)
		}
	})
}

func TestListRoles(t *testing.T) {
	identityClient := MockIdentityClient()

	do := func(ctx context.Context, ic *identity.IdentityClient) ([]*identity.Role, error) {
		defer testutils.DeactivateClient()

		roles, err := ic.Roles().List(ctx, &identity.ListRolesInput{})
		if err != nil {
			return nil, err
		}
		return roles, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountUrl, "roles"), listRolesSuccess)

		resp, err := do(context.Background(), identityClient)
		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatalf("Expected an output but got nil")
		}
	})

	t.Run("eof", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountUrl, "roles"), listRolesEmpty)

		_, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "EOF") {
			t.Errorf("expected error to contain EOF: found %s", err)
		}
	})

	t.Run("bad_decode", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountUrl, "roles"), listRolesBadeDecode)

		_, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "invalid character") {
			t.Errorf("expected decode to fail: found %s", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountUrl, "roles"), listRolesError)

		resp, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Error("expected resp to be nil")
		}

		if !strings.Contains(err.Error(), "unable to list roles") {
			t.Errorf("expected error to equal testError: found %s", err)
		}
	})
}

func TestGetRole(t *testing.T) {
	identityClient := MockIdentityClient()

	do := func(ctx context.Context, ic *identity.IdentityClient) (*identity.Role, error) {
		defer testutils.DeactivateClient()

		role, err := ic.Roles().Get(ctx, &identity.GetRoleInput{
			RoleID: fakeRoleId,
		})
		if err != nil {
			return nil, err
		}
		return role, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountUrl, "roles", fakeRoleId), getRoleSuccess)

		resp, err := do(context.Background(), identityClient)
		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatalf("Expected an output but got nil")
		}
	})

	t.Run("eof", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountUrl, "roles", fakeRoleId), getRoleEmpty)

		_, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "EOF") {
			t.Errorf("expected error to contain EOF: found %s", err)
		}
	})

	t.Run("bad_decode", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountUrl, "roles", fakeRoleId), getRoleBadeDecode)

		_, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "invalid character") {
			t.Errorf("expected decode to fail: found %s", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("GET", path.Join("/", accountUrl, "roles"), getRoleError)

		resp, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Error("expected resp to be nil")
		}

		if !strings.Contains(err.Error(), "unable to get role") {
			t.Errorf("expected error to equal testError: found %s", err)
		}
	})
}

func deleteRoleSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	return &http.Response{
		StatusCode: 204,
		Header:     header,
	}, nil
}

func deleteRoleError(req *http.Request) (*http.Response, error) {
	return nil, deleteRoleErrorType
}

func createRoleSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
  "name": "readable",
  "id": "e53b8fec-e661-4ded-a21e-959c9ba08cb2"
}
`)

	return &http.Response{
		StatusCode: 201,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func createRoleError(req *http.Request) (*http.Response, error) {
	return nil, createRoleErrorType
}

func updateRoleSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
  "name": "updated-role-name",
  "id": "e53b8fec-e661-4ded-a21e-959c9ba08cb2"
}
`)

	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func updateRoleError(req *http.Request) (*http.Response, error) {
	return nil, updateRoleErrorType
}

func listRolesEmpty(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}, nil
}

func listRolesSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`[
	{
    "name": "readable",
    "id": "e53b8fec-e661-4ded-a21e-959c9ba08cb2",
    "members": [
      "foo"
    ],
    "default_members": [
      "foo"
    ],
    "policies": [
      "readinstance"
    ]
  }
]`)
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func listRolesBadeDecode(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{[
	{
    "name": "readable",
    "id": "e53b8fec-e661-4ded-a21e-959c9ba08cb2",
    "members": [
      "foo"
    ],
    "default_members": [
      "foo"
    ],
    "policies": [
      "readinstance"
    ]
  }
]}`)
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func listRolesError(req *http.Request) (*http.Response, error) {
	return nil, listRolesErrorType
}

func getRoleSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
    "name": "readable",
    "id": "e53b8fec-e661-4ded-a21e-959c9ba08cb2",
    "members": [
      "foo"
    ],
    "default_members": [
      "foo"
    ],
    "policies": [
      "readinstance"
    ]
  }
`)
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func getRoleError(req *http.Request) (*http.Response, error) {
	return nil, getRoleErrorType
}

func getRoleBadeDecode(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
    "name": "readable",
    "id": "e53b8fec-e661-4ded-a21e-959c9ba08cb2",
    "members": [
      "foo"
    ],
    "default_members": [
      "foo"
    ],
    "policies": [
      "readinstance"
    ],
  }`)
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func getRoleEmpty(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}, nil
}
