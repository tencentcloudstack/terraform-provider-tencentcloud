---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_share_unit_member"
sidebar_current: "docs-tencentcloud-resource-organization_org_share_unit_member"
description: |-
  Provides a resource to create a organization org_share_unit_member
---

# tencentcloud_organization_org_share_unit_member

Provides a resource to create a organization org_share_unit_member

## Example Usage

```hcl
resource "tencentcloud_organization_org_share_unit" "org_share_unit" {
  name        = "iac-test"
  area        = "ap-guangzhou"
  description = "iac-test"
}
resource "tencentcloud_organization_org_share_unit_member" "org_share_unit_member" {
  unit_id = tencentcloud_organization_org_share_unit.org_share_unit.unit_id
  area    = tencentcloud_organization_org_share_unit.org_share_unit.area
  members {
    share_member_uin = 100035309479
  }
}
```

## Argument Reference

The following arguments are supported:

* `area` - (Required, String, ForceNew) Shared unit region.
* `members` - (Required, List, ForceNew) Shared member list. Up to 10 items.
* `unit_id` - (Required, String, ForceNew) Shared unit ID.

The `members` object supports the following:

* `share_member_uin` - (Required, Int) Member uin.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



