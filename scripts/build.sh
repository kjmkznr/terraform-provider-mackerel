#!/bin/bash

set -e

XC_ARCH=${XC_ARCH:-"386 amd64 arm"}
XC_OS=${XC_OS:-linux darwin windows freebsd openbsd solaris}
XC_EXCLUDE_OSARCH="!darwin/arm !darwin/386"

echo "==> Removing old directory..."
rm -f pkg/*
mkdir -p pkg/

if ! which gox > /dev/null; then
  echo "==> Installing gox..."
  go get -u github.com/mitchellh/gox
fi

export CGO_ENABLED=0
export LD_FLAGS="-s -w"

echo "==> Building..."
gox \
    -os="${XC_OS}" \
    -arch="${XC_ARCH}" \
    -osarch="${XC_EXCLUDE_OSARCH}" \
    -ldflags "${LD_FLAGS}" \
    -output "pkg/{{.OS}}_{{.Arch}}/terraform-provider-mackerel" \
    ./builtin/bins/provider-mackerel/

echo "==> Packaging..."
for PLATFORM in $(find ./pkg -mindepth 1 -maxdepth 1 -type d); do
    OSARCH=$(basename ${PLATFORM})
    echo "--> ${OSARCH}"

    pushd $PLATFORM >/dev/null 2>&1
    zip ../${OSARCH}.zip ./*
    popd >/dev/null 2>&1
done

# Done!
echo
echo "==> Results:"
ls -hl pkg/
