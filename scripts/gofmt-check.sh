#!/bin/bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
#

#
# Copyright 2019 Joyent, Inc.
#

#
# Determine the list of Go source files to check for gofmt problems.  Note that
# we do not check the formatting of files in the vendor/ tree.
#
if ! files=$(git ls-files '*.go' ':!:vendor/*') || [[ -z $files ]]; then
	printf 'ERROR: could not find go file list\n' >&2
	exit 1
fi

if ! diff=$(gofmt -d $files); then
	printf 'ERROR: could not run "gofmt -d"\n' >&2
	exit 1
fi

if [[ -z $diff ]]; then
	printf 'gofmt check ok\n'
	exit 0
fi

#
# Report the detected formatting issues and exit non-zero so that "make check"
# will fail.
#
printf -- '---- GOFMT ERRORS -------\n'
printf -- '%s\n' "$diff"
exit 2
