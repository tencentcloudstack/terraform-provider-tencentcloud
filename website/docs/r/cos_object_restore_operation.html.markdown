---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_object_restore_operation"
sidebar_current: "docs-tencentcloud-resource-cos_object_restore_operation"
description: |-
  Provides a resource to restore object
---

# tencentcloud_cos_object_restore_operation

Provides a resource to restore object

## Example Usage

```hcl
resource "tencentcloud_cos_object_restore_operation" "object_restore" {
  bucket = "keep-test-1308919341"
  key    = "test-restore.txt"
  tier   = "Expedited"
  days   = 2
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) Bucket.
* `days` - (Required, Int, ForceNew) Specifies the valid duration of the restored temporary copy in days.
* `key` - (Required, String, ForceNew) Object key.
* `tier` - (Required, String, ForceNew) when restoring, Tier can be specified as the supported recovery model.
There are three recovery models for recovering archived storage type data, which are:
- Expedited: quick retrieval mode, and the recovery task can be completed in 1-5 minutes.
- Standard: standard retrieval mode. Recovery task is completed within 3-5 hours.
- Bulk: batch retrieval mode, and the recovery task is completed within 5-12 hours.
For deep recovery archive storage type data, there are two recovery models, which are:
- Standard: standard retrieval mode, recovery time is 12-24 hours.
- Bulk: batch retrieval mode, recovery time is 24-48 hours.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



