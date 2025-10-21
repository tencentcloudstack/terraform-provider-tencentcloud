---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_share_unit_member"
sidebar_current: "docs-tencentcloud-resource-organization_org_share_unit_member"
description: |-
  Provides a resource to create a Organization share unit member
---

# tencentcloud_organization_org_share_unit_member

Provides a resource to create a Organization share unit member

~> **NOTE:** This resource has been deprecated in Terraform TencentCloud provider version 1.82.28, Please use `tencentcloud_organization_org_share_unit_member_v2` instead.

## Example Usage

```hcl
resource "tencentcloud_organization_org_share_unit" "example" {
  name        = "tf-example"
  area        = "ap-guangzhou"
  description = "description."
}

resource "tencentcloud_organization_org_share_unit_member" "example" {
  unit_id = tencentcloud_organization_org_share_unit.example.unit_id
  area    = tencentcloud_organization_org_share_unit.example.area
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



