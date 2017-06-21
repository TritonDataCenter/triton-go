package triton

import (
	"github.com/joyent/triton-go/authentication"
)

type ClientConfig struct {
	Endpoint    string
	AccountName string
	Signers     []authentication.Signer
}
