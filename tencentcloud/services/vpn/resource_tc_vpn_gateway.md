Provides a resource to create a VPN gateway.

-> **NOTE:** The prepaid VPN gateway do not support renew operation or delete operation with terraform.

Example Usage

VPC SSL VPN gateway

```hcl
resource "tencentcloud_vpn_gateway" "example" {
  name      = "tf-example"
  bandwidth = 5
  zone      = "ap-guangzhou-3"
  type      = "SSL"
  vpc_id    = "vpc-86v957zb"

  tags = {
    createBy = "Terraform"
  }
}
```

CCN IPSEC VPN gateway

```hcl
resource "tencentcloud_vpn_gateway" "example" {
  name      = "tf-example"
  bandwidth = 5
  zone      = "ap-guangzhou-3"
  type      = "IPSEC"

  tags = {
    createBy = "Terraform"
  }
}
```

CCN SSL VPN gateway

```hcl
resource "tencentcloud_vpn_gateway" "example" {
  name      = "tf-example"
  bandwidth = 5
  zone      = "ap-guangzhou-3"
  type      = "SSL_CCN"

  tags = {
    createBy = "Terraform"
  }
}
```

CCN VPN gateway

```hcl
resource "tencentcloud_vpn_gateway" "example" {
  name      = "tf-example"
  bandwidth = 200
  type      = "CCN"
  bgp_asn   = 9000

  tags = {
    createBy = "Terraform"
  }
}
```

POSTPAID_BY_HOUR VPN gateway

```hcl
resource "tencentcloud_vpn_gateway" "example" {
  name      = "tf-example"
  vpc_id    = "vpc-dk8zmwuf"
  bandwidth = 5
  zone      = "ap-guangzhou-3"

  tags = {
    createBy = "Terraform"
  }
}
```

PREPAID VPN gateway

```hcl
resource "tencentcloud_vpn_gateway" "example" {
  name           = "tf-example"
  vpc_id         = "vpc-dk8zmwuf"
  bandwidth      = 5
  zone           = "ap-guangzhou-3"
  charge_type    = "PREPAID"
  prepaid_period = 1

  tags = {
    createBy = "Terraform"
  }
}
```

Import

VPN gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_gateway.example vpngw-8ccsnclt
```