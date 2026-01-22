Provides a resource to create a vpc eni ipv6 address

Example Usage

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = "ap-guangzhou-6"
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
    ipv6_cidr_block = "2402:4e00:1015:7500::/64"
  }
}

resource "tencentcloud_eni" "example" {
  name        = "tf-example"
  vpc_id      = tencentcloud_vpc.vpc.id
  subnet_id   = tencentcloud_subnet.subnet.id
  description = "eni description."
  ipv4_count  = 1

  depends_on = [
    tencentcloud_vpc_ipv6_cidr_block.example,
    tencentcloud_vpc_ipv6_subnet_cidr_block.example
  ]
}

resource "tencentcloud_eni_ipv6_address" "example" {
  network_interface_id = tencentcloud_eni.example.id
  ipv6_address_count   = 1
}
```

Import

vpc eni ipv6 address can be imported using the id, e.g.

```
terraform import tencentcloud_eni_ipv6_address.example eni-fxrx5d1d
```