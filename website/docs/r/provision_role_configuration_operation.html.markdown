---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_provision_role_configuration_operation"
sidebar_current: "docs-tencentcloud-resource-provision_role_configuration_operation"
description: |-
  Provides a resource to create a organization provision_role_configuration_operation
---

# tencentcloud_provision_role_configuration_operation

Provides a resource to create a organization provision_role_configuration_operation

## Example Usage

```hcl
resource "tencentcloud_provision_role_configuration_operation" "provision_role_configuration_operation" {
  zone_id               = "xxxxxx"
  role_configuration_id = "xxxxxx"
  target_type           = "MemberUin"
  target_uin            = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `role_configuration_id` - (Required, String, ForceNew) Permission configuration ID.
* `target_type` - (Required, String, ForceNew) Type of the synchronized target account of the Tencent Cloud Organization. ManagerUin: admin account; MemberUin: member account.
* `target_uin` - (Required, Int, ForceNew) UIN of the target account of the Tencent Cloud Organization.
* `zone_id` - (Required, String, ForceNew) Space ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



