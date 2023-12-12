Provides a resource to create a vpc ipv6_cidr_block

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

Import

vpc ipv6_cidr_block can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_ipv6_cidr_block.ipv6_cidr_block vpc_id
```