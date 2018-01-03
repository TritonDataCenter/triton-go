package identity_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/joyent/triton-go/identity"
	"github.com/joyent/triton-go/testutils"
)

var (
	listPoliciesErrorType = errors.New("Error executing ListPolicies request:")
	getPolicyErrorType    = errors.New("Error executing GetPolicy request:")
	deletePolicyErrorType = errors.New("Error executing DeletePolicy request:")
	updatePolicyErrorType = errors.New("Error executing UpdatePolicy request:")
	createPolicyErrorType = errors.New("Error executing CreatePolicy request:")
)

func TestListPolicies(t *testing.T) {
	identityClient := MockIdentityClient()

	do := func(ctx context.Context, ic *identity.IdentityClient) ([]*identity.Policy, error) {
		defer testutils.DeactivateClient()

		policies, err := ic.Policies().List(ctx, &identity.ListPoliciesInput{})
		if err != nil {
			return nil, err
		}
		return policies, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("GET", fmt.Sprintf("/%s/policies", accountUrl), listPoliciesSuccess)

		resp, err := do(context.Background(), identityClient)
		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatalf("Expected an output but got nil")
		}
	})

	t.Run("eof", func(t *testing.T) {
		testutils.RegisterResponder("GET", fmt.Sprintf("/%s/policies", accountUrl), listPoliciesEmpty)

		_, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "EOF") {
			t.Errorf("expected error to contain EOF: found %s", err)
		}
	})

	t.Run("bad_decode", func(t *testing.T) {
		testutils.RegisterResponder("GET", fmt.Sprintf("/%s/policies", accountUrl), listPoliciesBadeDecode)

		_, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "invalid character") {
			t.Errorf("expected decode to fail: found %s", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("GET", fmt.Sprintf("/%s/policies", accountUrl), listPoliciesError)

		resp, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Error("expected resp to be nil")
		}

		if !strings.Contains(err.Error(), "Error executing ListPolicies request:") {
			t.Errorf("expected error to equal testError: found %s", err)
		}
	})
}

func TestGetPolicy(t *testing.T) {
	identityClient := MockIdentityClient()

	do := func(ctx context.Context, ic *identity.IdentityClient) (*identity.Policy, error) {
		defer testutils.DeactivateClient()

		user, err := ic.Policies().Get(ctx, &identity.GetPolicyInput{
			PolicyID: "95ca7b25-5c8f-4c1b-92da-4276f23807f3",
		})
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("GET", fmt.Sprintf("/%s/policies/%s", accountUrl, "95ca7b25-5c8f-4c1b-92da-4276f23807f3"), getPolicySuccess)

		resp, err := do(context.Background(), identityClient)
		if err != nil {
			t.Fatal(err)
		}

		if resp == nil {
			t.Fatalf("Expected an output but got nil")
		}
	})

	t.Run("eof", func(t *testing.T) {
		testutils.RegisterResponder("GET", fmt.Sprintf("/%s/policies/%s", accountUrl, "95ca7b25-5c8f-4c1b-92da-4276f23807f3"), getPolicyEmpty)

		_, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "EOF") {
			t.Errorf("expected error to contain EOF: found %s", err)
		}
	})

	t.Run("bad_decode", func(t *testing.T) {
		testutils.RegisterResponder("GET", fmt.Sprintf("/%s/policies/%s", accountUrl, "95ca7b25-5c8f-4c1b-92da-4276f23807f3"), getPolicyBadeDecode)

		_, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "invalid character") {
			t.Errorf("expected decode to fail: found %s", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("GET", fmt.Sprintf("/%s/policies", accountUrl), getPolicyError)

		resp, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Error("expected resp to be nil")
		}

		if !strings.Contains(err.Error(), "Error executing GetPolicy request:") {
			t.Errorf("expected error to equal testError: found %s", err)
		}
	})
}

func TestDeletePolicy(t *testing.T) {
	identityClient := MockIdentityClient()

	do := func(ctx context.Context, ic *identity.IdentityClient) error {
		defer testutils.DeactivateClient()

		return ic.Policies().Delete(ctx, &identity.DeletePolicyInput{
			PolicyID: "8700e959-4cb3-4337-8afa-fb0a53b5366e",
		})
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("DELETE", fmt.Sprintf("/%s/policies/%s", accountUrl, "8700e959-4cb3-4337-8afa-fb0a53b5366e"), deletePolicySuccess)

		err := do(context.Background(), identityClient)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("DELETE", fmt.Sprintf("/%s/policies", accountUrl), deletePolicyError)

		err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "Error executing DeletePolicy request:") {
			t.Errorf("expected error to equal testError: found %s", err)
		}
	})
}

func TestUpdatePolicy(t *testing.T) {
	identityClient := MockIdentityClient()

	do := func(ctx context.Context, ic *identity.IdentityClient) (*identity.Policy, error) {
		defer testutils.DeactivateClient()

		policy, err := ic.Policies().Update(ctx, &identity.UpdatePolicyInput{
			PolicyID:    "95ca7b25-5c8f-4c1b-92da-4276f23807f3",
			Description: "Updated Description",
			Name:        "Updated Name",
		})
		if err != nil {
			return nil, err
		}
		return policy, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("POST", fmt.Sprintf("/%s/policies/%s", accountUrl, "95ca7b25-5c8f-4c1b-92da-4276f23807f3"), updatePolicySuccess)

		_, err := do(context.Background(), identityClient)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("POST", fmt.Sprintf("/%s/policies/%s", accountUrl, "95ca7b25-5c8f-4c1b-92da-4276f23807f3"), updatePolicyError)

		_, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "Error executing UpdatePolicy request:") {
			t.Errorf("expected error to equal testError: found %s", err)
		}
	})
}

func TestCreatePolocy(t *testing.T) {
	identityClient := MockIdentityClient()

	do := func(ctx context.Context, ic *identity.IdentityClient) (*identity.Policy, error) {
		defer testutils.DeactivateClient()

		policy, err := ic.Policies().Create(ctx, &identity.CreatePolicyInput{
			Name:        "Test Policy",
			Description: "Test Description",
			Rules:       []string{"CAN rebootmachine"},
		})

		if err != nil {
			return nil, err
		}
		return policy, nil
	}

	t.Run("successful", func(t *testing.T) {
		testutils.RegisterResponder("POST", fmt.Sprintf("/%s/policies", accountUrl), createPolicySuccess)

		_, err := do(context.Background(), identityClient)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("error", func(t *testing.T) {
		testutils.RegisterResponder("POST", fmt.Sprintf("/%s/policies", accountUrl), createPolicyError)

		_, err := do(context.Background(), identityClient)
		if err == nil {
			t.Fatal(err)
		}

		if !strings.Contains(err.Error(), "Error executing CreatePolicy request:") {
			t.Errorf("expected error to equal testError: found %s", err)
		}
	})
}

func listPoliciesEmpty(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}, nil
}

func listPoliciesSuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`[
	{
    "name": "readinstance",
    "id": "95ca7b25-5c8f-4c1b-92da-4276f23807f3",
    "rules": [
      "can listmachine and getmachine"
    ]
  },
  {
    "name": "createinstance",
    "id": "95ca7b25-5c8f-4c1b-92da-4276f23805ds",
    "rules": [
      "can createinstance"
    ]
  }
]`)
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func listPoliciesBadeDecode(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{[
	{
    "name": "createinstance",
    "id": "95ca7b25-5c8f-4c1b-92da-4276f23805ds",
    "rules": [
      "can createinstance"
    ]
  }
]}`)
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func listPoliciesError(req *http.Request) (*http.Response, error) {
	return nil, listPoliciesErrorType
}

func getPolicySuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
    "name": "readinstance",
    "id": "95ca7b25-5c8f-4c1b-92da-4276f23807f3",
    "rules": [
      "can listmachine and getmachine"
    ]
  }
`)
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func getPolicyError(req *http.Request) (*http.Response, error) {
	return nil, getPolicyErrorType
}

func getPolicyBadeDecode(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
    "name": "readinstance",
    "id": "95ca7b25-5c8f-4c1b-92da-4276f23807f3",
    "rules": [
      "can listmachine and getmachine"
    ],
  }`)
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func getPolicyEmpty(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}, nil
}

func deletePolicySuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	return &http.Response{
		StatusCode: 204,
		Header:     header,
	}, nil
}

func deletePolicyError(req *http.Request) (*http.Response, error) {
	return nil, deletePolicyErrorType
}

func updatePolicySuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
  "name": "Updated Name",
  "id": "95ca7b25-5c8f-4c1b-92da-4276f23807f3",
  "rules": [
    "can rebootMachine"
  ],
  "description": "Updated Description"
}
`)

	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func updatePolicyError(req *http.Request) (*http.Response, error) {
	return nil, updatePolicyErrorType
}

func createPolicySuccess(req *http.Request) (*http.Response, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body := strings.NewReader(`{
  "name": "Test Policy",
  "id": "8700e959-4cb3-4337-8afa-fb0a53b5366e",
  "rules": [
    "CAN rebootmachine"
  ],
  "description": "Test Description"
}
`)

	return &http.Response{
		StatusCode: 201,
		Header:     header,
		Body:       ioutil.NopCloser(body),
	}, nil
}

func createPolicyError(req *http.Request) (*http.Response, error) {
	return nil, createPolicyErrorType
}
