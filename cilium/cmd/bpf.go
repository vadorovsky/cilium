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

// newBpfCommand returns the bpf command.
func newBpfCommand() *cobra.Command {
	bpfCmd := &cobra.Command{
		Use:   "bpf",
		Short: "Direct access to local BPF maps",
	}
	bpfCmd.AddCommand(newBpfConfigCmd())
	bpfCmd.AddCommand(newBpfCtCmd())
	bpfCmd.AddCommand(newBpfEndpointCommand())
	bpfCmd.AddCommand(newBpfIPCacheCommand())
	bpfCmd.AddCommand(newBpfLBCommand())
	bpfCmd.AddCommand(newBpfMetricsCommand())
	bpfCmd.AddCommand(newBpfPolicyCommand())
	bpfCmd.AddCommand(newBpfProxyCommand())
	bpfCmd.AddCommand(newBpfTunnelCommand())
	return bpfCmd
}
