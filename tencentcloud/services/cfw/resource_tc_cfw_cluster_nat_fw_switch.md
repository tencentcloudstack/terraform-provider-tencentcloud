Provides a resource to create a CFW cluster NAT firewall switch

Example Usage

```hcl
resource "tencentcloud_cfw_cluster_nat_fw_switch" "example" {
  nat_ccn_switch {
    nat_ins_id   = "cfwnat-xxxxxxxx"
    ccn_id       = "ccn-xxxxxxxx"
    switch_mode  = 1
    routing_mode = 1
  }
}
```

With access instance list

```hcl
resource "tencentcloud_cfw_cluster_nat_fw_switch" "example" {
  nat_ccn_switch {
    nat_ins_id   = "cfwnat-xxxxxxxx"
    ccn_id       = "ccn-xxxxxxxx"
    switch_mode  = 2
    routing_mode = 0
    access_instance_list {
      instance_id      = "vpc-xxxxxxxx"
      instance_type    = "VPC"
      instance_region  = "ap-guangzhou"
      access_cidr_mode = 1
      access_cidr_list = ["10.0.0.0/16"]
    }
  }
}
```

Import

CFW cluster NAT firewall switch can be imported using the composite id `nat_ins_id#ccn_id`, e.g.

```
terraform import tencentcloud_cfw_cluster_nat_fw_switch.example cfwnat-xxxxxxxx#ccn-xxxxxxxx
```
