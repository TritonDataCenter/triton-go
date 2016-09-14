package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jen20/triton-go"
	"github.com/jen20/triton-go/authentication"
	"log"
)

func main() {
	sshKeySigner, err := authentication.NewSSHAgentSigner(
		"1b:bc:29:48:89:af:72:63:f0:83:b8:11:b6:4d:ff:3f", "hashicorp")
	if err != nil {
		log.Fatalf("NewSSHAgentSigner: %s", err)
	}

	client, err := triton.NewClient("https://us-sw-1.api.joyent.com/", "hashicorp", sshKeySigner)
	if err != nil {
		log.Fatalf("NewClient: %s", err)
	}

	dcs, err := client.Datacenters().ListDataCenters(&triton.ListDataCentersInput{})
	if err != nil {
		log.Fatalf("ListDatacenters(): %s", err)
	}

	dc0, err := client.Datacenters().GetDataCenter(&triton.GetDataCenterInput{dcs[0].Name})
	if err != nil {
		log.Fatalf("GetDatacenter(): %s", err)
	}

	spew.Dump(dcs)
	spew.Dump(dc0)
}
