#!/bin/bash

set -euo pipefail

oldpwd=`pwd`

echo "[+] Setting up Golang test"

mkdir /tmp/go-test
cd /tmp/go-test

# write test source
cat <<EOF >test.go
package main
import "fmt"
func main() {
  fmt.Print("test")
}
EOF

echo "[+] Running Golang test"

# compile & run
result=`go run test.go`

# clean up
cd $oldpwd
rm -rf /tmp/go-test

# validate test
if [ "$result" != "test" ]; then
  echo "[^] Test failed: $result"
  exit 1
fi

echo "[!] Test passed"