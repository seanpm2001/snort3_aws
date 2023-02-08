#!/usr/bin/python3

import os
import sys

def main(argv):
    if len(argv) != 1:
        print ("must specify container to be removed")
        os.exit(1)
    print ("remove container " + argv[0])
    cmd = "docker stop " + argv[0]
    os.system(cmd)
    cmd = "docker rm " + argv[0]
    os.system(cmd)

if __name__ == '__main__':
    main(sys.argv[1:])
