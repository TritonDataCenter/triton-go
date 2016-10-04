package triton

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hashicorp/errwrap"
)

type AccountsClient struct {
	*Client
}

// Accounts returns a client used for accessing functions pertaining
// to Account functionality in the Triton API.
func (client *Client) Accounts() *AccountsClient {
	return &AccountsClient{client}
}

type Account struct {
	ID               string    `json:"id"`
	Login            string    `json:"login"`
	Email            string    `json:"email"`
	CompanyName      string    `json:"companyName"`
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"lastName"`
	Address          string    `json:"address"`
	PostalCode       string    `json:"postalCode"`
	City             string    `json:"city"`
	State            string    `json:"state"`
	Country          string    `json:"country"`
	Phone            string    `json:"phone"`
	Created          time.Time `json:"created"`
	Updated          time.Time `json:"updated"`
	TritonCNSEnabled bool      `json:"triton_cns_enabled"`
}

type GetAccountInput struct{}

func (client *AccountsClient) GetAccount(input *GetAccountInput) (*Account, error) {
	respReader, err := client.executeRequest(http.MethodGet, "/my", nil)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing GetAccount request: {{err}}", err)
	}

	var result *Account
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding GetAccount response: {{err}}", err)
	}

	return result, nil
}
