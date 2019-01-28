#!/bin/bash

# This script sets up devstack (development environment for OpenStack) with
# kuryr and kuryr-kubernetes which is ready to use with Cilium. It should be
# used in Vagrant dev VM, there is no guarantee that it will work on any other
# environment.

set -e

CILIUM_LOCATION=${CILIUM_LOCATION:-/home/vagrant/go/src/github.com/cilium/cilium}
DEVSTACK_LOCATION=${DEVSTACK_LOCATION:-/opt/stack/devstack}
DEVSTACK_PULL=${DEVSTACK_PULL:-True}
KURYR_K8S_LOCATION=${KURYR_K8S_LOCATION:-/opt/stack/kuryr-kubernetes}
KURYR_K8S_PULL=${KURYR_K8S_PULL:-True}

# Create stack user
if ! getent passwd stack > /dev/null 2>&1; then
	sudo useradd -s /bin/bash -d /opt/stack -m stack
	echo "stack ALL=(ALL) NOPASSWD: ALL" | sudo tee /etc/sudoers.d/stack
fi

sudo chown -R stack:stack /opt/stack

# Clone devstack and kuryr-kubernetes repositories and pull newest changes
if [ ! -e "${DEVSTACK_LOCATION}" ]; then
	sudo su stack -c "git clone https://git.openstack.org/openstack-dev/devstack ${DEVSTACK_LOCATION}"
fi
if [ "${DEVSTACK_PULL}" == "True" ]; then
	sudo su stack -c "cd ${DEVSTACK_LOCATION}; git pull"
fi
if [ ! -e "${KURYR_K8S_LOCATION}" ]; then
	sudo su stack -c "git clone https://git.openstack.org/openstack/kuryr-kubernetes ${KURYR_K8S_LOCATION}"
fi
if [ "${KURYR_K8S_PULL}" ]; then
	sudo su stack -c "cd ${KURYR_K8S_LOCATION}; git pull"
fi

# Copy local.conf
sudo su stack -c "cp ${CILIUM_LOCATION}/contrib/kuryr/local.conf ${DEVSTACK_LOCATION}"

# Execution of this script takes very long time (at least 20-30 minutes) and it
# tends to fail for random  reasons when it's done as a part of Vagrant
# provisioning.
# sudo su stack -c "cd ${DEVSTACK_LOCATION}; ${DEVSTACK_LOCATION}/stack.sh"
