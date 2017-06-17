package compute

import (
	"github.com/joyent/triton-go/client"
)

type Compute struct {
	Client *client.Client
}

// Datacenters returns a Compute client used for accessing functions pertaining
// to DataCenter functionality in the Triton API.
func (c *Compute) Datacenters() *DataCentersClient {
	return &DataCentersClient{c.Client}
}

// Images returns a Compute client used for accessing functions pertaining to
// Images functionality in the Triton API.
func (c *Compute) Images() *ImagesClient {
	return &ImagesClient{c.Client}
}

// Machines returns a Compute client used for accessing functions pertaining to
// machine functionality in the Triton API.
func (c *Compute) Machines() *MachinesClient {
	return &MachinesClient{c.Client}
}

// Packages returns a Compute client used for accessing functions pertaining to
// Packages functionality in the Triton API.
func (c *Compute) Packages() *PackagesClient {
	return &PackagesClient{c.Client}
}

// Services returns a Compute client used for accessing functions pertaining to
// Services functionality in the Triton API.
func (c *Compute) Services() *ServicesClient {
	return &ServicesClient{c.Client}
}
