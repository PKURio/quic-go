#!/usr/bin/env bash

# Verify that no non-main files import Ginkgo or Gomega.

set -e

HAS_TESTING=false

cd ..
for f in $(find . -name "*.go" ! -name "*_test.go"); do
	if grep -q "github.com/onsi/ginkgo" $f; then
    echo "$f imports github.com/onsi/ginkgo"
    HAS_TESTING=true
	fi
	if grep -q "github.com/onsi/gomega" $f; then
    echo "$f imports github.com/onsi/gomega"
    HAS_TESTING=true
	fi
done

if "$HAS_TESTING"; then
	exit 1
fi
exit 0
