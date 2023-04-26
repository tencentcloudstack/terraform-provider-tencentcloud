---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_apply_instance_snapshot"
sidebar_current: "docs-tencentcloud-resource-lighthouse_apply_instance_snapshot"
description: |-
  Provides a resource to create a lighthouse apply_instance_snapshot
---

# tencentcloud_lighthouse_apply_instance_snapshot

Provides a resource to create a lighthouse apply_instance_snapshot

## Example Usage

```hcl
resource "tencentcloud_lighthouse_apply_instance_snapshot" "apply_instance_snapshot" {
  instance_id = "lhins-123456"
  snapshot_id = "lhsnap-123456"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `snapshot_id` - (Required, String, ForceNew) Snapshot ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



