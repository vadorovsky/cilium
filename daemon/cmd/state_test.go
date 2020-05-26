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

// +build !privileged_tests

package cmd

import (
	"context"

	"github.com/cilium/cilium/pkg/endpoint"
	"github.com/cilium/cilium/pkg/k8s"
	"github.com/cilium/cilium/pkg/policy"
	"github.com/cilium/cilium/pkg/testutils/allocator"
	testEndpoint "github.com/cilium/cilium/pkg/testutils/endpoint"

	. "gopkg.in/check.v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ds *DaemonSuite) TestValidateEndpoint(c *C) {
	// Pod for ep2.
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: metav1.NamespaceDefault,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "foo",
					Image: "docker.io/cilium/cilium:latest",
				},
			},
		},
	}

	k8s.InitFake()
	clientset := k8s.WatcherCli()
	clientset.CoreV1().Pods(metav1.NamespaceDefault).Create(context.Background(), pod, metav1.CreateOptions{})

	var (
		valid bool
		err   error
	)

	owner := &testEndpoint.DummyOwner{Repo: policy.NewPolicyRepository(nil, nil)}
	// Health endpoint - should be determined invalid.
	/* ep1 := endpoint.NewEndpointWithState(owner, &endpoint.FakeEndpointProxy{}, &allocator.FakeIdentityAllocator{}, 10, endpoint.StateReady)
	ep1.UpdateLabels(context.Background(), labels.LabelHealth, nil, true)
	valid, err = ds.d.validateEndpoint(ep1)
	c.Assert(err, IsNil)
	c.Assert(valid, Equals, false) */

	// Endpoint with an existing pod - should be determined valid.
	ep2 := endpoint.NewEndpointWithState(owner, &endpoint.FakeEndpointProxy{}, &allocator.FakeIdentityAllocator{}, 11, endpoint.StateReady)
	ep2.K8sPodName = "foo"
	ep2.K8sNamespace = metav1.NamespaceDefault
	valid, err = ds.d.validateEndpoint(ep2)
	c.Assert(err, IsNil)
	c.Assert(valid, Equals, true)
}
