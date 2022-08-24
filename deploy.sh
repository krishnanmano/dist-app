#!/bin/bash

function help() {
  echo "Usage:"
  echo "  deploy -n <number_of_cluster_instance>"
  echo "Example:"
  echo "  deploy -n 10"
  echo "Note: "
  echo "  docker should be installed in the system"
}


function deploy() {
  set -xe
  numberOfInstances=${1}; shift
  init_instance_name=${INSTANCE_PREFIX}"00"
  docker run -d -p8080:8080 -p7960:7960 -e APP_PORT=8080 -e GOSSIP_PORT=7960 --name ${init_instance_name} distapp:0.0.1
  distapp01_ip=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${init_instance_name})

  numberOfInstances=$((numberOfInstances-1))
  for count in $(seq -f %02g 01 ${numberOfInstances})
  do
    exposeable_app_port=$(echo "8080 + ${count}" | bc -l)
    exposeable_gossip_port=$(echo "7960 + ${count}" | bc -l)
    docker run -d -p${exposeable_app_port}:8080 -p${exposeable_gossip_port}:7960 -e APP_PORT=8080 -e GOSSIP_PORT=7960 -e GOSSIP_LEADER=${distapp01_ip}:7960 --name ${INSTANCE_PREFIX}${count} distapp:0.0.1
  done

  set +xe
}

function destroy() {
  set -xe
	docker stop $(docker container ls -a -f "name=${INSTANCE_PREFIX}*" --format '{{.Names}}')
	docker rm $(docker container ls -a -f "name=${INSTANCE_PREFIX}*" --format '{{.Names}}')
	set +xe
}

INSTANCE_PREFIX="distapp"
case ${1} in
  "-n")
    shift
    deploy $1
    ;;
  "-d")
    shift
    destroy
    ;;
  *)
    help
    ;;
esac
