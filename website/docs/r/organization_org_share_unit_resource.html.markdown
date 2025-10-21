---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_share_unit_resource"
sidebar_current: "docs-tencentcloud-resource-organization_org_share_unit_resource"
description: |-
  Provides a resource to create a organization organization_org_share_unit_resource
---

# tencentcloud_organization_org_share_unit_resource

Provides a resource to create a organization organization_org_share_unit_resource

## Example Usage

```hcl
resource "tencentcloud_organization_org_share_unit_resource" "organization_org_share_unit_resource" {
  unit_id             = "xxxxxx"
  area                = "ap-guangzhou"
  type                = "secret"
  product_resource_id = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `area` - (Required, String, ForceNew) Shared unit area.
* `product_resource_id` - (Required, String, ForceNew) Product Resource ID.
* `type` - (Required, String, ForceNew) Shared resource type.
* `unit_id` - (Required, String, ForceNew) Shared unit ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `resource_id` - Shared resource ID.
* `share_manager_uin` - Sharing administrator OwnerUin.
* `shared_member_num` - Number of shared unit members.
* `shared_member_use_num` - Number of shared unit members in use.


## Import

organization organization_org_share_unit_resource can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_share_unit_resource.organization_org_share_unit_resource ${unit_id}#${area}#${share_resource_type}#${product_resource_id}
```

