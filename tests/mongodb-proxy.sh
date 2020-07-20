#!/bin/bash

dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
source "${dir}/helpers.bash"
# dir might have been overwritten by helpers.bash
dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

TEST_NAME=$(get_filename_without_extension $0)
LOGS_DIR="${dir}/cilium-files/${TEST_NAME}/logs"
redirect_debug_logs ${LOGS_DIR}

set -ex

function cleanup {
  monitor_stop
  cilium policy delete --all 2> /dev/null || true
  docker rm -f mongo-server mongo-client 2> /dev/null || true
}

function finish_test {
  echo cleanup
}

trap finish_test EXIT

IMAGE="mongo"
SERVER_LABEL="mongo-server"
CLIENT_LABEL="mongo-client"
TAG="3.6.18"
#TAG="4.2.8"

CLIENT_RUN="docker run --rm -t --net=cilium --name $CLIENT_LABEL -l $CLIENT_LABEL $IMAGE:$TAG mongo"
cleanup
logs_clear

function proxy_init {
  log "beginning proxy_init"
  create_cilium_docker_network

  docker run -dt --net=cilium --name $SERVER_LABEL -l $SERVER_LABEL --publish 27017:27017 $IMAGE:$TAG
  wait_for_docker_ipv6_addr $SERVER_LABEL

  log "waiting for mongodb server"
  while ! cilium endpoint list -o jsonpath='{range [*]}{.status.identity.id}{" "}{.status.identity.labels}{"\n"}' | grep '[0-9].*mongo-server' ; do
    log "waiting..."
    sleep 1
  done

  echo "probing until mongodb server is responsive"
  until docker exec -i mongo-server mongo --eval "db" 2>/dev/null >/dev/null; do
    echo "."
    sleep 1
  done

  echo "Creating a sample collection"
  docker exec -i mongo-server mongo cilium --eval "db.coll.insert({x: 1})"

  SERVER_IP4=$(docker inspect --format '{{ .NetworkSettings.Networks.cilium.IPAddress }}' mongo-server)

  echo "Testing client without policy"
  $CLIENT_RUN $SERVER_IP4/cilium --eval "db.coll.find({x: 1})"

  monitor_start
  log "finished proxy_init"
}

function policy_single_egress {
    cilium policy delete --all
    cat <<EOF | policy_import_and_wait -
[{
    "endpointSelector": {"matchLabels":{"id.server":""}},
    "ingress": [{
        "fromEndpoints": [
            {"matchLabels":{"reserved:host":""}},
            {"matchLabels":{"mongo-server":""}}
        ]
    }]
},{
    "endpointSelector": {"matchLabels":{"mongo-client":""}},
    "egress": [{
        "toPorts": [{
            "ports": [{"port": "27017", "protocol": "TCP"}],
            "rules": {
                "l7proto": "envoy.filters.network.mongo_proxy",
                "l7": [{
                    "action": "deny",
                    "*": "query"
                }]
            }
        }]
    }]
}]
EOF
}

function proxy_test {
  log "beginning Mongo proxy test"
  monitor_clear

  log "trying to reach Mongo server at $SERVER_IP4 from client"
  if $CLIENT_RUN $SERVER_IP4 --eval "db"; then
    echo "Success"
  else
    abort "Mongo query failed"
  fi

  log "trying to query denied collection at $SERVER_IP4 from client"
  if $CLIENT_RUN $SERVER_IP4/cilium --eval "db.coll.find()"; then
    abort "Mongo query should have failed, but it succeeded"
  else
    echo "Mongo query failed as expected"
  fi

  monitor_dump

  log "finished Mongo proxy test"
}

proxy_init

policy_single_egress

proxy_test

# Leave test setup behind for manual testing
#
# log "deleting all policies from Cilium"
# cilium policy delete --all 2> /dev/null || true
# log "removing containers"
# docker rm -f mongo-server mongo-client 2> /dev/null || true

test_succeeded "${TEST_NAME}"
