Provides a resource to creating dedicated tunnels instances.

~> **NOTE:** 1. ID of the DC is queried, can only apply for this resource offline.

Example Usage

If network_type is VPC

```hcl
resource "tencentcloud_dcx" "example" {
  dc_id            = "dc-ink7y3qf"
  name             = "tf-example"
  dc_owner_account = "100017971194"
  network_type     = "VPC"
  network_region   = "ap-guangzhou"
  vpc_id           = "vpc-nzuu8dyj"
  dcg_id           = "dcg-ehr22qfb"
  bandwidth        = 100
  route_type       = "BGP"
  bgp_asn          = 64511
  vlan             = 60
  tencent_address  = "10.8.254.14/30"
  customer_address = "10.8.254.13/30"
}
```

If network_type is CCN

```hcl
resource "tencentcloud_dcx" "example" {
  dc_id            = "dc-ink7y3qf"
  name             = "tf-example"
  dc_owner_account = "100017971194"
  network_type     = "CCN"
  network_region   = "ap-guangzhou"
  dcg_id           = "dcg-6d4uaubp"
  bandwidth        = 100
  route_type       = "BGP"
  bgp_asn          = 64511
  vlan             = 10
  tencent_address  = "10.8.254.10/30"
  customer_address = "10.8.254.9/30"
}
```

Import

DCX instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_dcx.example dcx-cbbr1gjk
```