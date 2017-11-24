// Copyright 2016-2017 Authors of Cilium
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
	"encoding/json"
	"time"

	"github.com/cilium/cilium/common/types"

	log "github.com/sirupsen/logrus"
)

// KeyValuePairs is a map of key=value pairs
type KeyValuePairs map[string][]byte

type KVClient interface {
	// BEGIN Obsolete API
	GetValue(k string) (json.RawMessage, error)
	SetValue(k string, v interface{}) error
	InitializeFreeID(path string, firstID uint32) error
	GetMaxID(key string, firstID uint32) (uint32, error)
	SetMaxID(key string, firstID, maxID uint32) error

	GASNewL3n4AddrID(basePath string, baseID uint32, lAddrID *types.L3n4AddrID) error
	// END Obsolete API

	// Status returns the status of he kvstore client including an
	// eventual error
	Status() (string, error)

	// LockPath locks the provided path
	LockPath(path string) (kvLocker, error)

	// Get returns value of key
	Get(key string) ([]byte, error)

	// GetPrefix returns the first key which matches the prefix
	GetPrefix(prefix string) ([]byte, error)

	// Set sets value of key
	Set(key string, value []byte) error

	// Delete deletes a key
	Delete(key string) error

	DeletePrefix(path string) error

	// CreateOnly atomically creates a key or fails if it already exists
	CreateOnly(key string, value []byte, lease bool) error

	// CreateIfExists creates a key with the value only if key condKey exists
	CreateIfExists(condKey, key string, value []byte, lease bool) error

	// ListPrefix returns a list of keys matching the prefix
	ListPrefix(prefix string) (KeyValuePairs, error)

	// Watch starts watching for changes in a prefix. If list is true, the
	// current keys matching the prefix will be listed and reported as new
	// keys first.
	Watch(w *Watcher)

	// CreateLease creates a lease with the specified ttl
	CreateLease(ttl time.Duration) (interface{}, error)

	// KeepAlive keeps a lease previously created with CreateLease alive
	// for the duration specified at CreateLease time
	KeepAlive(lease interface{}) error

	// DeleteLease deletes a lease
	DeleteLease(interface{}) error

	// Close closes the kvstore client
	Close()

	// GetCapabilities returns the capabilities of the backend
	GetCapabilities() Capabilities

	// Encodes a binary slice into a character set that the backend
	// supports
	Encode(in []byte) string

	// Decodes a key previously encoded back into the original binary slice
	Decode(in string) ([]byte, error)
}

// Capabilities is a bitmask to indicate the capabilities of a backend
type Capabilities uint32

const (
	// CapabilityCreateIfExists is true if CreateIfExists is functional
	CapabilityCreateIfExists Capabilities = 1 << 0

	// CapabilityDeleteOnZeroCount is true if DeleteOnZeroCount is functional
	CapabilityDeleteOnZeroCount Capabilities = 1 << 1

	// BaseKeyPrefix is the base prefix that should be used for all keys
	BaseKeyPrefix = "cilium"
)

// Get returns value of key
func Get(key string) ([]byte, error) {
	v, err := Client().Get(key)
	Trace("Get", err, log.Fields{fieldKey: key, fieldValue: string(v)})
	return v, err
}

// GetPrefix returns the first key which matches the prefix
func GetPrefix(prefix string) ([]byte, error) {
	v, err := Client().GetPrefix(prefix)
	Trace("GetPrefix", err, log.Fields{fieldPrefix: prefix, fieldValue: string(v)})
	return v, err
}

// ListPrefix returns the list of keys matching the prefix
func ListPrefix(prefix string) (KeyValuePairs, error) {
	v, err := Client().ListPrefix(prefix)
	Trace("ListPrefix", err, log.Fields{fieldPrefix: prefix, fieldNumEntries: len(v)})
	return v, err
}

// CreateOnly atomically creates a key or fails if it already exists
func CreateOnly(key string, value []byte, lease bool) error {
	err := Client().CreateOnly(key, value, lease)
	Trace("CreateOnly", err, log.Fields{fieldKey: key, fieldValue: string(value), fieldAttachLease: lease})
	return err
}

// CreateIfExists creates a key with the value only if key condKey exists
func CreateIfExists(condKey, key string, value []byte, lease bool) error {
	err := Client().CreateIfExists(condKey, key, value, lease)
	Trace("CreateIfExists", err, log.Fields{fieldKey: key, fieldValue: string(value), fieldCondition: condKey, fieldAttachLease: lease})
	return err
}

// Set sets the value of a key
func Set(key string, value []byte) error {
	err := Client().Set(key, value)
	Trace("Set", err, log.Fields{fieldKey: key, fieldValue: string(value)})
	return err
}

// Delete deletes a key
func Delete(key string) error {
	err := Client().Delete(key)
	Trace("Delete", err, log.Fields{fieldKey: key})
	return err
}

// DeletePrefix deletes all keys matching a prefix
func DeletePrefix(prefix string) error {
	err := Client().DeletePrefix(prefix)
	Trace("DeletePrefix", err, log.Fields{fieldPrefix: prefix})
	return err
}

// GetCapabilities returns the capabilities of the backend
func GetCapabilities() Capabilities {
	return Client().GetCapabilities()
}

// Encode encodes a binary slice into a character set that the backend supports
func Encode(in []byte) string {
	out := Client().Encode(in)
	Trace("Encode", nil, log.Fields{"in": in, "out": out})
	return out
}

// Decode decodes a key previously encoded back into the original binary slice
func Decode(in string) ([]byte, error) {
	out, err := Client().Decode(in)
	Trace("Decode", err, log.Fields{"in": in, "out": out})
	return out, err
}
