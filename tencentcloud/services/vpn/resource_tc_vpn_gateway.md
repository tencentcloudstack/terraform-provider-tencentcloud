Provides a resource to create a VPN gateway.

-> **NOTE:** The prepaid VPN gateway do not support renew operation or delete operation with terraform.

Example Usage

VPC SSL VPN gateway
```hcl
resource "tencentcloud_vpn_gateway" "my_cgw" {
  name      = "test"
  bandwidth = 5
  zone      = "ap-guangzhou-3"
  type      = "SSL"
  vpc_id    = "vpc-86v957zb"

  tags = {
    test = "test"
  }
}
```

CCN IPEC VPN gateway
```hcl
resource "tencentcloud_vpn_gateway" "my_cgw" {
  name      = "test"
  bandwidth = 5
  zone      = "ap-guangzhou-3"
  type      = "CCN"

  tags      = {
    test = "test"
  }
}
```

CCN SSL VPN gateway
```hcl
resource "tencentcloud_vpn_gateway" "my_cgw" {
  name      = "test"
  bandwidth = 5
  zone      = "ap-guangzhou-3"
  type      = "SSL_CCN"

  tags      = {
    test = "test"
  }
}
```

POSTPAID_BY_HOUR VPN gateway
```hcl
resource "tencentcloud_vpn_gateway" "my_cgw" {
  name      = "test"
  vpc_id    = "vpc-dk8zmwuf"
  bandwidth = 5
  zone      = "ap-guangzhou-3"

  tags = {
    test = "test"
  }
}
```

PREPAID VPN gateway
```hcl
resource "tencentcloud_vpn_gateway" "my_cgw" {
  name           = "test"
  vpc_id         = "vpc-dk8zmwuf"
  bandwidth      = 5
  zone           = "ap-guangzhou-3"
  charge_type    = "PREPAID"
  prepaid_period = 1

  tags = {
    test = "test"
  }
}
```

Import

VPN gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_gateway.foo vpngw-8ccsnclt
```