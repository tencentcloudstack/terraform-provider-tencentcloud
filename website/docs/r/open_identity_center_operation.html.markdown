---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_open_identity_center_operation"
sidebar_current: "docs-tencentcloud-resource-open_identity_center_operation"
description: |-
  Provides a resource to open identity center
---

# tencentcloud_open_identity_center_operation

Provides a resource to open identity center

## Example Usage

```hcl
resource "tencentcloud_open_identity_center_operation" "open_identity_center_operation" {
  zone_name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `zone_name` - (Required, String, ForceNew) Space name, which must be globally unique and contain 2-64 characters including lowercase letters, digits, and hyphens (-). It can neither start or end with a hyphen (-) nor contain two consecutive hyphens (-).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `zone_id` - Space ID. z-Prefix starts with 12 random numbers/lowercase letters followed by.


