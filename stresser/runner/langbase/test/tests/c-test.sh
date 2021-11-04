#!/bin/bash

set -euo pipefail

oldpwd=`pwd`

echo "[+] Setting up C test"

mkdir /tmp/c-test
cd /tmp/c-test

# write test source
cat <<EOF >test.c
#include <stdio.h>
int main() {
  printf("test");
  return 0;
}
EOF

echo "[+] Running C test"

# compile & run
gcc -o test test.c
chmod +x test
result=`./test`

# clean up
cd $oldpwd
rm -rf /tmp/c-test

# validate test
if [ "$result" != "test" ]; then
  echo "[^] Test failed: $result"
  exit 1
fi

echo "[!] Test passed"