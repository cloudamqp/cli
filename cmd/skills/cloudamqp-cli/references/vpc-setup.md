# VPC and network setup

## Create a VPC

```bash
cloudamqp vpc create --name=prod-vpc --region=amazon-web-services::us-east-1 --subnet=10.56.72.0/24
```

The region must match the region of any instances you want to place in the VPC.

## Create an instance in a VPC

```bash
cloudamqp instance create \
  --name=prod-broker \
  --plan=rabbit-1 \
  --region=amazon-web-services::us-east-1 \
  --vpc-id=5678 \
  --wait
```

## List and inspect VPCs

```bash
cloudamqp vpc list
cloudamqp vpc get --id 5678
```

## Update VPC name

```bash
cloudamqp vpc update --id 5678 --name=new-vpc-name
```

## Delete a VPC

```bash
cloudamqp vpc delete --id 5678
```

Remove all instances from the VPC before deleting it.

## Typical setup workflow

```bash
# 1. pick a region
cloudamqp regions --provider=amazon-web-services

# 2. create VPC
cloudamqp vpc create --name=prod-vpc --region=amazon-web-services::us-east-1 --subnet=10.0.0.0/24

# 3. create instance in the VPC
cloudamqp instance create \
  --name=prod-broker \
  --plan=rabbit-1 \
  --region=amazon-web-services::us-east-1 \
  --vpc-id=<vpc-id> \
  --wait

# 4. verify
cloudamqp instance get --id <instance-id>
cloudamqp vpc get --id <vpc-id>
```
