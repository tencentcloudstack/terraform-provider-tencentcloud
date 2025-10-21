---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_snapshot"
sidebar_current: "docs-tencentcloud-resource-lighthouse_snapshot"
description: |-
  Provides a resource to create a lighthouse snapshot
---

# tencentcloud_lighthouse_snapshot

Provides a resource to create a lighthouse snapshot

## Example Usage

```hcl
resource "tencentcloud_lighthouse_snapshot" "snapshot" {
  instance_id   = "lhins-acd1234"
  snapshot_name = "snap_20200903"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of the instance for which to create a snapshot.
* `snapshot_name` - (Optional, String) Snapshot name, which can contain up to 60 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



