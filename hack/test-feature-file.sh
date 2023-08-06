#!/bin/env bash


[ -z "$1" ] && echo "feature file is required, pass it as first argument" && exit 1

export EXTRA_BEHAVE_ARGS="-i $1 -k"
export PRIMAZA_CONTROLLER_IMAGE_REF=primaza-controller:testing
export PRIMAZA_AGENTAPP_IMAGE_REF=agentapp:testing
export PRIMAZA_AGENTSVC_IMAGE_REF=agentsvc:testing
export CLUSTER_PROVIDER=external
export MAIN_KUBECONFIG=out/main-kubeconfig
export WORKER_KUBECONFIG=out/worker-kubeconfig

kind delete clusters main worker
kind create cluster --name main
kind create cluster --name worker
kind get kubeconfig --name main > "${MAIN_KUBECONFIG}"
kind get kubeconfig --name worker > "${WORKER_KUBECONFIG}"

bin/yq -i ".clusters[0].cluster.server = \"https://$(docker container inspect main-control-plane | bin/yq '.[0].NetworkSettings.Networks.kind.IPAddress'):6443\"" "${MAIN_KUBECONFIG}"
bin/yq -i ".clusters[0].cluster.server = \"https://$(docker container inspect worker-control-plane | bin/yq '.[0].NetworkSettings.Networks.kind.IPAddress'):6443\"" "${WORKER_KUBECONFIG}"

IMG="$PRIMAZA_CONTROLLER_IMAGE_REF" make primaza docker-build
IMG="$PRIMAZA_AGENTAPP_IMAGE_REF" make agentapp docker-build
IMG="$PRIMAZA_AGENTSVC_IMAGE_REF" make agentsvc docker-build

kind load docker-image --name main "$PRIMAZA_CONTROLLER_IMAGE_REF"
kind load docker-image --name main "$PRIMAZA_AGENTAPP_IMAGE_REF"
kind load docker-image --name main "$PRIMAZA_AGENTSVC_IMAGE_REF"
kind load docker-image --name worker "$PRIMAZA_CONTROLLER_IMAGE_REF"
kind load docker-image --name worker "$PRIMAZA_AGENTAPP_IMAGE_REF"
kind load docker-image --name worker "$PRIMAZA_AGENTSVC_IMAGE_REF"

make kustomize test-acceptance
