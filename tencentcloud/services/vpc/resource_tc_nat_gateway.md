Provides a resource to create a NAT gateway.

~> **NOTE:** If `nat_product_version` is `1`, `max_concurrent` valid values is `1000000`, `3000000`, `10000000`.

~> **NOTE:** If set `stock_public_ip_addresses_bandwidth_out`, do not set the `internet_max_bandwidth_out` parameter of resource `tencentcloud_eip` at the same time, otherwise conflicts may occur.

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
  name                = "tf_example_nat_gateway"
  vpc_id              = tencentcloud_vpc.vpc.id
  nat_product_version = 1
  bandwidth           = 100
  max_concurrent      = 1000000
  assigned_eip_set = [
    tencentcloud_eip.eip_example1.public_ip,
    tencentcloud_eip.eip_example2.public_ip,
  ]
  tags = {
    createBy = "terraform"
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
  name                = "tf_example_nat_gateway"
  vpc_id              = tencentcloud_vpc.vpc.id
  nat_product_version = 2
  assigned_eip_set = [
    tencentcloud_eip.eip_example1.public_ip,
    tencentcloud_eip.eip_example2.public_ip,
  ]
  tags = {
    createBy = "terraform"
  }
}
```

Or set stock public ip addresses bandwidth out

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
  name                                    = "tf_example_nat_gateway"
  vpc_id                                  = tencentcloud_vpc.vpc.id
  nat_product_version                     = 2
  stock_public_ip_addresses_bandwidth_out = 100
  assigned_eip_set = [
    tencentcloud_eip.eip_example1.public_ip,
    tencentcloud_eip.eip_example2.public_ip,
  ]
  tags = {
    createBy = "terraform"
  }
}
```

Import

NAT gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_nat_gateway.example nat-1asg3t63
```
