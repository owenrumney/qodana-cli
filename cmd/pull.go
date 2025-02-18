/*
 * Copyright 2021-2022 JetBrains s.r.o.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"context"

	"github.com/JetBrains/qodana-cli/core"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// PullOptions represents pull command options.
type PullOptions struct {
	Linter     string
	ProjectDir string
}

// NewPullCommand returns a new instance of the show command.
func NewPullCommand() *cobra.Command {
	options := &PullOptions{}
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull latest version of linter",
		Long:  `An alternative to docker pull.`,
		Run: func(cmd *cobra.Command, args []string) {
			if options.Linter == "" {
				qodanaYaml := core.GetQodanaYaml(options.ProjectDir)
				if qodanaYaml.Linter == "" {
					core.WarningMessage(
						"No valid qodana.yaml found. Have you run %s? Running that for you...",
						core.PrimaryBold("qodana init"),
					)
					options.Linter = core.GetLinter(options.ProjectDir)
					core.EmptyMessage()
				} else {
					options.Linter = qodanaYaml.Linter
				}
			}
			docker, err := client.NewClientWithOpts()
			if err != nil {
				log.Fatal("couldn't connect to docker ", err)
			}
			core.PullImage(context.Background(), docker, options.Linter)
			core.SuccessMessage("Pulled the latest version of linter")
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&options.Linter, "linter", "l", "", "Override linter to use")
	flags.StringVarP(&options.ProjectDir, "project-dir", "i", ".", "Root directory of the inspected project")
	return cmd
}
