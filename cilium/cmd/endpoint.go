// Copyright 2017 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
)

// newEndpointCommand returns the endpoint command.
func newEndpointCommand() *cobra.Command {
	endpointCmd := &cobra.Command{
		Use:   "endpoint",
		Short: "Manage endpoints",
	}
	endpointCmd.AddCommand(newEndpointConfigCommand())
	endpointCmd.AddCommand(newEndpointDisconnectCommand())
	endpointCmd.AddCommand(newEndpointGetCommand())
	endpointCmd.AddCommand(newEndpointHealthCommand())
	endpointCmd.AddCommand(newEndpointLabelsCommand())
	endpointCmd.AddCommand(newEndpointListCommand())
	endpointCmd.AddCommand(newEndpointLogCommand())
	endpointCmd.AddCommand(newEndpointRegenerateCommand())
	return endpointCmd
}
