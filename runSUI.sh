#!/bin/sh
set -e

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

"$SCRIPT_DIR/build.sh"
cd "$SCRIPT_DIR"
SUI_DB_FOLDER="db" SUI_DEBUG=true ./sui  -local
