#!/bin/bash

set -eu

# shellcheck disable=SC2086
cd "$(dirname $0)"

for CMD_DIR in ./cmd/*; do
  if [ -d "${CMD_DIR}" ]; then
    CMD_NAME="$(basename "${CMD_DIR}")"
    echo "installing ${CMD_NAME}"
    go install "${CMD_DIR}"
  fi
done
