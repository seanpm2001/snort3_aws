#!/usr/bin/python3

import argparse
import time
import cmd

def check_status(args, stack_name):
    desc_cmd = ["aws", "cloudformation", "describe-stacks"]
    desc_cmd += ["--region", args.region]
    desc_cmd += ["--stack-name", stack_name]
    for i in range(120):
        time.sleep(30)
        out = cmd.Run(desc_cmd)
        try:
            result = cmd.ParseJson(out)
            if "Stacks" in result.keys():
                if len(result["Stacks"]) > 0:
                    if "StackStatus" in result["Stacks"][0].keys():
                        print (stack_name, ": ", result["Stacks"][0]["StackStatus"])
                        if result["Stacks"][0]["StackStatus"] == "DELETE_COMPLETE":
                            break
        except:
            break

parser = argparse.ArgumentParser()
parser.add_argument('--stack-name-prefix', help='cloudformation stack name prefix')
parser.add_argument('--region', help='aws cloud region')
args = parser.parse_args()

if args.stack_name_prefix is None or args.region is None:
    parser.print_help()
    exit()

# delete test tgw stack
delete_cmd = ["aws", "cloudformation", "delete-stack"]
delete_cmd += ["--stack-name", args.stack_name_prefix + "-tgw"]
delete_cmd += ["--region", args.region]
cmd.Run(delete_cmd)
print ("test client stack is being deleted...")
check_status(args, args.stack_name_prefix + "-tgw")

# delete test client stack
delete_cmd = ["aws", "cloudformation", "delete-stack"]
delete_cmd += ["--stack-name", args.stack_name_prefix + "-client"]
delete_cmd += ["--region", args.region]
cmd.Run(delete_cmd)
print ("test client stack is being deleted...")
check_status(args, args.stack_name_prefix + "-client")

# delete test server stack
delete_cmd = ["aws", "cloudformation", "delete-stack"]
delete_cmd += ["--stack-name", args.stack_name_prefix + "-server"]
delete_cmd += ["--region", args.region]
cmd.Run(delete_cmd)
print ("test server stack is being deleted...")
check_status(args, args.stack_name_prefix + "-server")
