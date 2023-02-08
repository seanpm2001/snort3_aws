#!/usr/bin/python3

import os
import sys

def check_container_running(container):
    cmd = "docker container ps | grep " + container
    rv = os.system(cmd)
    if rv == 0:
        return True
    return False
    
def main(argv):
    if len(argv) != 3:
        print ("must specify directory to be mounted, container name and image")
        os.exit(1)
    if check_container_running(argv[1]):
        return

    cmd = "docker run -t -v " + argv[0] + ":/snort3_aws --name " + argv[1] + " -d " + argv[2]
    print ("starting container " + argv[1])
    os.system(cmd)

if __name__ == '__main__':
    main(sys.argv[1:])
