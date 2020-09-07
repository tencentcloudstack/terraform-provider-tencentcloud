---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_snapshot_policies"
sidebar_current: "docs-tencentcloud-datasource-cbs_snapshot_policies"
description: |-
  Use this data source to query detailed information of CBS snapshot policies.
---

# tencentcloud_cbs_snapshot_policies

Use this data source to query detailed information of CBS snapshot policies.

## Example Usage

```hcl
data "tencentcloud_cbs_snapshot_policies" "policies" {
  snapshot_policy_id   = "snap-f3io7adt"
  snapshot_policy_name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional) Used to save results.
* `snapshot_policy_id` - (Optional) ID of the snapshot policy to be queried.
* `snapshot_policy_name` - (Optional) Name of the snapshot policy to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `snapshot_policy_list` - A list of snapshot policy. Each element contains the following attributes:
  * `attached_storage_ids` - Storage ids that the snapshot policy attached.
  * `create_time` - Create time of the snapshot policy.
  * `repeat_hours` - Trigger hours of periodic snapshot.
  * `repeat_weekdays` - Trigger days of periodic snapshot.
  * `retention_days` - Retention days of the snapshot.
  * `snapshot_policy_id` - ID of the snapshot policy.
  * `snapshot_policy_name` - Name of the snapshot policy.
  * `status` - Status of the snapshot policy.


