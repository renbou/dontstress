#!/bin/bash

set -euo pipefail

oldpwd=`pwd`

echo "[+] Setting up Python test"

mkdir /tmp/python-test
cd /tmp/python-test

# write test source
cat <<EOF >test.py
print("test",end="")
EOF

echo "[+] Running Python test"

# compile & run
result=`python test.py`

# clean up
cd $oldpwd
rm -rf /tmp/python-test

# validate test
if [ "$result" != "test" ]; then
  echo "[^] Test failed: $result"
  exit 1
fi

echo "[!] Test passed"