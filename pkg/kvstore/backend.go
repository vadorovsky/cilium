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

package kvstore

import (
	log "github.com/sirupsen/logrus"
)

type backendOption struct {
	// description is the description of the option
	description string

	// value is the value the option has been configured to
	value string

	// validate, if set, is called to validate the value before assignment
	validate func(value string) error
}

type backendOptions map[string]*backendOption

// backendModule is the interface that each kvstore plugin has to implement.
type backendModule interface {
	// getName returns name of the backend
	getName() string

	// setConfig configures the backend with the specified options. This
	// must be called before init().
	setConfig(opts map[string]string) error

	// getConfig must return the backend configuration.
	getConfig() map[string]string

	// setDummyConfig configured the backend with dummy configuration for
	// testing purposes. This can be used instead of setConfig().
	setDummyConfig()

	// newClient initializes the backend implementation and creates the
	// client structure
	newClient() (KVClient, error)
}

var (
	registeredBackends = map[string]backendModule{}
)

// registerBackend is called by kvstore plugins to register themselves
func registerBackend(name string, module backendModule) {
	if _, ok := registeredBackends[name]; ok {
		log.Panicf("backend with name '%s' already registered", name)
	}

	registeredBackends[name] = module
}

// findBackend returns a kvstore backend by name
func findBackend(name string) backendModule {
	if backend, ok := registeredBackends[name]; ok {
		return backend
	}

	return nil
}
