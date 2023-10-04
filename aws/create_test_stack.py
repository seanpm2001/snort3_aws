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
        result = cmd.ParseJson(out)
        if "Stacks" in result.keys():
            if len(result["Stacks"]) > 0:
                if "StackStatus" in result["Stacks"][0].keys():
                    print (stack_name, ": ", result["Stacks"][0]["StackStatus"])
                    if result["Stacks"][0]["StackStatus"] == "CREATE_COMPLETE":
                        break

def get_stack_output(args, stack_name):
    desc_cmd = ["aws", "cloudformation", "describe-stacks"]
    desc_cmd += ["--region", args.region]
    desc_cmd += ["--stack-name", stack_name]
    out = cmd.Run(desc_cmd)
    if out == "":
        print (stack_name, " does not exist")
        exit()
    result = cmd.ParseJson(out)
    output = {}
    for kv in result["Stacks"][0]["Outputs"]:
        output[kv["OutputKey"]] = kv["OutputValue"]
    return output

parser = argparse.ArgumentParser()
parser.add_argument('--stack-name-prefix', help='cloudformation stack name prefix')
parser.add_argument('--key-pair', help='aws ec2 key pair name')
parser.add_argument('--region', help='aws cloud region')
parser.add_argument('--az1', help='availability zone1')
parser.add_argument('--gwlb-stack-name', help='GWLB stack name')
args = parser.parse_args()

if args.stack_name_prefix is None or args.key_pair is None or args.region is None or args.az1 is None:
    parser.print_help()
    exit()
if args.gwlb_stack_name is None:
    parser.print_help()
    exit()

# create test client VPC
create_cmd = ["aws", "cloudformation", "create-stack"]
create_cmd += ["--stack-name", args.stack_name_prefix + "-client"]
create_cmd += ["--region", args.region]
create_cmd += ["--template-body", "file://gwlb/centralized/client-vpc.yaml"]
create_cmd += ["--parameters", "ParameterKey=KeyPairName,ParameterValue=" + args.key_pair]
create_cmd += ["ParameterKey=AvailabilityZone1,ParameterValue=" + args.az1]
create_cmd += ["--capabilities", "CAPABILITY_NAMED_IAM"]
out = cmd.Run(create_cmd)
result = cmd.ParseJson(out)
print ("test client VPC ", result["StackId"], " is being created...")
check_status(args, args.stack_name_prefix + "-client")

# create test server VPC
create_cmd = ["aws", "cloudformation", "create-stack"]
create_cmd += ["--stack-name", args.stack_name_prefix + "-server"]
create_cmd += ["--region", args.region]
create_cmd += ["--template-body", "file://gwlb/centralized/server-vpc.yaml"]
create_cmd += ["--parameters", "ParameterKey=KeyPairName,ParameterValue=" + args.key_pair]
create_cmd += ["ParameterKey=AvailabilityZone1,ParameterValue=" + args.az1]
create_cmd += ["--capabilities", "CAPABILITY_NAMED_IAM"]
out = cmd.Run(create_cmd)
result = cmd.ParseJson(out)
print ("test server VPC ", result["StackId"], " is being created...")
check_status(args, args.stack_name_prefix + "-server")

# create test tgw stack
gwlb_output = get_stack_output(args, args.gwlb_stack_name)
print (gwlb_output)
client_output = get_stack_output(args, args.stack_name_prefix + "-client")
print (client_output)
server_output = get_stack_output(args, args.stack_name_prefix + "-server")
print (server_output)
create_cmd = ["aws", "cloudformation", "create-stack"]
create_cmd += ["--stack-name", args.stack_name_prefix + "-tgw"]
create_cmd += ["--region", args.region]
create_cmd += ["--template-body", "file://gwlb/centralized/tgw-ca.yaml"]
create_cmd += ["--parameters", "ParameterKey=ApplianceVpcId,ParameterValue=" + gwlb_output["ApplianceVpcId"]]
create_cmd += ["ParameterKey=ApplianceVpcTgwAttachSubnet1Id,ParameterValue=" + gwlb_output["ApplianceTgwAttachSubnet1Id"]]
create_cmd += ["ParameterKey=ApplianceVpcApplianceRtb1Id,ParameterValue=" + gwlb_output["ApplianceRtb1Id"]]
create_cmd += ["ParameterKey=Spoke1VpcId,ParameterValue=" + client_output["SpokeVpcId"]]
create_cmd += ["ParameterKey=Spoke1VpcTgwAttachSubnet1Id,ParameterValue=" + client_output["SpokeTgwAttachSubnet1Id"]]
create_cmd += ["ParameterKey=Spoke1VpcRtb1Id,ParameterValue=" + client_output["SpokeApplicationRouteTableId"]]
create_cmd += ["ParameterKey=Spoke2VpcId,ParameterValue=" + server_output["SpokeVpcId"]]
create_cmd += ["ParameterKey=Spoke2VpcTgwAttachSubnet1Id,ParameterValue=" + server_output["SpokeTgwAttachSubnet1Id"]]
create_cmd += ["ParameterKey=Spoke2VpcRtb1Id,ParameterValue=" + server_output["SpokeApplicationRouteTableId"]]
create_cmd += ["--capabilities", "CAPABILITY_NAMED_IAM"]
out = cmd.Run(create_cmd)
result = cmd.ParseJson(out)
print ("test tgw ", result["StackId"], " is being created...")
check_status(args, args.stack_name_prefix + "-tgw")
