import cmd

def CreateStack(args):
    gwlb_cmd = ["aws", "cloudformation", "create-stack"]
    gwlb_cmd += ["--stack-name", args.stack_name_prefix + "-GWLB"]
    gwlb_cmd += ["--region", args.region]
    gwlb_cmd += ["--template-body", "file://gwlb/centralized/gwlb-ca.yaml"]
    gwlb_cmd += ["--parameters", "ParameterKey=KeyPairName,ParameterValue=" + args.key_pair]
    gwlb_cmd += ["ParameterKey=AvailabilityZone1,ParameterValue=" + args.az1]
    gwlb_cmd += ["ParameterKey=AvailabilityZone2,ParameterValue=" + args.az2]
    gwlb_cmd += ["--capabilities", "CAPABILITY_NAMED_IAM"]
    out = cmd.Run(gwlb_cmd)
    if out == "":
        print ("failed to create GWLB stack")
        exit()
    print ("GWLB stack output: ", out)
    result = cmd.ParseJson(out)
    return result["StackId"]

def CheckStatus(args):
    desc_cmd = ["aws", "cloudformation", "describe-stacks"]
    desc_cmd += ["--region", args.region]
    desc_cmd += ["--stack-name", args.stack_name_prefix + "-GWLB"]
    out = cmd.Run(desc_cmd)
    if out == "":
        return "Status_Unknown"
    result = cmd.ParseJson(out)
    if "Stacks" in result.keys():
        if len(result["Stacks"]) > 0:
             if "StackStatus" in result["Stacks"][0].keys():
                 return result["Stacks"][0]["StackStatus"]
 
    return "Status_Unknown"

def LoadStackOutput(args):
    desc_cmd = ["aws", "cloudformation", "describe-stacks"]
    desc_cmd += ["--region", args.region]
    desc_cmd += ["--stack-name", args.stack_name_prefix + "-GWLB"]
    out = cmd.Run(desc_cmd)
    if out == "":
        print("failed to describe GWLB stack")
        exit()
    result = cmd.ParseJson(out)
    gwlb_output = {}
    for kv in result["Stacks"][0]["Outputs"]:
        gwlb_output[kv["OutputKey"]] = kv["OutputValue"]
    return gwlb_output
