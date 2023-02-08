#!/bin/bash

POLICY_VER=3.1.11.0-0
MODULE_VER=3.1.15.0

stdbuf -oL /usr/local/bin/snort -c /var/snort/lightspd/policies/$POLICY_VER/snort.lua --daq-dir /usr/local/lib/daq --daq gwlb -Q -i $ENV_SNORT3_INTERFACE -s $ENV_SNORT3_SNAPLEN --plugin-path /var/snort/lightspd/modules/$MODULE_VER/ubuntu-x64 -l /var/snort/events --lua "alert_full = { file = true }" --lua "alert_fast = { file = true }" --lua "alert_json = { file = true }"
