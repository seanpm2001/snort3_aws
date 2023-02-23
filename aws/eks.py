import cmd
import json

def CreateStack(args):
    template_url = "https://" + args.qs_bucket_name + ".s3." + args.qs_bucket_region + ".amazonaws.com/" + args.qs_bucket_prefix + "/templates/amazon-eks-entrypoint-existing-vpc.template.yaml"
    eks_cmd = ["aws", "cloudformation", "create-stack"]
    eks_cmd += ["--stack-name", args.stack_name_prefix + "-eks"]
    eks_cmd += ["--region", args.region]
    eks_cmd += ["--template-url", template_url]
    eks_cmd += ["--capabilities", "CAPABILITY_NAMED_IAM", "CAPABILITY_AUTO_EXPAND", "CAPABILITY_IAM"]
    eks_cmd += ["--parameters", "file://eks/" + args.stack_name_prefix + ".json"]
    out = cmd.Run(eks_cmd)
    if out == "":
        print ("failed to create GWLB stack")
        exit()
    print ("GWLB stack output: ", out)
    result = cmd.ParseJson(out)
    return result["StackId"]

def WriteParameters(args, gwlb):
    params = {}
    params["KeyPairName"] = args.key_pair
    params["EKSClusterName"] = args.stack_name_prefix + "-eks"
    params["QSS3BucketRegion"] = args.qs_bucket_region
    params["QSS3BucketName"] = args.qs_bucket_name
    params["QSS3KeyPrefix"] = args.qs_bucket_prefix + "/"
    params["NodeInstanceType"] = "m5.large"
    params["ProvisionBastionHost"] = "Disabled"
    params["PerAccountSharedResources"] = "No"
    params["PerRegionSharedResources"] = "No"
    params["NumberOfNodes"] = "1"
    params["VPCID"] = gwlb["ApplianceVpcId"]
    params["PrivateSubnet1ID"] = gwlb["PrivateSubnet1"]
    params["PrivateSubnet2ID"] = gwlb["PrivateSubnet2"]
    params["PublicSubnet1ID"] = gwlb["PublicSubnet1"]
    params["PublicSubnet2ID"] = gwlb["PublicSubnet2"]

    json_data = []
    for key in params:
        json_data += [{"ParameterKey": key, "ParameterValue": params[key]}]
    param_path = "./eks/" + args.stack_name_prefix + ".json"
    print ("EKS parameters: ", json.dumps(json_data))
    param_file = open(param_path, "w")
    param_file.write(json.dumps(json_data, indent=4))
    param_file.close()

def SyncQuickStart(args):
    sync_cmd = ["aws", "s3", "sync", "eks/quickstart/templates"]
    dest = "s3://" + args.qs_bucket_name + "/" + args.qs_bucket_prefix + "/templates"
    sync_cmd += [dest]
    cmd.Run(sync_cmd) 
    sync_cmd = ["aws", "s3", "sync", "eks/quickstart/submodules"]
    dest = "s3://" + args.qs_bucket_name + "/" + args.qs_bucket_prefix + "/submodules"
    sync_cmd += [dest]
    cmd.Run(sync_cmd) 

def CheckStatus(args):
    desc_cmd = ["aws", "cloudformation", "describe-stacks"]
    desc_cmd += ["--region", args.region]
    desc_cmd += ["--stack-name", args.stack_name_prefix + "-eks"]
    out = cmd.Run(desc_cmd)
    if out == "":
        return "Status_Unknown"
    result = cmd.ParseJson(out)
    if "Stacks" in result.keys():
        if len(result["Stacks"]) > 0:
             if "StackStatus" in result["Stacks"][0].keys():
                 return result["Stacks"][0]["StackStatus"]
 
    return "Status_Unknown"
