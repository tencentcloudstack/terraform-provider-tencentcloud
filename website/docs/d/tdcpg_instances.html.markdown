---
subcategory: "TDSQL-C for PostgreSQL(TDCPG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdcpg_instances"
sidebar_current: "docs-tencentcloud-datasource-tdcpg_instances"
description: |-
  Use this data source to query detailed information of tdcpg instances.
---

# tencentcloud_tdcpg_instances

Use this data source to query detailed information of tdcpg instances.

~> **NOTE:** This data source is still in internal testing. To experience its functions, you need to apply for a whitelist from Tencent Cloud.

## Example Usage

```hcl
data "tencentcloud_tdcpg_instances" "instances" {
  cluster_id    = ""
  instance_id   = ""
  instance_name = ""
  status        = ""
  instance_type = ""
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) instance id.
* `instance_id` - (Optional, String) instance id.
* `instance_name` - (Optional, String) instance name.
* `instance_type` - (Optional, String) instance type.
* `result_output_file` - (Optional, String) Used to save results.
* `status` - (Optional, String) instance status.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - instance list.
  * `cluster_id` - cluster id.
  * `cpu` - cpu cores.
  * `create_time` - create time.
  * `db_kernel_version` - db kernel version.
  * `db_major_version` - db major version.
  * `db_version` - db version.
  * `endpoint_id` - endpoint id.
  * `instance_id` - instance id.
  * `instance_name` - instance name.
  * `instance_type` - instance type.
  * `memory` - memory size, unit is GiB.
  * `pay_mode` - pay mode.
  * `pay_period_end_time` - pay period expired time.
  * `region` - region.
  * `status_desc` - status description.
  * `status` - status.
  * `zone` - zone.


