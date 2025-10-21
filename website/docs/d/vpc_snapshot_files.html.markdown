---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_snapshot_files"
sidebar_current: "docs-tencentcloud-datasource-vpc_snapshot_files"
description: |-
  Use this data source to query detailed information of vpc snapshot_files
---

# tencentcloud_vpc_snapshot_files

Use this data source to query detailed information of vpc snapshot_files

## Example Usage

```hcl
data "tencentcloud_vpc_snapshot_files" "snapshot_files" {
  business_type = "securitygroup"
  instance_id   = "sg-902tl7t7"
  start_date    = "2022-10-10 00:00:00"
  end_date      = "2023-10-30 19:00:00"
}
```

## Argument Reference

The following arguments are supported:

* `business_type` - (Required, String) Business type, currently supports security group:securitygroup.
* `end_date` - (Required, String) End date in the format %Y-%m-%d %H:%M:%S.
* `instance_id` - (Required, String) InstanceId.
* `start_date` - (Required, String) Start date in the format %Y-%m-%d %H:%M:%S.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `snapshot_file_set` - snap shot file set.
  * `backup_time` - backup time.
  * `instance_id` - instance id.
  * `operator` - Uin of operator.
  * `snapshot_file_id` - snap shot file id.
  * `snapshot_policy_id` - Snapshot Policy Id.


