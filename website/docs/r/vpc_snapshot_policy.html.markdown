---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_snapshot_policy"
sidebar_current: "docs-tencentcloud-resource-vpc_snapshot_policy"
description: |-
  Provides a resource to create a vpc snapshot_policy
---

# tencentcloud_vpc_snapshot_policy

Provides a resource to create a vpc snapshot_policy

## Example Usage

```hcl
resource "tencentcloud_vpc_snapshot_policy" "snapshot_policy" {
  snapshot_policy_name = "terraform-test"
  backup_type          = "time"
  cos_bucket           = "cos-lock-1308919341"
  cos_region           = "ap-guangzhou"
  create_new_cos       = false
  keep_time            = 2

  backup_policies {
    backup_day  = "monday"
    backup_time = "00:00:00"
  }
  backup_policies {
    backup_day  = "tuesday"
    backup_time = "02:03:03"
  }
  backup_policies {
    backup_day  = "wednesday"
    backup_time = "04:13:23"
  }
}
```

## Argument Reference

The following arguments are supported:

* `backup_type` - (Required, String) Backup strategy type, `operate`: operate backup, `time`: schedule backup.
* `cos_bucket` - (Required, String) cos bucket.
* `cos_region` - (Required, String) The region where the cos bucket is located.
* `create_new_cos` - (Required, Bool) Whether to create a new cos bucket, the default is False.Note: This field may return null, indicating that no valid value can be obtained.
* `keep_time` - (Required, Int) The retention time supports 1 to 365 days.
* `snapshot_policy_name` - (Required, String) Snapshot policy name.
* `backup_policies` - (Optional, List) Time backup strategy. Note: This field may return null, indicating that no valid value can be obtained.

The `backup_policies` object supports the following:

* `backup_day` - (Required, String) Backup cycle time, the value can be monday, tuesday, wednesday, thursday, friday, saturday, sunday.
* `backup_time` - (Required, String) Backup time point, format:HH:mm:ss.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time.Note: This field may return null, indicating that no valid value can be obtained.
* `enable` - Enabled state, True-enabled, False-disabled, the default is True.
* `snapshot_policy_id` - Snapshot policy Id.


## Import

vpc snapshot_policy can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_snapshot_policy.snapshot_policy snapshot_policy_id
```

