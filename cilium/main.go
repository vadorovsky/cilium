// Copyright 2018 Authors of Cilium
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

package main

import (
	"fmt"
	"os"
	"path"

	ciliumCmd "github.com/cilium/cilium/cilium/cmd"
	"github.com/cilium/cilium/daemon"

	gops "github.com/google/gops/agent"
	"github.com/spf13/cobra"
)

func startGops() {
	if err := gops.Listen(gops.Options{}); err != nil {
		fmt.Fprintf(os.Stderr, "unable to start gops: %s", err)
		os.Exit(-1)
	}
}

func main() {
	base := path.Base(os.Args[0])

	var cmd *cobra.Command
	switch base {
	case "cilium-agent":
		startGops()
		cmd = daemon.NewCommand()
	case "cilium":
		cmd = ciliumCmd.NewRootCommand()
	case "cilium-node-monitor":
		startGops()

	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "unable to execute command %s: %s", base, err)
		os.Exit(-1)
	}
}
