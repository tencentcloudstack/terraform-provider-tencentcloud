Provides a resource to creating dedicated tunnels instances.

~> **NOTE:** 1. ID of the DC is queried, can only apply for this resource offline.

Example Usage

```hcl
variable "dc_id" {
  default = "dc-kax48sg7"
}

variable "dcg_id" {
  default = "dcg-dmbhf7jf"
}

variable "vpc_id" {
  default = "vpc-4h9v4mo3"
}

resource "tencentcloud_dcx" "bgp_main" {
  bandwidth    = 900
  dc_id        = var.dc_id
  dcg_id       = var.dcg_id
  name         = "bgp_main"
  network_type = "VPC"
  route_type   = "BGP"
  vlan         = 306
  vpc_id       = var.vpc_id
}

resource "tencentcloud_dcx" "static_main" {
  bandwidth             = 900
  dc_id                 = var.dc_id
  dcg_id                = var.dcg_id
  name                  = "static_main"
  dc_owner_account      = "xxxxxxxx"
  network_type          = "VPC"
  route_type            = "STATIC"
  vlan                  = 301
  vpc_id                = var.vpc_id
  route_filter_prefixes = ["10.10.10.101/32"]
  tencent_address       = "100.93.46.1/30"
  customer_address      = "100.93.46.2/30"
}
```