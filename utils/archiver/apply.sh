#!/usr/bin/env bash

# standard bash error handling
set -o nounset  # treat unset variables as an error and exit immediately.
set -o errexit  # exit immediately when a command fails.
set -E          # needs to be set if we want the ERR trap
set -o pipefail # prevents errors in a pipeline from being masked

deploymentName=kcp-kyma-environment-broker
namespace=kcp-system
kebContainerName=kyma-environment-broker
cloudsqlProxyContainerName=cloudsql-proxy
host=kyma-env-broker

SCRIPT_CLOUDSQL_PROXY_COMMAND=$(kubectl get deployment $deploymentName -n $namespace -o jsonpath=\
"{.spec.template.spec.containers[?(@.name==\"$cloudsqlProxyContainerName\")].command}")
SCRIPT_CLOUDSQL_PROXY_IMAGE=$(kubectl get deployment $deploymentName -n $namespace -o jsonpath=\
"{.spec.template.spec.containers[?(@.name==\"$cloudsqlProxyContainerName\")].image}")

export SCRIPT_CLOUDSQL_PROXY_COMMAND
export SCRIPT_CLOUDSQL_PROXY_IMAGE

envsubst < kyma-environment-broker-archiver.yaml | kubectl apply -f -
