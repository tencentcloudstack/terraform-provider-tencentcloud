---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_accept_join_share_unit_invitation_operation"
sidebar_current: "docs-tencentcloud-resource-accept_join_share_unit_invitation_operation"
description: |-
  Provides a resource to create a organization accept_join_share_unit_invitation_operation
---

# tencentcloud_accept_join_share_unit_invitation_operation

Provides a resource to create a organization accept_join_share_unit_invitation_operation

## Example Usage

```hcl
resource "tencentcloud_accept_join_share_unit_invitation_operation" "accept_join_share_unit_invitation_operation" {
  unit_id = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `unit_id` - (Required, String, ForceNew) Shared unit ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



