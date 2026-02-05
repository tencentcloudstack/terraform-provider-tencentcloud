Provide a resource to create a VPC subnet.

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
