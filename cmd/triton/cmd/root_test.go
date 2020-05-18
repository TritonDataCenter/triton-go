package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/joyent/triton-go/v2/testutils"
)

func TestListDataCenters_Cmd(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{
			&testutils.StepAssertFunc{
				AssertFunc: func(state testutils.TritonStateBag) error {
					cmd := exec.Command("../triton", "datacenters")
					out := bytes.NewBuffer([]byte{})

					cmd.Stdout = out
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						return fmt.Errorf("%v", err)
					}
					re, err := regexp.Compile(`(?i)url`)
					if err != nil {
						return fmt.Errorf("Error compiling Regexp: %v", err)
					}

					if !re.MatchString(out.String()) {
						return fmt.Errorf("Unexpected command stdout:\n%s", out.String())
					}

					t.Logf("\n%s =>\n%s", strings.Join(cmd.Args[:], " "), out.String())
					return nil
				},
			},
		},
	})
}

func TestListInstances_Cmd(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{
			&testutils.StepAssertFunc{
				AssertFunc: func(state testutils.TritonStateBag) error {
					cmd := exec.Command("../triton", "instances", "list")
					out := bytes.NewBuffer([]byte{})

					cmd.Stdout = out
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						return fmt.Errorf("%v", err)
					}
					re, err := regexp.Compile(`(?i)shortid`)
					if err != nil {
						return fmt.Errorf("Error compiling Regexp: %v", err)
					}

					if !re.MatchString(out.String()) {
						return fmt.Errorf("Unexpected command stdout:\n%s", out.String())
					}
					t.Logf("\n%s =>\n%s", strings.Join(cmd.Args[:], " "), out.String())
					return nil
				},
			},
		},
	})
}

func TestListPackages_Cmd(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{
			&testutils.StepAssertFunc{
				AssertFunc: func(state testutils.TritonStateBag) error {
					cmd := exec.Command("../triton", "package", "list", "--memory", "128")
					out := bytes.NewBuffer([]byte{})

					cmd.Stdout = out
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						return fmt.Errorf("%v", err)
					}
					re, err := regexp.Compile(`(?i)shortid`)
					if err != nil {
						return fmt.Errorf("Error compiling Regexp: %v", err)
					}

					if !re.MatchString(out.String()) {
						return fmt.Errorf("Unexpected command stdout:\n%s", out.String())
					}
					t.Logf("\n%s =>\n%s", strings.Join(cmd.Args[:], " "), out.String())
					return nil
				},
			},
		},
	})
}

func TestGetPackageByName_Cmd(t *testing.T) {
	testutils.AccTest(t, testutils.TestCase{
		Steps: []testutils.Step{
			&testutils.StepAssertFunc{
				AssertFunc: func(state testutils.TritonStateBag) error {
					cmd := exec.Command("../triton", "package", "get", "--name", "sample-128M")
					out := bytes.NewBuffer([]byte{})

					cmd.Stdout = out
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						return fmt.Errorf("%v", err)
					}
					re, err := regexp.Compile(`(?i)shortid`)
					if err != nil {
						return fmt.Errorf("Error compiling Regexp: %v", err)
					}

					t.Logf("\n%s =>\n%s", strings.Join(cmd.Args[:], " "), out.String())
					if !re.MatchString(out.String()) {
						return fmt.Errorf("Unexpected command stdout:\n%s", out.String())
					}
					t.Logf("\n%s =>\n%s", strings.Join(cmd.Args[:], " "), out.String())
					return nil
				},
			},
		},
	})
}
