#!/usr/bin/python3

import argparse
import time
import eks
import gwlb

parser = argparse.ArgumentParser()
parser.add_argument('--stack-name-prefix', help='cloudformation stack name prefix')
parser.add_argument('--key-pair', help='aws ec2 key pair name')
parser.add_argument('--region', help='aws cloud region')
parser.add_argument('--az1', help='availability zone1')
parser.add_argument('--az2', help='availability zone2')
parser.add_argument('--qs-bucket-name', help='quickstart bucket name')
parser.add_argument('--qs-bucket-region', help='quickstart bucket region')
parser.add_argument('--qs-bucket-prefix', help='quickstart bucket prefix')
parser.add_argument('--skip-gwlb', help='yes or no')
args = parser.parse_args()

if args.stack_name_prefix is None or args.key_pair is None or args.region is None or args.az1 is None or args.az2 is None:
    parser.print_help()
    exit()
if args.qs_bucket_name is None or args.qs_bucket_region is None or args.qs_bucket_prefix is None:
    parser.print_help()
    exit()

if args.skip_gwlb != "yes":
    # create gwlb in a new VPC
    gwlb_sid = gwlb.CreateStack(args)
    print (gwlb_sid + " is being created...")
    success = False
    for i in range(60):
        gwlb_status = gwlb.CheckStatus(args)
        print ("GWLB stack status: ", gwlb_status)
        if gwlb_status == "CREATE_COMPLETE":
            success = True
            break
        time.sleep(60)
    if not success:
        print ("GWLB stack was not created successfully")
        exit()
else:
    print ("skipping GWLB stack")

print ("GWLB stack created successfully")
gwlb_output = gwlb.LoadStackOutput(args)

# Sync quickstart to S3 bucket
eks.SyncQuickStart(args)
# Update EKS parameters
eks.WriteParameters(args, gwlb_output)
# create EKS stack in the GWLB VPC
eks_sid = eks.CreateStack(args)
print (eks_sid + " stack is being created...")
success = False
for i in range(60):
    eks_status = eks.CheckStatus(args)
    print ("EKS stack status: ", eks_status)
    if eks_status == "CREATE_COMPLETE":
        success = True
        break
    time.sleep(60)
if not success:
    print ("EKS stack was not created successfully")
    exit()
print ("EKS stack created successfully")
