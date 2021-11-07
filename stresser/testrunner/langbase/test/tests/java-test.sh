#!/bin/bash

set -euo pipefail

oldpwd=`pwd`

echo "[+] Setting up Java test"

mkdir /tmp/java-test
cd /tmp/java-test

# write test source
cat <<EOF >Test.java
import java.lang.System;
public class Test {
  public static void main(String[] args) {
    System.out.print("test");
  }
}
EOF

echo "[+] Running Java test"

# compile & run
javac Test.java
result=`java Test`

# clean up
cd $oldpwd
rm -rf /tmp/java-test

# validate test
if [ "$result" != "test" ]; then
  echo "[^] Test failed: $result"
  exit 1
fi

echo "[!] Test passed"