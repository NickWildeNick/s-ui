#!/bin/sh
set -e

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR"

cd frontend
npm install
npm run build

cd ..
echo "Backend"

mkdir -p web/html
rm -fr web/html/*
cp -R frontend/dist/* web/html/

export GOCACHE="$SCRIPT_DIR/.cache/go-build"
mkdir -p "$GOCACHE"

BUILD_TAGS="with_quic,with_grpc,with_utls,with_acme,with_gvisor,with_naive_outbound,with_musl,badlinkname,tfogo_checklinkname0,with_tailscale"
LDFLAGS='-w -s -checklinkname=0'
case "$(uname -s)" in
    Darwin)
        LDFLAGS="$LDFLAGS -extldflags \"-Wl,-no_warn_duplicate_libraries\""
        ;;
esac

go build -ldflags "$LDFLAGS" -tags "$BUILD_TAGS" -o sui main.go
