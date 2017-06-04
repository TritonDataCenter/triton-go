package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jen20/manta-go"
	"github.com/jen20/manta-go/authentication"
)

func main() {
	sshKeySigner, err := authentication.NewSSHAgentSigner(
		"fd:9e:9a:9c:28:99:57:05:18:9f:b6:44:6b:cc:fd:3a", "tritongo")
	if err != nil {
		log.Fatalf("NewSSHAgentSigner: %s", err)
	}

	client, err := manta.NewClient(&manta.ClientOptions{
		Endpoint:    "https://us-east.manta.joyent.com/",
		AccountName: "tritongo",
		Signers:     []authentication.Signer{sshKeySigner},
	})
	if err != nil {
		log.Fatalf("NewClient: %s", err)
	}

	output, err := client.GetObject(&manta.GetObjectInput{
		ObjectPath: "tempstate.tfstate",
	})
	if err != nil {
		log.Fatalf("GetObject(): %s", err)
	}

	defer output.ObjectReader.Close()
	body, err := ioutil.ReadAll(output.ObjectReader)
	if err != nil {
		log.Fatalf("Reading Object: %s", err)
	}

	fmt.Printf("Content-Length: %d\n", output.ContentLength)
	fmt.Printf("Content-MD5: %s\n", output.ContentMD5)
	fmt.Printf("Content-Type: %s\n", output.ContentType)
	fmt.Printf("ETag: %s\n", output.ETag)
	fmt.Printf("Date-Modified: %s\n", output.LastModified.String())
	fmt.Printf("Object:\n\n%s", string(body))
}
