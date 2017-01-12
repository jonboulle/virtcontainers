//
// Copyright (c) 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package virtcontainers

import (
	"fmt"
	"os"
)

// ProxyType describes a proxy type.
type ProxyType string

const (
	// CCProxyType is the ccProxy.
	CCProxyType ProxyType = "ccProxy"

	// NoopProxyType is the noopProxy.
	NoopProxyType ProxyType = "noopProxy"
)

// Set sets a proxy type based on the input string.
func (pType *ProxyType) Set(value string) error {
	switch value {
	case "noopProxy":
		*pType = NoopProxyType
		return nil
	case "ccProxy":
		*pType = CCProxyType
		return nil
	default:
		return fmt.Errorf("Unknown proxy type %s", value)
	}
}

// String converts a proxy type to a string.
func (pType *ProxyType) String() string {
	switch *pType {
	case NoopProxyType:
		return string(NoopProxyType)
	case CCProxyType:
		return string(CCProxyType)
	default:
		return ""
	}
}

// newProxy returns a proxy from a proxy type.
func newProxy(pType ProxyType) (proxy, error) {
	switch pType {
	case NoopProxyType:
		return &noopProxy{}, nil
	case CCProxyType:
		return &ccProxy{}, nil
	default:
		return &noopProxy{}, nil
	}
}

// IOStream holds three file descriptors returned by the proxy.
// Those file descriptors will be given to the calling process,
// so that it can interact with a workload running on a container.
type IOStream struct {
	Stdin    *os.File
	Stdout   *os.File
	Stderr   *os.File
	StdinID  uint64
	StdoutID uint64
	StderrID uint64
}

// proxy is the virtcontainers proxy interface.
type proxy interface {
	// register connects and registers the proxy to the given VM.
	// It also returns streams related to containers workloads.
	register(pod Pod) ([]IOStream, error)

	// unregister unregisters and disconnects the proxy from the given VM.
	unregister(pod Pod) error

	// connect gets the proxy a handle to a previously registered VM.
	// It also returns streams related to containers workloads.
	connect(pod Pod) (IOStream, error)

	// disconnect disconnects from the proxy.
	disconnect() error

	// sendCmd sends a command to the agent inside the VM through the proxy.
	sendCmd(cmd interface{}) (interface{}, error)
}
