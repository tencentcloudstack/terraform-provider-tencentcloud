---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_share_unit_member_v2"
sidebar_current: "docs-tencentcloud-resource-organization_org_share_unit_member_v2"
description: |-
  Provides a resource to create a Organization share unit member
---

# tencentcloud_organization_org_share_unit_member_v2

Provides a resource to create a Organization share unit member

~> **NOTE:** This resource must exclusive in one share unit, do not declare additional members resources of this share unit elsewhere.

## Example Usage

```hcl
resource "tencentcloud_organization_org_share_unit" "example" {
  name        = "tf-example"
  area        = "ap-guangzhou"
  description = "description."
}

resource "tencentcloud_organization_org_share_unit_member_v2" "example" {
  unit_id = tencentcloud_organization_org_share_unit.example.unit_id
  area    = tencentcloud_organization_org_share_unit.example.area
  members {
    share_member_uin = 100042257812
  }

  members {
    share_member_uin = 100043990767
  }

  members {
    share_member_uin = 100042234123
  }
}
```

## Argument Reference

The following arguments are supported:

* `area` - (Required, String, ForceNew) Shared unit region.
* `members` - (Required, Set) Shared member list.
* `unit_id` - (Required, String, ForceNew) Shared unit ID.

The `members` object supports the following:

* `share_member_uin` - (Required, Int) Member uin.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AddShareUnitMembers](https://cloud.tencent.com/document/api/850/103055) - 添加共享单元成员
* Read: [DescribeShareUnitMembers](https://cloud.tencent.com/document/api/850/103049) - 获取共享单元成员列表
* Update: [AddShareUnitMembers](https://cloud.tencent.com/document/api/850/103055) - 添加共享单元成员
* Update: [DeleteShareUnitMembers](https://cloud.tencent.com/document/api/850/103052) - 删除共享单元成员
* Delete: [DeleteShareUnitMembers](https://cloud.tencent.com/document/api/850/103052) - 删除共享单元成员
* Delete: [DescribeShareUnitMembers](https://cloud.tencent.com/document/api/850/103049) - 获取共享单元成员列表

## Import

Organization share unit member can be imported using the unitId#area, e.g.

```
terraform import tencentcloud_organization_org_share_unit_member_v2.example shareUnit-switt8i4s4#ap-guangzhou
```

