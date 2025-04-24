Provides a resource to create a VPC ipv6 cidr block

Example Usage

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_vpc_ipv6_cidr_block" "example" {
  vpc_id = tencentcloud_vpc.vpc.id
}
```

Or

```hcl
resource "tencentcloud_vpc_ipv6_cidr_block" "example" {
  vpc_id       = tencentcloud_vpc.vpc.id
  address_type = "ULA"
}
```

Import

vpc ipv6_cidr_block can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_ipv6_cidr_block.example vpc-826mi3hd
```
