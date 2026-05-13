Provide a resource to create a VPC subnet.

~> **NOTE:** In accordance with VPC business requirements, the default value for `is_multicast` has been updated to `false`(previously `true`) in version `v1.82.93` of the provider. If you wish to utilize this feature, you must first contact the VPC product team to have your account added to the whitelist, and then set the `is_multicast` field to `true`.

Example Usage

Create a normal VPC subnet

```hcl
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name         = "vpc-example"
  cidr_block   = "10.0.0.0/16"
  is_multicast = false
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  is_multicast      = false
}
```

Create a CDC instance VPC subnet

```hcl
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name         = "vpc-example"
  cidr_block   = "10.0.0.0/16"
  is_multicast = false
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  cdc_id            = "cluster-lchwgxhs"
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  is_multicast      = false
}

```

Import

Vpc subnet instance can be imported, e.g.

```
$ terraform import tencentcloud_subnet.subnet subnet-b8j03v0c
```
