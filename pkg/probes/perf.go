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

package probes

import (
	"github.com/cilium/cilium/pkg/bpf"
)

func Perf() error {
	events, err := bpf.NewPerCpuEvents(&bpf.PerfEventConfig{
		MapName:      "perf_checker",
		Type:         bpf.PERF_TYPE_SOFTWARE,
		Config:       bpf.PERF_COUNT_SW_BPF_OUTPUT,
		SampleType:   bpf.PERF_SAMPLE_RAW,
		WakeupEvents: 1,
	})
	if err != nil {
		return err
	}

	for {
		todo, err := events.Poll(-1)
		if err != nil {
			return err
		}
		if todo > 0 {
			events.ReadAll(nil, nil, nil)
		}
	}
}
