Provides a resource to create a vpc ipv6_eni_address

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

resource "tencentcloud_eni" "eni" {
  name        = "eni-example"
  vpc_id      = tencentcloud_vpc.vpc.id
  subnet_id   = tencentcloud_subnet.subnet.id
  description = "eni desc."
  ipv4_count  = 1
}

resource "tencentcloud_vpc_ipv6_cidr_block" "example" {
  vpc_id = tencentcloud_vpc.vpc.id
}

resource "tencentcloud_vpc_ipv6_eni_address" "ipv6_eni_address" {
  vpc_id               = tencentcloud_vpc.vpc.id
  network_interface_id = tencentcloud_eni.eni.id
  ipv6_addresses {
    address     = tencentcloud_vpc_ipv6_cidr_block.example.ipv6_cidr_block
    description = "desc."
  }
}
```