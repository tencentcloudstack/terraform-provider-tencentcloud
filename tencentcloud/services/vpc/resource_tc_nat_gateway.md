Provides a resource to create a NAT gateway.

Example Usage

Create a traditional NAT gateway.

```hcl
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_nat_gateway_vpc"
}

resource "tencentcloud_eip" "eip_example1" {
  name = "tf_nat_gateway_eip1"
}

resource "tencentcloud_eip" "eip_example2" {
  name = "tf_nat_gateway_eip2"
}

resource "tencentcloud_nat_gateway" "example" {
  name             = "tf_example_nat_gateway"
  vpc_id           = tencentcloud_vpc.vpc.id
  bandwidth        = 100
  max_concurrent   = 1000000
  assigned_eip_set = [
    tencentcloud_eip.eip_example1.public_ip,
    tencentcloud_eip.eip_example2.public_ip,
  ]
  tags = {
    tf_tag_key = "tf_tag_value"
  }
}
```

Create a standard NAT gateway.

```hcl
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_nat_gateway_vpc"
}

resource "tencentcloud_eip" "eip_example1" {
  name = "tf_nat_gateway_eip1"
}

resource "tencentcloud_eip" "eip_example2" {
  name = "tf_nat_gateway_eip2"
}

resource "tencentcloud_nat_gateway" "example" {
  name             = "tf_example_nat_gateway"
  vpc_id           = tencentcloud_vpc.vpc.id
  assigned_eip_set = [
    tencentcloud_eip.eip_example1.public_ip,
    tencentcloud_eip.eip_example2.public_ip,
  ]
  nat_product_version = 2
  tags                = {
    tf_tag_key = "tf_tag_value"
  }
  lifecycle {
    ignore_changes = [
      // standard nat will set default values for bandwidth and max_concurrent
      bandwidth,
      max_concurrent,
    ]
  }
}
```

Import

NAT gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_nat_gateway.foo nat-1asg3t63
```