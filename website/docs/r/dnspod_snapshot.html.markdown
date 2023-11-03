---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_snapshot"
sidebar_current: "docs-tencentcloud-resource-dnspod_snapshot"
description: |-
  Provides a resource to create a dnspod snapshot
---

# tencentcloud_dnspod_snapshot

Provides a resource to create a dnspod snapshot

## Example Usage

```hcl
resource "tencentcloud_dnspod_snapshot" "snapshot" {
  domain = "dnspod.cn"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dnspod snapshot can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_snapshot.snapshot domain#snapshot_id
```

