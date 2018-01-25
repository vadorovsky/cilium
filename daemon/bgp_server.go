// Copyright 2018 Authors of Cilium
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

package main

import (
	"fmt"
	"net"
	"time"

	"github.com/cilium/cilium/pkg/node"

	"github.com/golang/glog"
	"github.com/osrg/gobgp/packet/bgp"
	gobgp "github.com/osrg/gobgp/server"
	"github.com/osrg/gobgp/table"
)

var (
	stopBGPServer = make(chan struct{})
)

func (d *Daemon) injectRoute(path *table.Path) error {
	return nil
}

func (d *Daemon) periodicSync() error {
	return nil
}

func (d *Daemon) runBGPServer(stopCh <-chan struct{}) {
	// Initial serving of BGP server.
	go d.bgpServer.Serve()
	defer d.bgpServer.Shutdown()

	watcher := d.bgpServer.Watch(gobgp.WatchBestPath(false))
	for {
		select {
		case ev := <-watcher.Event():
			fmt.Println(ev)
			switch msg := ev.(type) {
			case *gobgp.WatchEventBestPath:
				for _, path := range msg.PathList {
					if path.IsLocal() {
						continue
					}
					if err := d.injectRoute(path); err != nil {
						glog.Errorf("Failed to inject route: %v", err)
						continue
					}
				}
			}
		case <-stopCh:
			return
		default:
			if err := d.periodicSync(); err != nil {
				glog.Errorf("Error during BGP server periodic sync: %c", err)
				if err := d.periodicSync(); err != nil {
					glog.Errorf("Error during BGP server periodic sync: %c", err)
				}
			}
		}
	}
}

func (d *Daemon) AddRoute(ipAddress net.IP, ipRange *net.IPNet) error {
	attrs := []bgp.PathAttributeInterface{
		bgp.NewPathAttributeOrigin(0),
		bgp.NewPathAttributeNextHop(ipAddress.String()),
	}
	cidrLen, _ := ipRange.Mask.Size()
	if _, err := d.bgpServer.AddPath("", []*table.Path{table.NewPath(nil, bgp.NewIPAddrPrefix(uint8(cidrLen), ipRange.IP.String()), false, attrs, time.Now(), false)}); err != nil {
		return err
	}

	return nil
}

func (d *Daemon) EnableBGPServer() error {
	if d.conf.Tunnel != "bgp" {
		return nil
	}

	// Initialize the BGP server.
	if d.bgpServer != nil {
		return fmt.Errorf("bgp server was already initialized")
	}
	d.bgpServer = gobgp.NewBgpServer()

	// Run the BGP server.
	go d.runBGPServer(stopBGPServer)

	// Add the cilium node and cluster ranges as a route.
	d.AddRoute(node.GetExternalIPv4(), node.GetIPv4AllocRange())
	d.AddRoute(node.GetIPv6(), node.GetIPv6NodeRange())
	d.AddRoute(node.GetIPv6(), node.GetIPv6ClusterRange())

	return nil
}
