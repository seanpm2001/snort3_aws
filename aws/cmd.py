import subprocess
import json

def Run(cmd, check=True):
    print ("running command: ", ' '.join(cmd))
    out = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    stdout_str = out.communicate()[0]
    if check and out.returncode != 0:
        print ("command return code: ", out.returncode)
        print ("failed to run command")
        exit()
    return stdout_str.decode('utf8')

def ParseJson(out):
    try:
        o = json.loads(out)
        return o
    except:
        print ("failed to parse json string: ", out)
        exit()
