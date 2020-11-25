---
subcategory: "TCR"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_instances"
sidebar_current: "docs-tencentcloud-datasource-tcr_instances"
description: |-
  Use this data source to query detailed information of TCR instances.
---

# tencentcloud_tcr_instances

Use this data source to query detailed information of TCR instances.

## Example Usage

```hcl
data "tencentcloud_tcr_instances" "name" {
  name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional) Id of the TCR instance to query.
* `name` - (Optional) Name of the TCR instance to query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - Information list of the dedicated TCR instances.
  * `id` - Id of the TCR instance.
  * `instance_type` - Instance type.
  * `internal_end_point` - Internal address for access of the TCR instance.
  * `name` - Name of TCR instance.
  * `public_domain` - Public address for access of the TCR instance.
  * `status` - Status of the TCR instance.
  * `tags` - Tags of the TCR instance.


