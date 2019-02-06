// Copyright 2019 Authors of Cilium
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

package workloads

import (
	"github.com/containers/libpod/libpod"
	"github.com/containers/storage/pkg/reexec"
	"golang.org/x/net/context"
)

const (
	Podman WorkloadRuntimeType = "podman"
)

var (
	podmanInstance = &podmanModule{}
)

type podmanModule struct {
}

func init() {
	registerWorkload(Podman, podmanInstance)
}

func (c *podmanModule) getName() string {
	return string(Podman)
}

func (c *podmanModule) setConfigDummy() {

}

type podmanClient struct {
	*libpod.Runtime
}

func newPodmanClient(opts workloadRuntimeOupts) (WorkloadRuntime, error) {
	if reexec.Init() {
		return
	}

	runtime, err := adapter.GetRuntime()
	if err != nil {
		return nil, err
	}

	return podmanClient{Runtime: runtime}, nil
}

func (p *podmanClient) IsRunning(ep *endpoint.Endpoint) bool {
	return true
}

func (p *podmanClient) Status() *models.Status {
	// libpod has no daemon, so it's always OK
	return &models.Status{State: models.StatusStateOk, Msg: "libpod: OK"}
}

func (p *podmanClient) containersList(ctx context.Context) (*libpod.Container, error) {
	cList, err := p.GetAllContainers()
	return cList, err
}

func (p *podmanClient) GetAllInfraContainersPID() (map[string]int, error) {
	timeoutCtx, cancel := ctx.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cList, err := Client().containersList(timeoutCtx)
	if err != nil {
		return nil, err
	}
	pids := map[string]int{}
	for _, c := range cList {
		inspectData, err := c.Inspect(true)
		if err != nil {
			return nil, err
		}
	}

	return pids, nil
}
