#!/bin/bash

set -euo pipefail

oldpwd=`pwd`

echo "[+] Setting up CPP test"

mkdir /tmp/cpp-test
cd /tmp/cpp-test

# write test source
cat <<EOF >test.cpp
#include <iostream>
int main() {
  std::cout << "test" << std::endl;
  return 0;
}
EOF

echo "[+] Running CPP test"

# compile & run
g++ -o test test.cpp
chmod +x test
result=`./test`

# clean up
cd $oldpwd
rm -rf /tmp/cpp-test

# validate test
if [ "$result" != "test" ]; then
  echo "[^] Test failed: $result"
  exit 1
fi

echo "[!] Test passed"