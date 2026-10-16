---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_accept_join_share_unit_invitation_operation"
sidebar_current: "docs-tencentcloud-resource-accept_join_share_unit_invitation_operation"
description: |-
  Provides a resource to create a organization accept join share unit invitation operation
---

# tencentcloud_accept_join_share_unit_invitation_operation

Provides a resource to create a organization accept join share unit invitation operation

## Example Usage

```hcl
resource "tencentcloud_accept_join_share_unit_invitation_operation" "example" {
  unit_id = "shareUnit-xhreo**2p"
}
```

## Argument Reference

The following arguments are supported:

* `unit_id` - (Required, String, ForceNew) Shared unit ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AcceptJoinShareUnitInvitation](https://cloud.tencent.com/document/api/850/107559)

