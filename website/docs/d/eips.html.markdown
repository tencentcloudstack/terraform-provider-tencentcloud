---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eips"
sidebar_current: "docs-tencentcloud-datasource-eips"
description: |-
  Use this data source to query eip instances.
---

# tencentcloud_eips

Use this data source to query eip instances.

## Example Usage

```hcl
data "tencentcloud_eips" "foo" {
  eip_id = "eip-ry9h95hg"
}
```

## Argument Reference

The following arguments are supported:

* `eip_id` - (Optional) ID of the eip to be queried.
* `eip_name` - (Optional) Name of the eip to be queried.
* `public_ip` - (Optional) The elastic ip address.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `eip_list` - An information list of eip. Each element contains the following attributes:
  * `create_time` - Creation time of the eip.
  * `eip_id` - ID of the eip.
  * `eip_name` - Name of the eip.
  * `eip_type` - Type of the eip.
  * `eni_id` - The eni id to bind with the eip.
  * `instance_id` - The instance id to bind with the eip.
  * `public_ip` - The elastic ip address.
  * `status` - The eip current status.


