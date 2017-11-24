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
	"fmt"

	"github.com/cilium/cilium/pkg/lock"

	log "github.com/sirupsen/logrus"
)

type kvLocker interface {
	Unlock() error
}

// getLockPath returns the lock path representation of the given path.
func getLockPath(path string) string {
	return path + ".lock"
}

// Lock is a lock return by LockPath
type Lock struct {
	path       string
	localMutex *lock.Mutex
	lock       kvLocker
}

var (
	lockPathsMutex lock.RWMutex
	localMutex     = map[string]*lock.Mutex{}
)

func getLocalPathMutex(path string) *lock.Mutex {
	lockPathsMutex.RLock()
	if localMutex[path] != nil {
		m := localMutex[path]
		lockPathsMutex.RUnlock()
		return m
	}
	lockPathsMutex.RUnlock()

	lockPathsMutex.Lock()
	defer lockPathsMutex.Unlock()

	// We've unlocked the reader lock so check if another writer has come
	// first in the meantime
	if localMutex[path] == nil {
		localMutex[path] = &lock.Mutex{}
	}

	return localMutex[path]
}

// LockPath locks the specified path. The key for the lock is not the path
// provided itself but the path with a suffix of ".lock" appended. The lock
// returned also contains a patch specific local Mutex which will be held.
//
// It is required to call Unlock() on the returned Lock to unlock
func LockPath(path string) (*Lock, error) {
	Trace("Creating lock", nil, log.Fields{fieldKey: path})

	// Take the local lock as kvstore backends don't lock multiple local
	// readers
	localMutex := getLocalPathMutex(path)
	localMutex.Lock()

	lock, err := Client().LockPath(path)
	if err != nil {
		Trace("Unsuccessful lock", err, log.Fields{fieldKey: path})
		localMutex.Unlock()
		return nil, fmt.Errorf("Error while locking path %s: %s", path, err)
	}

	Trace("Successful lock", nil, log.Fields{fieldKey: path})
	return &Lock{lock: lock, localMutex: localMutex, path: path}, nil
}

// Unlock unlocks a lock
func (l *Lock) Unlock() error {
	err := l.lock.Unlock()

	l.localMutex.Unlock()
	if err == nil {
		Trace("Unlocked", nil, log.Fields{fieldKey: l.path})
	}
	return err
}
