package triton

import (
	"fmt"
	"runtime"
)

// The main version number of the current released Triton-go SDK.
const Version = "0.9.0"

// A pre-release marker for the version. If this is "" (empty string)
// then it means that it is a final release. Otherwise, this is a pre-release
// such as "dev" (in development), "beta", "rc1", etc.
var Prerelease = ""

func UserAgent() string {
	if Prerelease != "" {
		return fmt.Sprintf("triton-go/%s-%s (%s-%s; %s)", Version, Prerelease, runtime.GOARCH, runtime.GOOS, runtime.Version())
	} else {
		return fmt.Sprintf("triton-go/%s (%s-%s; %s)", Version, runtime.GOARCH, runtime.GOOS, runtime.Version())
	}
}

const CloudAPIMajorVersion = "8"
