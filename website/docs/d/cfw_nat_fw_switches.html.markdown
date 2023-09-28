---
subcategory: "Cfw"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_nat_fw_switches"
sidebar_current: "docs-tencentcloud-datasource-cfw_nat_fw_switches"
description: |-
  Use this data source to query detailed information of cfw nat_fw_switches
---

# tencentcloud_cfw_nat_fw_switches

Use this data source to query detailed information of cfw nat_fw_switches

## Example Usage

### Query Nat instance'switch by instance id

```hcl
data "tencentcloud_cfw_nat_fw_switches" "example" {
  nat_ins_id = "cfwnat-18d2ba18"
}
```

### Or filter by switch status

```hcl
data "tencentcloud_cfw_nat_fw_switches" "example" {
  nat_ins_id = "cfwnat-18d2ba18"
  status     = 1
}
```

## Argument Reference

The following arguments are supported:

* `nat_ins_id` - (Optional, String) Filter the NAT firewall instance to which the NAT firewall subnet switch belongs.
* `result_output_file` - (Optional, String) Used to save results.
* `status` - (Optional, Int) Switch status, 1 open; 0 close.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - NAT border firewall switch list data.
  * `abnormal` - Whether the switch is abnormal, 0: normal, 1: abnormal.
  * `cvm_num` - Cvm Num.
  * `enable` - Effective status.
  * `id` - ID.
  * `nat_id` - NAT gatway Id.
  * `nat_ins_id` - NAT firewall instance Id.
  * `nat_ins_name` - NAT firewall instance name.
  * `nat_name` - NAT gatway name.
  * `region` - Region.
  * `route_id` - Route Id.
  * `route_name` - Route Name.
  * `status` - Switch status.
  * `subnet_cidr` - IPv4 CIDR.
  * `subnet_id` - Subnet Id.
  * `subnet_name` - Subnet Name.
  * `vpc_id` - Vpc Id.
  * `vpc_name` - Vpc Name.


