Provides a resource to create a vpc ipv6_subnet_cidr_block

Example Usage

```hcl
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_vpc_ipv6_cidr_block" "example" {
  vpc_id = tencentcloud_vpc.vpc.id
}

resource "tencentcloud_vpc_ipv6_subnet_cidr_block" "example" {
  vpc_id = tencentcloud_vpc.vpc.id
  ipv6_subnet_cidr_blocks {
    subnet_id       = tencentcloud_subnet.subnet.id
    ipv6_cidr_block = tencentcloud_vpc_ipv6_cidr_block.example.ipv6_cidr_block
  }
}
```

Import

vpc ipv6_subnet_cidr_block can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_ipv6_subnet_cidr_block.ipv6_subnet_cidr_block ipv6_subnet_cidr_block_id
```