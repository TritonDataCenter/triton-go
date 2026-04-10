//
// Copyright 2020 Joyent, Inc.
// Copyright 2026 Edgecast Cloud LLC.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package triton

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// Version represents main version number of the current release
// of the Triton-go SDK.
const Version = "2.0.0"

// Prerelease adds a pre-release marker to the version.
//
// If this is "" (empty string) then it means that it is a final release.
// Otherwise, this is a pre-release such as "dev" (in development), "beta",
// "rc1", etc.
var Prerelease = "pre4"

// UserAgent returns a Triton-go characteristic string that allows the
// network protocol peers to identify the version, release and runtime
// of the Triton-go client from which the requests originate.
func UserAgent() string {
	if Prerelease != "" {
		return fmt.Sprintf("triton-go/%s-%s (%s-%s; %s)", Version, Prerelease,
			runtime.GOARCH, runtime.GOOS, runtime.Version())
	}

	return fmt.Sprintf("triton-go/%s (%s-%s; %s)", Version, runtime.GOARCH,
		runtime.GOOS, runtime.Version())
}

// CloudAPIMajorVersion specifies the CloudAPI version compatibility
// for current release of the Triton-go SDK.
const CloudAPIMajorVersion = "8"

// VCSInfo contains version control metadata embedded by the Go toolchain
// at build time.
type VCSInfo struct {
	Revision string
	Time     string
	Modified bool
}

// BuildInfo returns VCS metadata embedded by the Go toolchain. If the
// binary was not built with module support or VCS info is unavailable,
// the returned fields will be empty/false.
func BuildInfo() VCSInfo {
	var info VCSInfo
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return info
	}
	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			info.Revision = s.Value
		case "vcs.time":
			info.Time = s.Value
		case "vcs.modified":
			info.Modified = s.Value == "true"
		}
	}
	return info
}
