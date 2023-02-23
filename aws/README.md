### To install AWS GWLB and EKS cluster

```
./install.py --stack-name-prefix your_stack_name_prefix --region us-west-1 --key-pair you_ec2_key_pair_name --az1 us-west-1a --az2 us-west-1b --qs-bucket-name your_quick_start_bucket_name --qs-bucket-region us-west-1 --qs-bucket-prefix your_quick_start_bucket_prefix
```

Once installation is complate, run the follow command to update kube config for kubectl access.

```
aws eks --region us-west-1 update-kubeconfig --name eks_cluster_name
```

### To install a helm chart

Make sure aws cli version >= 2.7.1 for running helm commands.

Go to project root directory and run

```
kubectl create namespace snort3
helm install your_helm_chart_deployment_name helm/snort3-ips/.
```

### Traffic test

Create test client VPC, test server VPC and TGW.

```
./create_test_stack.py --stack-name-prefix snort3_test --key-pair your_ec2_key_pair_name --region us-west-1 --az1 us-west-1a
```

1. ssh to the bastion host in client VPC.

2. ssh to the application instance in client VPC from the bastion host

3. ping the server instance running in the server VPC

4. ping 8.8.8.8

If the pings can reach the other end, it means snort3 is inspecting geneve traffic.
