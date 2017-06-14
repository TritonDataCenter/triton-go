package identity

import "github.com/joyent/triton-go/client"

type Identity struct {
	client *client.Client
}

// Roles returns a Compute client used for accessing functions pertaining to
// Role functionality in the Triton API.
func (c *Identity) Roles() *RolesClient {
	return &RolesClient{c}
}
