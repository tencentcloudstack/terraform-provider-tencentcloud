---
subcategory: "Cfw"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_edge_fw_switches"
sidebar_current: "docs-tencentcloud-datasource-cfw_edge_fw_switches"
description: |-
  Use this data source to query detailed information of cfw edge_fw_switches
---

# tencentcloud_cfw_edge_fw_switches

Use this data source to query detailed information of cfw edge_fw_switches

## Example Usage

```hcl
data "tencentcloud_cfw_edge_fw_switches" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Ip switch list.
  * `asset_type` - Asset Type.
  * `instance_id` - Instance Id.
  * `instance_name` - Instance Name.
  * `intranet_ip` - Intranet Ip.
  * `public_ip_type` - Public IP type.
  * `public_ip` - public ip.
  * `region` - region.
  * `status` - status.
  * `switch_mode` - switch mode.
  * `vpc_id` - vpc id.


