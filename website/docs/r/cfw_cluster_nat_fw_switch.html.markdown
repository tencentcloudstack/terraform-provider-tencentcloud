---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_cluster_nat_fw_switch"
sidebar_current: "docs-tencentcloud-resource-cfw_cluster_nat_fw_switch"
description: |-
  Provides a resource to create a CFW cluster NAT firewall switch
---

# tencentcloud_cfw_cluster_nat_fw_switch

Provides a resource to create a CFW cluster NAT firewall switch

## Example Usage

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

### With access instance list

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

## Argument Reference

The following arguments are supported:

* `nat_ccn_switch` - (Required, List) NAT CCN switch configuration.

The `access_instance_list` object of `nat_ccn_switch` supports the following:

* `access_cidr_mode` - (Required, Int) Network segment mode for accessing firewall: 0-no access, 1-access all network segments associated with the instance, 2-access user-defined network segments.
* `instance_id` - (Required, String) Instance ID.
* `instance_region` - (Required, String) Region where the instance is located.
* `instance_type` - (Required, String) Instance type such as VPC or DIRECTCONNECT.
* `access_cidr_list` - (Optional, List) List of network segments for accessing firewall.

The `nat_ccn_switch` object supports the following:

* `ccn_id` - (Required, String, ForceNew) CCN instance ID.
* `nat_ins_id` - (Required, String, ForceNew) NAT firewall instance ID.
* `switch_mode` - (Required, Int) Switch access mode, 1: automatic access, 2: manual access.
* `access_instance_list` - (Optional, List) List of access instances.
* `lead_vpc_cidr` - (Optional, String) CIDR of the lead VPC.
* `routing_mode` - (Optional, Int) Traffic steering routing method, 0: multi-route table, 1: policy routing. Automatic access mode only supports policy routing (1); manual access mode supports both multi-route table (0) and policy routing (1).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `bypass` - Bypass status, 0: disabled, 1: enabled.
* `ccn_id` - CCN instance ID.
* `routing_mode` - Traffic steering routing method, 0: multi-route table, 1: policy routing.
* `switch_mode` - Switch access mode, 1: automatic access, 2: manual access.


## Import

CFW cluster NAT firewall switch can be imported using the composite id `nat_ins_id#ccn_id`, e.g.

```
terraform import tencentcloud_cfw_cluster_nat_fw_switch.example cfwnat-xxxxxxxx#ccn-xxxxxxxx
```

