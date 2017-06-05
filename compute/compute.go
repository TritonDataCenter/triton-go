package compute

import "github.com/joyent/triton-go/client"

type Compute struct {
	*client.Client
}

type Compute interface {
	Accounts() *AccountsClient
	Config() *ConfigClient
	Datacenters() *DataCentersClient
	Fabrics() *FabricsClient
	Firewall() *FirewallClient
	Images() *ImagesClient
	Keys() *KeysClient
	Machines() *MachinesClient
	Networks() *NetworksClient
	Packages() *PackagesClient
	Roles() *RolesClient
	Services() *ServicesClient
}

// Accounts returns a Compute client used for accessing functions pertaining to
// Account functionality in the Triton API.
func (c *Compute) Accounts() *AccountsClient {
	return &AccountsClient{c}
}

// Config returns a c used for accessing functions pertaining
// to Config functionality in the Triton API.
func (c *Compute) Config() *ConfigClient {
	return &ConfigClient{c}
}

// Datacenters returns a Compute client used for accessing functions pertaining
// to DataCenter functionality in the Triton API.
func (c *Compute) Datacenters() *DataCentersClient {
	return &DataCentersClient{c}
}

// Fabrics returns a Compute client used for accessing functions pertaining to
// Fabric functionality in the Triton API.
func (c *Compute) Fabrics() *FabricsClient {
	return &FabricsClient{c}
}

// Firewall returns a Compute client used for accessing functions pertaining to
// firewall functionality in the Triton API.
func (c *Compute) Firewall() *FirewallClient {
	return &FirewallClient{c}
}

// Images returns a Compute client used for accessing functions pertaining to
// Images functionality in the Triton API.
func (c *Compute) Images() *ImagesClient {
	return &ImagesClient{c}
}

// Keys returns a Compute client used for accessing functions pertaining to SSH
// key functionality in the Triton API.
func (c *Compute) Keys() *KeysClient {
	return &KeysClient{c}
}

// Machines returns a Compute client used for accessing functions pertaining to
// machine functionality in the Triton API.
func (c *Compute) Machines() *MachinesClient {
	return &MachinesClient{c}
}

// Networks returns a Compute client used for accessing functions pertaining to
// Network functionality in the Triton API.
func (c *Compute) Networks() *NetworksClient {
	return &NetworksClient{c}
}

// Packages returns a Compute client used for accessing functions pertaining to
// Packages functionality in the Triton API.
func (c *Compute) Packages() *PackagesClient {
	return &PackagesClient{c}
}

// Roles returns a Compute client used for accessing functions pertaining to
// Role functionality in the Triton API.
func (c *Compute) Roles() *RolesClient {
	return &RolesClient{c}
}

// Services returns a Compute client used for accessing functions pertaining to
// Services functionality in the Triton API.
func (c *Client) Services() *ServicesClient {
	return &ServicesClient{c}
}
