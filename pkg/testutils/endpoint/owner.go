// Copyright 2020 Authors of Cilium
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

package endpoint

import (
	"context"

	"github.com/cilium/cilium/pkg/datapath"
	"github.com/cilium/cilium/pkg/identity/cache"
	"github.com/cilium/cilium/pkg/lock"
	monitorAPI "github.com/cilium/cilium/pkg/monitor/api"
	"github.com/cilium/cilium/pkg/policy"
)

// DummyOwner implements pkg/endpoint/regeneration/Owner. Used for unit testing.
type DummyOwner struct {
	Repo *policy.Repository
}

// GetPolicyRepository returns the policy repository of the owner.
func (d *DummyOwner) GetPolicyRepository() *policy.Repository {
	return d.Repo
}

// QueueEndpointBuild does nothing.
func (d *DummyOwner) QueueEndpointBuild(ctx context.Context, epID uint64) (func(), error) {
	return nil, nil
}

// GetCompilationLock does nothing.
func (d *DummyOwner) GetCompilationLock() *lock.RWMutex {
	return nil
}

// GetCIDRPrefixLengths does nothing.
func (d *DummyOwner) GetCIDRPrefixLengths() (s6, s4 []int) {
	return nil, nil
}

// SendNotification does nothing.
func (d *DummyOwner) SendNotification(typ monitorAPI.AgentNotification, text string) error {
	return nil
}

// Datapath returns a nil datapath.
func (d *DummyOwner) Datapath() datapath.Datapath {
	return nil
}

// GetNodeSuffix does nothing.
func (d *DummyOwner) GetNodeSuffix() string {
	return ""
}

// UpdateIdentities does nothing.
func (d *DummyOwner) UpdateIdentities(added, deleted cache.IdentityCache) {}
