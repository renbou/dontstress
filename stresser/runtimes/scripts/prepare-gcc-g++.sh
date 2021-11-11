#!/bin/bash

dir=`dirname "$0"`
source $dir/util.sh

MUSL_VERSION=x86_64-linux-musl-native

# Get temporary directory
tmpdir=`mktemp -d`
oldwd=`pwd`
cd "$tmpdir"

info "Getting musl"
wget "https://musl.cc/$MUSL_VERSION.tgz" -O musl.tgz
tar xzvf musl.tgz
mv "$MUSL_VERSION" musl

info "Removing useless binaries"
cd musl/bin
### Prepare binaries ##
# Get gcc version
gccname=`ls x86_64-linux-musl-gcc-[1-9]* | tr -d '\n'`
gccver=${gccname#*x86_64-linux-musl-gcc-}

# Delete garbage in bin
rm -f addr2line c++ c++filt cpp elfedit gcc-nm gcov gcov-dump gcov-tool \
      gfortran gprof lto-dump nm objcopy objdump readelf strings strip \
      x86_64-linux-musl-gcc-nm x86_64-linux-musl-gfortran

# Delete hard links and make soft ones instead
rm -f ar as g++ gcc gcc-ar gcc-ranlib ld ld.bfd ld.gold ranlib \
      x86_64-linux-musl-c++ x86_64-linux-musl-cc x86_64-linux-musl-gcc
ln -s ../x86_64-linux-musl/bin/ar       ar
ln -s ../x86_64-linux-musl/bin/as       as
ln -s x86_64-linux-musl-g++             g++
ln -s x86_64-linux-musl-gcc-$gccver     gcc
ln -s x86_64-linux-musl-gcc-ar          gcc-ar
ln -s x86_64-linux-musl-gcc-ranlib      gcc-ranlib
ln -s ../x86_64-linux-musl/bin/ld.bfd   ld
ln -s ../x86_64-linux-musl/bin/ld.bfd   ld.bfd
ln -s ../x86_64-linux-musl/bin/ld.gold  ld.gold
ln -s ../x86_64-linux-musl/bin/ranlib   ranlib
ln -s x86_64-linux-musl-g++             x86_64-linux-musl-c++
ln -s x86_64-linux-musl-gcc             x86_64-linux-musl-cc
ln -s x86_64-linux-musl-gcc-$gccver     x86_64-linux-musl-gcc
### --- ###
cd ../..

info "Removing docs"
### Remove docs ###
rm -rf musl/share
### --- ###

info "Cleaning up"
### Remove sources ###
rm -rf musl.tgz
### --- ###

success "Musl packaging complete, get it from here:"
success "$tmpdir/musl"
cd "$oldwd"