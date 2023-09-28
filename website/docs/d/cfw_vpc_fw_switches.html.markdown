---
subcategory: "Cfw"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_vpc_fw_switches"
sidebar_current: "docs-tencentcloud-datasource-cfw_vpc_fw_switches"
description: |-
  Use this data source to query detailed information of cfw vpc_fw_switches
---

# tencentcloud_cfw_vpc_fw_switches

Use this data source to query detailed information of cfw vpc_fw_switches

## Example Usage

```hcl
data "tencentcloud_cfw_vpc_fw_switches" "example" {
  vpc_ins_id = "cfwg-c8c2de41"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_ins_id` - (Required, String) Firewall instance id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `switch_list` - Switch list.
  * `enable` - Switch status 0: off, 1: on.
  * `status` - Switch status 0: normal, 1: switching.
  * `switch_id` - Firewall switch ID.
  * `switch_mode` - switch mode.
  * `switch_name` - Firewall switch name.


