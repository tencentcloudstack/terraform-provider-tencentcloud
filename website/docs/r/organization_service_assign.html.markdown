---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_service_assign"
sidebar_current: "docs-tencentcloud-resource-organization_service_assign"
description: |-
  Provides a resource to create a organization service assign
---

# tencentcloud_organization_service_assign

Provides a resource to create a organization service assign

## Example Usage

```hcl
resource "tencentcloud_organization_service_assign" "example" {
  service_id       = 15
  management_scope = 1
  member_uins      = [100037235241, 100033738111]
}
```



```hcl
resource "tencentcloud_organization_service_assign" "example" {
  service_id                = 15
  management_scope          = 2
  member_uins               = [100013415241, 100078908111]
  management_scope_uins     = [100019287759, 100020537485]
  management_scope_node_ids = [2024256, 2024259]
}
```

## Argument Reference

The following arguments are supported:

* `member_uins` - (Required, List: [`Int`]) Uin list of the delegated admins, Including up to 20 items.
* `service_id` - (Required, Int) Organization service ID.
* `management_scope_node_ids` - (Optional, List: [`Int`]) ID list of the managed departments. This parameter is valid when `management_scope` is `2`.
* `management_scope_uins` - (Optional, List: [`Int`]) Uin list of the managed members. This parameter is valid when `management_scope` is `2`.
* `management_scope` - (Optional, Int) Management scope of the delegated admin. Valid values: 1 (all members), 2 (partial members). Default value: `1`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

organization service assign can be imported using the id, e.g.
```
$ terraform import tencentcloud_organization_service_assign.example 15
```

