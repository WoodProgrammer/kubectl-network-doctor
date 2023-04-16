#!/bin/bash
## TODO: traceroute args and options
set -e
export ADDRESSES=$(cat hosts.txt)
yellow=$(tput setaf 3 || true)

traceRoute(){
    addr=$1
    traceroute $addr
}

for addr in ${ADDRESSES[@]};
do
    echo "$yellow TraceRoute debug on this address:: $addr"
    traceRoute $addr
done