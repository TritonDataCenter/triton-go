//
//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package docs

import (
	"fmt"
	"os"

	"path"
	"path/filepath"
	"strings"
	time "time"

	"github.com/joyent/triton-go/cmd/internal/command"
	"github.com/joyent/triton-go/cmd/internal/config"
	"github.com/joyent/triton-go/cmd/internal/console_writer"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

const gendocFrontmatterTemplate = `---
date: %s
title: "%s"
slug: %s
url: %s
---
`

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:   "doc",
		Short: "Generates and installs triton cli documentation in markdown",
		Long: `Generate Markdown documentation for the Triton CLI.

It creates one Markdown file per command `,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cons := console_writer.GetTerminal()

			location := viper.GetString(config.KeyMarkdownDirectory)
			if location == "" {
				location = "docs/md"
			}
			if _, err := os.Stat(location); os.IsNotExist(err) {
				os.Mkdir(location, 0777)
			}

			now := time.Now().UTC().Format(time.RFC3339)
			prepender := func(filename string) string {
				name := filepath.Base(filename)
				base := strings.TrimSuffix(name, path.Ext(name))
				url := "/commands/" + strings.ToLower(base) + "/"
				return fmt.Sprintf(gendocFrontmatterTemplate, now, strings.Replace(base, "_", " ", -1), base, url)
			}

			linkHandler := func(name string) string {
				base := strings.TrimSuffix(name, path.Ext(name))
				return "/commands/" + strings.ToLower(base) + "/"
			}

			cons.Write([]byte(fmt.Sprintf("Generating documention in markdown to %s", location)))

			doc.GenMarkdownTreeCustom(cmd.Root(), location, prepender, linkHandler)

			cons.Write([]byte("\nDocumentation generation complete"))

			return nil
		},
	},
	Setup: func(parent *command.Command) error {
		return nil
	},
}
