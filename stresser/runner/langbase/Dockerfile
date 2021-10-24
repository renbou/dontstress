FROM ubuntu:20.04

LABEL org.opencontainers.image.source="https://github.com/renbou/dontstress/tree/master/stresser/runner/langbase"
LABEL org.opencontainers.image.authors="renbou"

ENV DEBIAN_FRONTEND=noninteractive

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
  apt-get -y install --no-install-recommends python3.9 && \
  update-alternatives --install /usr/bin/python python /usr/bin/python3 100

# Install Golang
RUN apt-get update && \
  apt-get -y install --no-install-recommends wget && \
  wget https://golang.org/dl/go1.17.2.linux-amd64.tar.gz -O /tmp/go.tar.gz && \
  tar -C /usr/local -xzf /tmp/go.tar.gz && \
  rm /tmp/go.tar.gz && apt-get -y remove wget

# Clean up
RUN apt-get -y remove software-properties-common && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

ENV JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64
ENV PATH="$PATH:$JAVA_HOME/bin:/usr/local/go/bin"

ENTRYPOINT []