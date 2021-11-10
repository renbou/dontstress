#!/bin/bash

export reset="\x1b[0m"
export green="\x1b[32m"
export cyan="\x1b[36m"

info() {
  echo -e "$cyan""[+] ""$1""$reset"
}

success() {
  echo -e "$green""[v] ""$1""$reset"
}