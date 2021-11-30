#!/bin/bash

# This script is used to collect the IPS alert files generated by snort3.
# Please note that after running this script, the alert files in the pods will be removed.
# Running this script will also help clean up alert files which could used up the disk space.

curr_time=$(date "+%Y.%m.%d-%H.%M.%S")

out=`kubectl get pods -n isv-namespace -o wide | grep snort- |cut -d" " -f 1`
array=($out)

for pod in "${array[@]}"; do
    tar_file="alert-$pod-$curr_time.tar.gz"
    echo "saving IPS alerts in pod $pod to $tar_file"
    alert_cmd="supervisorctl stop snort && tar czvfP /var/tmp/$tar_file /var/snort/events/*_alert_*.txt && rm /var/snort/events/*_alert_*.txt && supervisorctl start snort"
    echo $alert_cmd
    cmd_out=`kubectl exec $pod -n isv-namespace -c snort3 -- bash -c "$alert_cmd"`
    echo "tar out: $cmd_out"
    copy_out=`kubectl cp $pod:/var/tmp/$tar_file ./$tar_file -n isv-namespace -c snort3`
    echo "copy out:  $copy_out"
    echo "alerts in $pod saved at $tar_file"
done
