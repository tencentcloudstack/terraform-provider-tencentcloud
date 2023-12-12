Provide a resource to create a VPC.

Example Usage

Create a basic VPC

```hcl
resource "tencentcloud_vpc" "vpc" {
  name         = "tf-example"
  cidr_block   = "10.0.0.0/16"
  dns_servers  = ["119.29.29.29", "8.8.8.8"]
  is_multicast = false

  tags = {
    "test" = "test"
  }
}
```

Using Assistant CIDR

```hcl
resource "tencentcloud_vpc" "vpc" {
  name            = "tf-example"
  cidr_block      = "10.0.0.0/16"
  is_multicast    = false
  assistant_cidrs = ["172.16.0.0/24"]

  tags = {
    "test" = "test"
  }
}
```

Import

Vpc instance can be imported, e.g.

```
$ terraform import tencentcloud_vpc.test vpc-id
```