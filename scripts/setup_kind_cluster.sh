#!/bin/bash

# create a local kind cluster
if [[ `kind get clusters` != *"snort3"* ]]; then
    echo "creating kind cluster"
    kind create cluster --config $1 --name snort3
fi
