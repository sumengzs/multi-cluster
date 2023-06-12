/*
Copyright 2023 The Multi Cluster Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cluster

type Code int

const (
	// Disabled Indicates that the cluster is disabled.
	Disabled Code = iota
	// Stopped Indicates that the cluster is stopped.
	Stopped
	// Started Indicates that the cluster is started.
	Started
	// Waiting Indicates that waiting the cluster resource ready.
	Waiting
	// Ready Indicates that the cluster resource is already ready.
	Ready
)

var codes = []string{"disabled", "stopped", "started", "waiting", "ready"}

func (c Code) String() string {
	if int(c) < len(codes)-1 {
		return codes[c]
	}
	return "unknown"
}
