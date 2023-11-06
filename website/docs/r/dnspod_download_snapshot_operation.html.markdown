---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_download_snapshot_operation"
sidebar_current: "docs-tencentcloud-resource-dnspod_download_snapshot_operation"
description: |-
  Provides a resource to create a dnspod download_snapshot
---

# tencentcloud_dnspod_download_snapshot_operation

Provides a resource to create a dnspod download_snapshot

## Example Usage

```hcl
resource "tencentcloud_dnspod_download_snapshot_operation" "download_snapshot" {
  domain      = "dnspod.cn"
  snapshot_id = "456"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain.
* `snapshot_id` - (Required, String, ForceNew) Snapshot ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cos_url` - Snapshot download url.


