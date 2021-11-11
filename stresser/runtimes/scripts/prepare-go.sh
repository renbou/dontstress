#!/bin/bash

dir=`dirname "$0"`
source $dir/util.sh

# Get temporary directory
tmpdir=`mktemp -d`
oldwd=`pwd`
cd "$tmpdir"

### Download ###
gover=`curl 'https://golang.org/VERSION?m=text'`
info "Downloading $gover"
wget "https://dl.google.com/go/$gover.linux-amd64.tar.gz" -O go.tar.gz
tar xzvf go.tar.gz
### --- ###

### Prepare package ###
cd go
info "Removing useless binaries and docs"
rm -f AUTHORS CONTRIBUTING.md CONTRIBUTORS LICENSE PATENTS \
      README.md SECURITY.md VERSION codereview.cfg
rm -rf api doc test
rm -f bin/gofmt
rm -rf pkg/linux_amd64_race
find . -type f -name "*.md" -delete

info "Removing test files"
find . -type f -name "*_test.go" -delete

info "Cleaning up std library"
rm -rf src/vendor src/testing src/testdata src/database src/archive \
       src/compress src/net src/log src/image src/mime src/text/tabwriter \
       src/text/template src/html src/crypto src/cmd src/encoding \
       src/regexp

info "Removing code for unused architectures"
find . -type f -name "*_windows.*" -delete
find . -type f -name "*_windows*.*" -delete
find . -type f -name "*_plan9.*" -delete
find . -type f -name "*_plan9*.*" -delete
find . -type f -name "*_arm64.*" -delete
find . -type f -name "*_arm.*" -delete
find . -type f -name "*_arm*.*" -delete
find . -type f -name "*_mips.*" -delete
find . -type f -name "*_mips*.*" -delete
find . -type f -name "*_android.*" -delete
find . -type f -name "*_android*.*" -delete
find . -type f -name "*_wasm.*" -delete
find . -type f -name "*_wasm*.*" -delete
find . -type f -name "*_ppc64.*" -delete
find . -type f -name "*_ppc64le.*" -delete
find . -type f -name "*_ppc64*.*" -delete
find . -type f -name "*_freebsd.*" -delete
find . -type f -name "*_freebsd*.*" -delete
find . -type f -name "*_openbsd.*" -delete
find . -type f -name "*_openbsd*.*" -delete
find . -type f -name "*_netbsd.*" -delete
find . -type f -name "*_netbsd*.*" -delete
find . -type f -name "*_solaris.*" -delete
find . -type f -name "*_solaris*.*" -delete
cd ..
### --- ###

### Clean up ###
info "Cleaning up"
rm -rf go.tar.gz
### --- ###

success "Go downloaded, get it from here:"
success "$tmpdir/go"
cd "$oldwd"