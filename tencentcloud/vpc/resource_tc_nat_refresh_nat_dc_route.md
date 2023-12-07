Provides a resource to create a vpc refresh_nat_dc_route

Example Usage

If `dry_run` is True

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_eip" "eip_example" {
  name = "eip_example"
}

resource "tencentcloud_nat_gateway" "nat" {
  vpc_id         = tencentcloud_vpc.vpc.id
  name           = "tf_example_nat_gateway"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_example.public_ip,
  ]
}

resource "tencentcloud_nat_refresh_nat_dc_route" "refresh_nat_dc_route" {
  nat_gateway_id = tencentcloud_nat_gateway.nat.id
  vpc_id         = tencentcloud_vpc.vpc.id
  dry_run        = true
}
```

Or `dry_run` is False

```hcl
resource "tencentcloud_nat_refresh_nat_dc_route" "refresh_nat_dc_route" {
  nat_gateway_id = tencentcloud_nat_gateway.nat.id
  vpc_id         = tencentcloud_vpc.vpc.id
  dry_run        = false
}
```

Import

vpc refresh_nat_dc_route can be imported using the id, e.g.

```
terraform import tencentcloud_nat_refresh_nat_dc_route.refresh_nat_dc_route vpc_id#nat_gateway_id
```