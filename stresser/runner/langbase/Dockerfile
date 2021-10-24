FROM ubuntu:20.04

ENV DEBIAN_FRONTEND noninteractive

# Install tools
RUN apt-get update && \
  apt-get install -y --no-install-recommends software-properties-common

# Install c and c++ compiler toolchain
RUN add-apt-repository ppa:ubuntu-toolchain-r/test && \
  apt-get update && \
  apt-get -y install --no-install-recommends gcc-11 g++-11 && \
  update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-11 100 --slave /usr/bin/g++ g++ /usr/bin/g++-11

# Install Java
RUN apt-get update && \
  apt-get -y install --no-install-recommends openjdk-17-jdk-headless

# Install Python
RUN add-apt-repository ppa:deadsnakes/ppa && \
  apt-get update && \
  apt-get -y install --no-install-recommends python3.9

# Install Golang
RUN apt-get update && \
  apt-get -y install --no-install-recommends wget && \
  wget https://golang.org/dl/go1.17.2.linux-amd64.tar.gz -O /tmp/go.tar.gz && \
  tar -C /usr/local -xzf /tmp/go.tar.gz && \
  rm /tmp/go.tar.gz && apt-get -y remove wget

# Clean up
RUN apt-get clean && \
  rm -rf /var/lib/apt/lists/*

ENTRYPOINT []