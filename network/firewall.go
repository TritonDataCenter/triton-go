package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/errwrap"
)

type FirewallClient struct {
	*NetworkService
}

// FirewallRule represents a firewall rule
type FirewallRule struct {
	// ID is a unique identifier for this rule
	ID string `json:"id"`

	// Enabled indicates if the rule is enabled
	Enabled bool `json:"enabled"`

	// Rule is the firewall rule text
	Rule string `json:"rule"`

	// Global indicates if the rule is global. Optional.
	Global bool `json:"global"`

	// Description is a human-readable description for the rule. Optional
	Description string `json:"description"`
}

type ListFirewallRulesInput struct{}

func (c *FirewallClient) ListFirewallRules(ctx context.Context, _ *ListFirewallRulesInput) ([]*FirewallRule, error) {
	path := fmt.Sprintf("/%s/fwrules", c.client.AccountName)
	respReader, err := c.executeRequest(ctx, http.MethodGet, path, nil)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing ListFirewallRules request: {{err}}", err)
	}

	var result []*FirewallRule
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding ListFirewallRules response: {{err}}", err)
	}

	return result, nil
}

type GetFirewallRuleInput struct {
	ID string
}

func (c *FirewallClient) GetFirewallRule(ctx context.Context, input *GetFirewallRuleInput) (*FirewallRule, error) {
	path := fmt.Sprintf("/%s/fwrules/%s", c.client.AccountName, input.ID)
	respReader, err := c.executeRequest(ctx, http.MethodGet, path, nil)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing GetFirewallRule request: {{err}}", err)
	}

	var result *FirewallRule
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding GetFirewallRule response: {{err}}", err)
	}

	return result, nil
}

type CreateFirewallRuleInput struct {
	Enabled     bool   `json:"enabled"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

func (c *FirewallClient) CreateFirewallRule(ctx context.Context, input *CreateFirewallRuleInput) (*FirewallRule, error) {
	path := fmt.Sprintf("/%s/fwrules", c.client.AccountName)
	respReader, err := c.executeRequest(ctx, http.MethodPost, path, input)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing CreateFirewallRule request: {{err}}", err)
	}

	var result *FirewallRule
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding CreateFirewallRule response: {{err}}", err)
	}

	return result, nil
}

type UpdateFirewallRuleInput struct {
	ID          string `json:"-"`
	Enabled     bool   `json:"enabled"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

func (c *FirewallClient) UpdateFirewallRule(ctx context.Context, input *UpdateFirewallRuleInput) (*FirewallRule, error) {
	path := fmt.Sprintf("/%s/fwrules/%s", c.client.AccountName, input.ID)
	respReader, err := c.executeRequest(ctx, http.MethodPost, path, input)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing UpdateFirewallRule request: {{err}}", err)
	}

	var result *FirewallRule
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding UpdateFirewallRule response: {{err}}", err)
	}

	return result, nil
}

type EnableFirewallRuleInput struct {
	ID string `json:"-"`
}

func (c *FirewallClient) EnableFirewallRule(ctx context.Context, input *EnableFirewallRuleInput) (*FirewallRule, error) {
	path := fmt.Sprintf("/%s/fwrules/%s/enable", c.client.AccountName, input.ID)
	respReader, err := c.executeRequest(ctx, http.MethodPost, path, input)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing EnableFirewallRule request: {{err}}", err)
	}

	var result *FirewallRule
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding EnableFirewallRule response: {{err}}", err)
	}

	return result, nil
}

type DisableFirewallRuleInput struct {
	ID string `json:"-"`
}

func (c *FirewallClient) DisableFirewallRule(ctx context.Context, input *DisableFirewallRuleInput) (*FirewallRule, error) {
	path := fmt.Sprintf("/%s/fwrules/%s/disable", c.client.AccountName, input.ID)
	respReader, err := c.executeRequest(ctx, http.MethodPost, path, input)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing DisableFirewallRule request: {{err}}", err)
	}

	var result *FirewallRule
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding DisableFirewallRule response: {{err}}", err)
	}

	return result, nil
}

type DeleteFirewallRuleInput struct {
	ID string
}

func (c *FirewallClient) DeleteFirewallRule(ctx context.Context, input *DeleteFirewallRuleInput) error {
	path := fmt.Sprintf("/%s/fwrules/%s", c.client.AccountName, input.ID)
	respReader, err := c.executeRequest(ctx, http.MethodDelete, path, nil)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing DeleteFirewallRule request: {{err}}", err)
	}

	return nil
}

type ListMachineFirewallRulesInput struct {
	MachineID string
}

func (c *FirewallClient) ListMachineFirewallRules(ctx context.Context, input *ListMachineFirewallRulesInput) ([]*FirewallRule, error) {
	path := fmt.Sprintf("/%s/machines/%s/firewallrules", c.client.AccountName, input.MachineID)
	respReader, err := c.executeRequest(ctx, http.MethodGet, path, nil)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing ListMachineFirewallRules request: {{err}}", err)
	}

	var result []*FirewallRule
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding ListFirewallRules response: {{err}}", err)
	}

	return result, nil
}
