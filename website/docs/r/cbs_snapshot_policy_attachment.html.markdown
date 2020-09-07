---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_snapshot_policy_attachment"
sidebar_current: "docs-tencentcloud-resource-cbs_snapshot_policy_attachment"
description: |-
  Provides a CBS snapshot policy attachment resource.
---

# tencentcloud_cbs_snapshot_policy_attachment

Provides a CBS snapshot policy attachment resource.

## Example Usage

```hcl
resource "tencentcloud_cbs_snapshot_policy_attachment" "foo" {
  storage_id         = tencentcloud_cbs_storage.foo.id
  snapshot_policy_id = tencentcloud_cbs_snapshot_policy.policy.id
}
```

## Argument Reference

The following arguments are supported:

* `snapshot_policy_id` - (Required, ForceNew) ID of CBS snapshot policy.
* `storage_id` - (Required, ForceNew) ID of CBS.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



