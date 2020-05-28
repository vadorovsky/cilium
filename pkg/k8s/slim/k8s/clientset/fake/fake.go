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

package fake

import (
	corev1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/client/clientset/versioned/typed/core/v1"
	fakecorev1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/client/clientset/versioned/typed/core/v1/fake"
	discoveryv1beta1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/client/clientset/versioned/typed/discovery/v1beta1"
	fakediscoveryv1beta1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/client/clientset/versioned/typed/discovery/v1beta1/fake"
	networkingv1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/client/clientset/versioned/typed/networking/v1"
	fakenetworkingv1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/client/clientset/versioned/typed/networking/v1/fake"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

type Clientset struct {
	*fake.Clientset
	coreV1           *fakecorev1.FakeCoreV1
	discoveryV1beta1 *fakediscoveryv1beta1.FakeDiscoveryV1beta1
	// NOTE(mrostecki): Huh, the `FakeMetaV1` struct name seems incorrect,
	// but that's what client-gen generated. ¯\_(ツ)_/¯
	networkingV1 *fakenetworkingv1.FakeMetaV1
}

func (c *Clientset) CoreV1() corev1.CoreV1Interface {
	return &fakecorev1.FakeCoreV1{Fake: &c.Fake}
}

func (c *Clientset) DiscoveryV1beta1() discoveryv1beta1.DiscoveryV1beta1Interface {
	return &fakediscoveryv1beta1.FakeDiscoveryV1beta1{Fake: &c.Fake}
}

// NOTE(mrostecki): Huh, the `MetaV1Interface` name seems incorrect,
// but that's what client-gen generated. ¯\_(ツ)_/¯
func (c *Clientset) NetworkingV1() networkingv1.MetaV1Interface {
	// NOTE(mrostecki): Huh, the `FakeMetaV1` struct name seems incorrect,
	// but that's what client-gen generated. ¯\_(ツ)_/¯
	return &fakenetworkingv1.FakeMetaV1{Fake: &c.Fake}
}

func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	var cs Clientset
	cs.Clientset = fake.NewSimpleClientset()
	cs.coreV1 = &fakecorev1.FakeCoreV1{Fake: &cs.Fake}
	cs.discoveryV1beta1 = &fakediscoveryv1beta1.FakeDiscoveryV1beta1{Fake: &cs.Fake}
	// NOTE(mrostecki): Huh, the `FakeMetaV1` struct name seems incorrect,
	// but that's what client-gen generated. ¯\_(ツ)_/¯
	cs.networkingV1 = &fakenetworkingv1.FakeMetaV1{Fake: &cs.Fake}
	return &cs
}
