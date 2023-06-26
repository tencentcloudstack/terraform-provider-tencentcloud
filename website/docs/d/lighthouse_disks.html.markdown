---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_disks"
sidebar_current: "docs-tencentcloud-datasource-lighthouse_disks"
description: |-
  Use this data source to query detailed information of lighthouse disk
---

# tencentcloud_lighthouse_disks

Use this data source to query detailed information of lighthouse disk

## Example Usage

```hcl
data "tencentcloud_lighthouse_disks" "disks" {
  disk_ids = ["lhdisk-xxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `disk_ids` - (Optional, Set: [`String`]) List of disk ids.
* `filters` - (Optional, List) Filter list.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Fields to be filtered. Valid names: `disk-id`: Filters by disk id; `instance-id`: Filter by instance id; `disk-name`: Filter by disk name; `zone`: Filter by zone; `disk-usage`: Filter by disk usage(Values: `SYSTEM_DISK` or `DATA_DISK`); `disk-state`: Filter by disk state.
* `values` - (Required, Set) Value of the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `disk_list` - Cloud disk information list.
  * `attached` - Disk attach state.
  * `created_time` - Created time. Expressed according to the ISO8601 standard, and using UTC time. The format is `YYYY-MM-DDThh:mm:ssZ`.
  * `delete_with_instance` - Whether to release with the instance.
  * `disk_backup_count` - Number of existing backup points of cloud disk.
  * `disk_backup_quota` - Number of backup points quota for cloud disk.
  * `disk_charge_type` - Disk charge type.
  * `disk_id` - Disk id.
  * `disk_name` - Disk name.
  * `disk_size` - Disk size.
  * `disk_state` - Disk state. Valid values:`PENDING`, `UNATTACHED`, `ATTACHING`, `ATTACHED`, `DETACHING`, `SHUTDOWN`, `CREATED_FAILED`, `TERMINATING`, `DELETING`, `FREEZING`.
  * `disk_type` - Disk type.
  * `disk_usage` - Disk usage.
  * `expired_time` - Expired time. Expressed according to the ISO8601 standard, and using UTC time. The format is `YYYY-MM-DDThh:mm:ssZ`.
  * `instance_id` - Instance id.
  * `isolated_time` - Isolated time. Expressed according to the ISO8601 standard, and using UTC time. The format is `YYYY-MM-DDThh:mm:ssZ`.
  * `latest_operation_request_id` - Latest operation request id.
  * `latest_operation_state` - Latest operation state.
  * `latest_operation` - Latest operation.
  * `renew_flag` - Renew flag.
  * `zone` - Availability zone.


