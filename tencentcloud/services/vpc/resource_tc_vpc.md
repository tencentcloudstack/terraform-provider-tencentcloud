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
    createBy = "Terraform"
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
    createBy = "Terraform"
  }
}
```

Enable route vpc publish

```hcl
resource "tencentcloud_vpc" "vpc" {
  name                          = "tf-example"
  cidr_block                    = "10.0.0.0/16"
  dns_servers                   = ["119.29.29.29", "8.8.8.8"]
  is_multicast                  = false
  enable_route_vpc_publish      = true

  tags = {
    createBy = "Terraform"
  }
}
```

Import

Vpc instance can be imported, e.g.

```
$ terraform import tencentcloud_vpc.vpc vpc-8vazrwjv
```
