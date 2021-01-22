/**
 * Copyright © 2014-2021 The SiteWhere Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"github.com/gosuri/uitable"
	"io"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/sitewhere/swctl/pkg/action"
	"github.com/sitewhere/swctl/pkg/instance"

	"helm.sh/helm/v3/cmd/helm/require"
	helmAction "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli/output"
)

var createInstanceDesc = `
Use this command to create an Instance of SiteWhere.
For example, to create an instance with name "sitewhere" use:

  swctl create instance sitewhere

To create an instance with the minimal profile use:

	swctl create instance sitewhere -m
`

func newCreateInstanceCmd(cfg *helmAction.Configuration, out io.Writer) *cobra.Command {
	client := action.NewCreateInstance(cfg)
	var outFmt output.Format

	cmd := &cobra.Command{
		Use:               "instance [NAME]",
		Short:             "create an instance",
		Long:              createInstanceDesc,
		Args:              require.ExactArgs(1),
		ValidArgsFunction: noCompletions,
		RunE: func(_ *cobra.Command, args []string) error {
			instanceName, err := client.ExtractInstanceName(args)
			if err != nil {
				return err
			}
			client.InstanceName = instanceName
			results, err := client.Run()
			if err != nil {
				return err
			}
			return outFmt.Write(out, newCreateInstanceWriter(results))
		},
	}

	addCreateInstanceFlags(cmd, cmd.Flags(), client)
	bindOutputFlag(cmd, &outFmt)

	return cmd
}

func addCreateInstanceFlags(cmd *cobra.Command, f *pflag.FlagSet, client *action.CreateInstance) {
	f.StringVarP(&client.Namespace, "namespace", "n", client.Namespace, "Namespace of the instance.")
	f.BoolVarP(&client.Minimal, "minimal", "m", client.Minimal, "Minimal installation.")
	f.StringVarP(&client.Tag, "tag", "t", client.Tag, "Docker image tag.")
	f.StringVar(&client.Registry, "registry", client.Registry, "Docker image registry.")
	f.BoolVarP(&client.Debug, "debug", "d", client.Debug, "Debug mode.")
	f.BoolVar(&client.SkipIstioInject, "skip-istio-inject", client.SkipIstioInject, "Skip Istio Inject namespace label.")
	f.Int32VarP(&client.Replicas, "replicas", "r", client.Replicas, "Number of replicas")
	f.StringVarP(&client.ConfigurationTemplate, "config-template", "c", client.ConfigurationTemplate, "Configuration template.")
	f.StringVarP(&client.DatasetTemplate, "dateset-template", "x", client.DatasetTemplate, "Dataset template.")
}

type createInstancePrinter struct {
	instance *instance.CreateSiteWhereInstance
}

func newCreateInstanceWriter(result *instance.CreateSiteWhereInstance) *createInstancePrinter {
	return &createInstancePrinter{instance: result}
}

func (s createInstancePrinter) WriteJSON(out io.Writer) error {
	return output.EncodeJSON(out, s.instance)
}

func (s createInstancePrinter) WriteYAML(out io.Writer) error {
	return output.EncodeYAML(out, s.instance)
}

func (s createInstancePrinter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("INSTANCE", "STATUS")
	table.AddRow(s.instance.InstanceName, color.Info.Render("Created"))
	return output.EncodeTable(out, table)
}
