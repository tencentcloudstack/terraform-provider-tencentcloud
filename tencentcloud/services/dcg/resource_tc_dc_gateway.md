Provides a resource to creating direct connect gateway instance.

Example Usage

If network_type is VPC

```hcl
// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create dc gateway
resource "tencentcloud_dc_gateway" "example" {
  name                = "tf-example"
  network_instance_id = tencentcloud_vpc.vpc.id
  network_type        = "VPC"
  gateway_type        = "NORMAL"
  
  tags = {
    Environment = "production"
    Owner       = "ops-team"
  }
}
```

If network_type is CCN

```hcl
// create ccn
resource "tencentcloud_ccn" "ccn" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "PREPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
  tags = {
    createBy = "terraform"
  }
}

// create dc gateway
resource "tencentcloud_dc_gateway" "example" {
  name                = "tf-example"
  network_instance_id = tencentcloud_ccn.ccn.id
  network_type        = "CCN"
  gateway_type        = "NORMAL"
  
  tags = {
    Team     = "networking"
    Purpose  = "production"
  }
}
```

Update tags

```hcl
resource "tencentcloud_dc_gateway" "example" {
  name                = "tf-example"
  network_instance_id = tencentcloud_ccn.ccn.id
  network_type        = "CCN"
  gateway_type        = "NORMAL"
  
  # Tags can be updated without recreating the gateway
  tags = {
    Environment = "staging"
    Team        = "devops"
    CostCenter  = "IT-001"
  }
}
```

Import

Direct connect gateway instance can be imported, e.g. Tags will be imported automatically.

```
$ terraform import tencentcloud_dc_gateway.example dcg-dr1y0hu7
```