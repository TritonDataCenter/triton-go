package networking

import "github.com/joyent/triton-go/client"

type Networking struct {
	client *client.Client
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

// Networks returns a Compute client used for accessing functions pertaining to
// Network functionality in the Triton API.
func (c *Compute) Networks() *NetworksClient {
	return &NetworksClient{c}
}
