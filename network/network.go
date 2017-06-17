package network

import "github.com/joyent/triton-go/client"

type NetworkService struct {
	client *client.Client
}

// Fabrics returns a Compute client used for accessing functions pertaining to
// Fabric functionality in the Triton API.
func (c *NetworkService) Fabrics() *FabricsClient {
	return &FabricsClient{c}
}

// Firewall returns a NetworkService client used for accessing functions pertaining to
// firewall functionality in the Triton API.
func (c *NetworkService) Firewall() *FirewallClient {
	return &FirewallClient{c}
}

// Networks returns a NetworkService client used for accessing functions pertaining to
// Network functionality in the Triton API.
func (c *NetworkService) Networks() *NetworksClient {
	return &NetworksClient{c}
}
