package testutils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
)

const TestEnvVar = "TRITON_TEST"

type TestCase struct {
	Steps []Step
	State TritonStateBag
}

func AccTest(t *testing.T, c TestCase) {
	// We only run acceptance tests if an env var is set because they're
	// slow and generally require some outside configuration.
	if os.Getenv(TestEnvVar) == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set",
			TestEnvVar))
		return
	}

	// We require verbose mode so that the user knows what is going on.
	if !testing.Verbose() {
		t.Fatal("Acceptance tests must be run with the -v flag on tests")
		return
	}

	sdcURL := os.Getenv("SDC_URL")
	sdcAccount := os.Getenv("SDC_ACCOUNT")
	sdcKeyId := os.Getenv("SDC_KEY_ID")
	sdcKeyMaterial := os.Getenv("SDC_KEY_MATERIAL")
	mantaURL := os.Getenv("MANTA_URL")

	var prerollErrors []error
	if sdcURL == "" {
		prerollErrors = append(prerollErrors,
			errors.New("The SDC_URL environment variable must be set to run acceptance tests"))
	}
	if sdcAccount == "" {
		prerollErrors = append(prerollErrors,
			errors.New("The SDC_ACCOUNT environment variable must be set to run acceptance tests"))
	}
	if sdcKeyId == "" {
		prerollErrors = append(prerollErrors,
			errors.New("The SDC_KEY_ID environment variable must be set to run acceptance tests"))
	}
	if len(prerollErrors) > 0 {
		for _, err := range prerollErrors {
			t.Error(err)
		}
		t.FailNow()
	}

	var signer authentication.Signer
	var err error
	if sdcKeyMaterial != "" {
		log.Println("[INFO] Creating Triton Client with Private Key Signer...")
		signer, err = authentication.NewPrivateKeySigner(sdcKeyId, []byte(sdcKeyMaterial), sdcAccount)
		if err != nil {
			t.Fatalf("Error creating private key signer: %s", err)
		}
	} else {
		log.Println("[INFO] Creating Triton Client with SSH Key Signer...")
		signer, err = authentication.NewSSHAgentSigner(sdcKeyId, sdcAccount)
		if err != nil {
			t.Fatalf("Error creating SSH Agent signer: %s", err)
		}
	}

	// Old world... we spun up a universal client. This is pushed deeper into
	// the process within `testutils.StepClient`.
	//
	// client, err := NewClient(sdcURL, sdcAccount, signer)
	// if err != nil {
	//         t.Fatalf("Error creating Triton Client: %s", err)
	// }

	config := &triton.ClientConfig{
		TritonURL:   sdcURL,
		MantaURL:    mantaURL,
		AccountName: sdcAccount,
		Signers:     []authentication.Signer{signer},
	}

	state := &basicTritonStateBag{
		TritonConfig: config,
	}

	runner := &basicRunner{
		Steps: c.Steps,
	}

	runner.Run(state)

	if errs := state.ErrorsOrNil(); errs != nil {
		for _, err := range errs {
			t.Error(err)
		}
		t.Fatal("\n\nThere may be dangling resources in your Triton account!")
	}
}
