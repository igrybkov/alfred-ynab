#!/usr/bin/env bash

if [[ -z "${alfred_workflow_bundleid}" ]]; then
  echo "No environment variables found. Running in development mode." > /dev/stderr
  source env.sh
  go run . "$@"
  exit $?
fi

if [[ -f "workflow" ]]; then
  ./workflow "$@"
  exit $?
else
  go run . "$@"
  exit $?
fi
