Provides details about a specific VPC subnet.

This resource can prove useful when a module accepts a subnet id as an input variable and needs to, for example, determine the id of the VPC that the subnet belongs to.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_vpc_subnets.

Example Usage

Query method 1

```hcl
data "tencentcloud_subnet" "subnet" {
  vpc_id    = "vpc-ha5l97e3"
  subnet_id = "subnet-ezgfompo"
}
```

Query method 2

```hcl
data "tencentcloud_subnet" "subnet" {
  vpc_id    = "vpc-ha5l97e3"
  subnet_id = "subnet-ezgfompo"
  cdc_id    = "cluster-lchwgxhs"
}
```
