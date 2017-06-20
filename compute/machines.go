package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/joyent/triton-go/client"
)

type MachinesClient struct {
	client *client.Client
}

const (
	machineCNSTagDisable    = "triton.cns.disable"
	machineCNSTagReversePTR = "triton.cns.reverse_ptr"
	machineCNSTagServices   = "triton.cns.services"
)

// MachineCNS is a container for the CNS-specific attributes.  In the API these
// values are embedded within a Machine's Tags attribute, however they are
// exposed to the caller as their native types.
type MachineCNS struct {
	Disable    *bool
	ReversePTR *string
	Services   []string
}

type Machine struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Type            string                 `json:"type"`
	Brand           string                 `json:"brand"`
	State           string                 `json:"state"`
	Image           string                 `json:"image"`
	Memory          int                    `json:"memory"`
	Disk            int                    `json:"disk"`
	Metadata        map[string]string      `json:"metadata"`
	Tags            map[string]interface{} `json:"tags"`
	Created         time.Time              `json:"created"`
	Updated         time.Time              `json:"updated"`
	Docker          bool                   `json:"docker"`
	IPs             []string               `json:"ips"`
	Networks        []string               `json:"networks"`
	PrimaryIP       string                 `json:"primaryIp"`
	FirewallEnabled bool                   `json:"firewall_enabled"`
	ComputeNode     string                 `json:"compute_node"`
	Package         string                 `json:"package"`
	DomainNames     []string               `json:"dns_names"`
	CNS             MachineCNS
}

// _Machine is a private facade over Machine that handles the necessary API
// overrides from VMAPI's machine endpoint(s).
type _Machine struct {
	Machine
	Tags map[string]interface{} `json:"tags"`
}

type NIC struct {
	IP      string `json:"ip"`
	MAC     string `json:"mac"`
	Primary bool   `json:"primary"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	State   string `json:"state"`
	Network string `json:"network"`
}

type GetMachineInput struct {
	ID string
}

func (gmi *GetMachineInput) Validate() error {
	if gmi.ID == "" {
		return fmt.Errorf("machine ID can not be empty")
	}

	return nil
}

func (c *MachinesClient) GetMachine(ctx context.Context, input *GetMachineInput) (*Machine, error) {
	if err := input.Validate(); err != nil {
		return nil, errwrap.Wrapf("unable to get machine: {{err}}", err)
	}

	path := fmt.Sprintf("/%s/machines/%s", c.client.AccountName, input.ID)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   path,
	}
	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if response != nil {
		defer response.Body.Close()
	}
	if response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusGone {
		return nil, &TritonError{
			StatusCode: response.StatusCode,
			Code:       "ResourceNotFound",
		}
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing GetMachine request: {{err}}",
			c.client.DecodeError(response.StatusCode, response.Body))
	}

	var result *_Machine
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding GetMachine response: {{err}}", err)
	}

	native, err := result.toNative()
	if err != nil {
		return nil, errwrap.Wrapf("unable to convert API response for machines to native type: {{err}}", err)
	}

	return native, nil
}

type ListMachinesInput struct{}

func (c *MachinesClient) ListMachines(ctx context.Context, _ *ListMachinesInput) ([]*Machine, error) {
	path := fmt.Sprintf("/%s/machines", c.client.AccountName)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   path,
	}
	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if response != nil {
		defer response.Body.Close()
	}
	if response.StatusCode == http.StatusNotFound {
		return nil, &TritonError{
			StatusCode: response.StatusCode,
			Code:       "ResourceNotFound",
		}
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing ListMachines request: {{err}}",
			c.client.DecodeError(response.StatusCode, response.Body))
	}

	var results []*_Machine
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&results); err != nil {
		return nil, errwrap.Wrapf("Error decoding ListMachines response: {{err}}", err)
	}

	machines := make([]*Machine, 0, len(results))
	for _, machineAPI := range results {
		native, err := machineAPI.toNative()
		if err != nil {
			return nil, errwrap.Wrapf("unable to convert API response for machines to native type: {{err}}", err)
		}
		machines = append(machines, native)
	}
	return machines, nil
}

type CreateMachineInput struct {
	Name            string
	Package         string
	Image           string
	Networks        []string
	LocalityStrict  bool
	LocalityNear    []string
	LocalityFar     []string
	Metadata        map[string]string
	Tags            map[string]string
	FirewallEnabled bool
	CNS             MachineCNS
}

func (input *CreateMachineInput) toAPI() map[string]interface{} {
	const numExtraParams = 8
	result := make(map[string]interface{}, numExtraParams+len(input.Metadata)+len(input.Tags))

	result["firewall_enabled"] = input.FirewallEnabled

	if input.Name != "" {
		result["name"] = input.Name
	}

	if input.Package != "" {
		result["package"] = input.Package
	}

	if input.Image != "" {
		result["image"] = input.Image
	}

	if len(input.Networks) > 0 {
		result["networks"] = input.Networks
	}

	locality := struct {
		Strict bool     `json:"strict"`
		Near   []string `json:"near,omitempty"`
		Far    []string `json:"far,omitempty"`
	}{
		Strict: input.LocalityStrict,
		Near:   input.LocalityNear,
		Far:    input.LocalityFar,
	}
	result["locality"] = locality
	for key, value := range input.Tags {
		result[fmt.Sprintf("tag.%s", key)] = value
	}

	// Deliberately clobber any user-specified Tags with the attributes from the
	// CNS struct.
	input.CNS.toTags(result)

	for key, value := range input.Metadata {
		result[fmt.Sprintf("metadata.%s", key)] = value
	}

	return result
}

func (c *MachinesClient) CreateMachine(ctx context.Context, input *CreateMachineInput) (*Machine, error) {
	path := fmt.Sprintf("/%s/machines", c.client.AccountName)
	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   path,
		Body:   input.toAPI(),
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing CreateMachine request: {{err}}", err)
	}

	var result *Machine
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding CreateMachine response: {{err}}", err)
	}

	return result, nil
}

type DeleteMachineInput struct {
	ID string
}

func (c *MachinesClient) DeleteMachine(ctx context.Context, input *DeleteMachineInput) error {
	path := fmt.Sprintf("/%s/machines/%s", c.client.AccountName, input.ID)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   path,
	}
	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if response.Body != nil {
		defer response.Body.Close()
	}
	if response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusGone {
		return nil
	}
	if err != nil {
		return errwrap.Wrapf("Error executing DeleteMachine request: {{err}}",
			c.client.DecodeError(response.StatusCode, response.Body))
	}

	return nil
}

type DeleteMachineTagsInput struct {
	ID string
}

func (c *MachinesClient) DeleteMachineTags(ctx context.Context, input *DeleteMachineTagsInput) error {
	path := fmt.Sprintf("/%s/machines/%s/tags", c.client.AccountName, input.ID)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   path,
	}
	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if response.Body != nil {
		defer response.Body.Close()
	}
	if response.StatusCode == http.StatusNotFound {
		return nil
	}
	if err != nil {
		return errwrap.Wrapf("Error executing DeleteMachineTags request: {{err}}",
			c.client.DecodeError(response.StatusCode, response.Body))
	}

	return nil
}

type DeleteMachineTagInput struct {
	ID  string
	Key string
}

func (c *MachinesClient) DeleteMachineTag(ctx context.Context, input *DeleteMachineTagInput) error {
	path := fmt.Sprintf("/%s/machines/%s/tags/%s", c.client.AccountName, input.ID, input.Key)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   path,
	}
	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if response.Body != nil {
		defer response.Body.Close()
	}
	if response.StatusCode == http.StatusNotFound {
		return nil
	}
	if err != nil {
		return errwrap.Wrapf("Error executing DeleteMachineTag request: {{err}}",
			c.client.DecodeError(response.StatusCode, response.Body))
	}

	return nil
}

type RenameMachineInput struct {
	ID   string
	Name string
}

func (c *MachinesClient) RenameMachine(ctx context.Context, input *RenameMachineInput) error {
	path := fmt.Sprintf("/%s/machines/%s", c.client.AccountName, input.ID)

	params := &url.Values{}
	params.Set("action", "rename")
	params.Set("name", input.Name)

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   path,
		Query:  params,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing RenameMachine request: {{err}}", err)
	}

	return nil
}

type ReplaceMachineTagsInput struct {
	ID   string
	Tags map[string]string
}

func (c *MachinesClient) ReplaceMachineTags(ctx context.Context, input *ReplaceMachineTagsInput) error {
	path := fmt.Sprintf("/%s/machines/%s/tags", c.client.AccountName, input.ID)
	reqInputs := client.RequestInput{
		Method: http.MethodPut,
		Path:   path,
		Body:   input.Tags,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing ReplaceMachineTags request: {{err}}", err)
	}

	return nil
}

type AddMachineTagsInput struct {
	ID   string
	Tags map[string]string
}

func (c *MachinesClient) AddMachineTags(ctx context.Context, input *AddMachineTagsInput) error {
	path := fmt.Sprintf("/%s/machines/%s/tags", c.client.AccountName, input.ID)
	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   path,
		Body:   input.Tags,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing AddMachineTags request: {{err}}", err)
	}

	return nil
}

type GetMachineTagInput struct {
	ID  string
	Key string
}

func (c *MachinesClient) GetMachineTag(ctx context.Context, input *GetMachineTagInput) (string, error) {
	path := fmt.Sprintf("/%s/machines/%s/tags/%s", c.client.AccountName, input.ID, input.Key)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   path,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return "", errwrap.Wrapf("Error executing GetMachineTag request: {{err}}", err)
	}

	var result string
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return "", errwrap.Wrapf("Error decoding GetMachineTag response: {{err}}", err)
	}

	return result, nil
}

type ListMachineTagsInput struct {
	ID string
}

func (c *MachinesClient) ListMachineTags(ctx context.Context, input *ListMachineTagsInput) (map[string]interface{}, error) {
	path := fmt.Sprintf("/%s/machines/%s/tags", c.client.AccountName, input.ID)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   path,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing ListMachineTags request: {{err}}", err)
	}

	var result map[string]interface{}
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding ListMachineTags response: {{err}}", err)
	}

	_, tags := machineTagsExtractMeta(result)
	return tags, nil
}

type UpdateMachineMetadataInput struct {
	ID       string
	Metadata map[string]string
}

func (c *MachinesClient) UpdateMachineMetadata(ctx context.Context, input *UpdateMachineMetadataInput) (map[string]string, error) {
	path := fmt.Sprintf("/%s/machines/%s/tags", c.client.AccountName, input.ID)
	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   path,
		Body:   input.Metadata,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing UpdateMachineMetadata request: {{err}}", err)
	}

	var result map[string]string
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding UpdateMachineMetadata response: {{err}}", err)
	}

	return result, nil
}

type ResizeMachineInput struct {
	ID      string
	Package string
}

func (c *MachinesClient) ResizeMachine(ctx context.Context, input *ResizeMachineInput) error {
	path := fmt.Sprintf("/%s/machines/%s", c.client.AccountName, input.ID)

	params := &url.Values{}
	params.Set("action", "resize")
	params.Set("package", input.Package)

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   path,
		Query:  params,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing ResizeMachine request: {{err}}", err)
	}

	return nil
}

type EnableMachineFirewallInput struct {
	ID string
}

func (c *MachinesClient) EnableMachineFirewall(ctx context.Context, input *EnableMachineFirewallInput) error {
	path := fmt.Sprintf("/%s/machines/%s", c.client.AccountName, input.ID)

	params := &url.Values{}
	params.Set("action", "enable_firewall")

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   path,
		Query:  params,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing EnableMachineFirewall request: {{err}}", err)
	}

	return nil
}

type DisableMachineFirewallInput struct {
	ID string
}

func (c *MachinesClient) DisableMachineFirewall(ctx context.Context, input *DisableMachineFirewallInput) error {
	path := fmt.Sprintf("/%s/machines/%s", c.client.AccountName, input.ID)

	params := &url.Values{}
	params.Set("action", "disable_firewall")

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   path,
		Query:  params,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing DisableMachineFirewall request: {{err}}", err)
	}

	return nil
}

type ListNICsInput struct {
	MachineID string
}

func (c *MachinesClient) ListNICs(ctx context.Context, input *ListNICsInput) ([]*NIC, error) {
	path := fmt.Sprintf("/%s/machines/%s/nics", c.client.AccountName, input.MachineID)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   path,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing ListNICs request: {{err}}", err)
	}

	var result []*NIC
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding ListNICs response: {{err}}", err)
	}

	return result, nil
}

type GetNICInput struct {
	MachineID string
	MAC       string
}

func (c *MachinesClient) GetNIC(ctx context.Context, input *GetNICInput) (*NIC, error) {
	mac := strings.Replace(input.MAC, ":", "", -1)
	path := fmt.Sprintf("/%s/machines/%s/nics/%s", c.client.AccountName, input.MachineID, mac)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   path,
	}
	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if response != nil {
		defer response.Body.Close()
	}
	switch response.StatusCode {
	case http.StatusNotFound:
		return nil, &TritonError{
			StatusCode: response.StatusCode,
			Code:       "ResourceNotFound",
		}
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing GetNIC request: {{err}}", err)
	}

	var result *NIC
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding ListNICs response: {{err}}", err)
	}

	return result, nil
}

type AddNICInput struct {
	MachineID string `json:"-"`
	Network   string `json:"network"`
}

// AddNIC asynchronously adds a NIC to a given machine.  If a NIC for a given
// network already exists, a ResourceFound error will be returned.  The status
// of the addition of a NIC can be polled by calling GetNIC()'s and testing NIC
// until its state is set to "running".  Only one NIC per network may exist.
// Warning: this operation causes the machine to restart.
func (c *MachinesClient) AddNIC(ctx context.Context, input *AddNICInput) (*NIC, error) {
	path := fmt.Sprintf("/%s/machines/%s/nics", c.client.AccountName, input.MachineID)
	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   path,
		Body:   input,
	}
	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if response != nil {
		defer response.Body.Close()
	}
	switch response.StatusCode {
	case http.StatusFound:
		return nil, &TritonError{
			StatusCode: response.StatusCode,
			Code:       "ResourceFound",
			Message:    response.Header.Get("Location"),
		}
	}
	if err != nil {
		return nil, errwrap.Wrapf("Error executing AddNIC request: {{err}}", err)
	}

	var result *NIC
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&result); err != nil {
		return nil, errwrap.Wrapf("Error decoding AddNIC response: {{err}}", err)
	}

	return result, nil
}

type RemoveNICInput struct {
	MachineID string
	MAC       string
}

// RemoveNIC removes a given NIC from a machine asynchronously.  The status of
// the removal can be polled via GetNIC().  When GetNIC() returns a 404, the NIC
// has been removed from the instance.  Warning: this operation causes the
// machine to restart.
func (c *MachinesClient) RemoveNIC(ctx context.Context, input *RemoveNICInput) error {
	mac := strings.Replace(input.MAC, ":", "", -1)
	path := fmt.Sprintf("/%s/machines/%s/nics/%s", c.client.AccountName, input.MachineID, mac)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   path,
	}
	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if response != nil {
		defer response.Body.Close()
	}
	switch response.StatusCode {
	case http.StatusNotFound:
		return &TritonError{
			StatusCode: response.StatusCode,
			Code:       "ResourceNotFound",
		}
	}
	if err != nil {
		return errwrap.Wrapf("Error executing RemoveNIC request: {{err}}", err)
	}

	return nil
}

type StopMachineInput struct {
	MachineID string
}

func (c *MachinesClient) StopMachine(ctx context.Context, input *StopMachineInput) error {
	path := fmt.Sprintf("/%s/machines/%s", c.client.AccountName, input.MachineID)

	params := &url.Values{}
	params.Set("action", "stop")

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   path,
		Query:  params,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing StopMachine request: {{err}}", err)
	}

	return nil
}

type StartMachineInput struct {
	MachineID string
}

func (c *MachinesClient) StartMachine(ctx context.Context, input *StartMachineInput) error {
	path := fmt.Sprintf("/%s/machines/%s", c.client.AccountName, input.MachineID)

	params := &url.Values{}
	params.Set("action", "start")

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   path,
		Query:  params,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return errwrap.Wrapf("Error executing StartMachine request: {{err}}", err)
	}

	return nil
}

var reservedMachineCNSTags = map[string]struct{}{
	machineCNSTagDisable:    {},
	machineCNSTagReversePTR: {},
	machineCNSTagServices:   {},
}

// machineTagsExtractMeta() extracts all of the misc parameters from Tags and
// returns a clean CNS and Tags struct.
func machineTagsExtractMeta(tags map[string]interface{}) (MachineCNS, map[string]interface{}) {
	nativeCNS := MachineCNS{}
	nativeTags := make(map[string]interface{}, len(tags))
	for k, raw := range tags {
		if _, found := reservedMachineCNSTags[k]; found {
			switch k {
			case machineCNSTagDisable:
				b := raw.(bool)
				nativeCNS.Disable = &b
			case machineCNSTagReversePTR:
				s := raw.(string)
				nativeCNS.ReversePTR = &s
			case machineCNSTagServices:
				nativeCNS.Services = strings.Split(raw.(string), ",")
			default:
				// TODO(seanc@): should assert, logic fail
			}
		} else {
			nativeTags[k] = raw
		}
	}

	return nativeCNS, nativeTags
}

// toNative() exports a given _Machine (API representation) to its native object
// format.
func (api *_Machine) toNative() (*Machine, error) {
	m := Machine(api.Machine)
	m.CNS, m.Tags = machineTagsExtractMeta(api.Tags)
	return &m, nil
}

// toTags() injects its state information into a Tags map suitable for use to
// submit an API call to the vmapi machine endpoint
func (mcns *MachineCNS) toTags(m map[string]interface{}) {
	if mcns.Disable != nil {
		s := fmt.Sprintf("%t", mcns.Disable)
		m[machineCNSTagDisable] = &s
	}

	if mcns.ReversePTR != nil {
		m[machineCNSTagReversePTR] = &mcns.ReversePTR
	}

	if len(mcns.Services) > 0 {
		m[machineCNSTagServices] = strings.Join(mcns.Services, ",")
	}
}
