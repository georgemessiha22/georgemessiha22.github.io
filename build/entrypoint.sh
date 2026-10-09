#!/bin/sh

set -o errexit
set -o nounset

# TODO: app bootstrapping checks

exec "$@"
