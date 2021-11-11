#!/bin/bash
docker build test/ --tag langbase-test
docker run --rm langbase-test
docker image rm langbase-test