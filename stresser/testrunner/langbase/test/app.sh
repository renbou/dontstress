#!/bin/bash

set -euo pipefail

for i in tests/*; do
  echo "[+] Running test $i"
  bash $i
done