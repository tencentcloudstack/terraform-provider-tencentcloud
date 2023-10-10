---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_member_auth_identity_attachment"
sidebar_current: "docs-tencentcloud-resource-organization_org_member_auth_identity_attachment"
description: |-
  Provides a resource to create a organization org_member_auth_identity
---

# tencentcloud_organization_org_member_auth_identity_attachment

Provides a resource to create a organization org_member_auth_identity

## Example Usage

```hcl
resource "tencentcloud_organization_org_member_auth_identity_attachment" "org_member_auth_identity" {
  member_uin   = 100033704327
  identity_ids = [1657]
}
```

## Argument Reference

The following arguments are supported:

* `identity_ids` - (Required, Set: [`Int`], ForceNew) Identity Id list. Up to 5.
* `member_uin` - (Required, Int, ForceNew) Member Uin.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

organization org_member_auth_identity can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_member_auth_identity.org_member_auth_identity org_member_auth_identity_id
```

