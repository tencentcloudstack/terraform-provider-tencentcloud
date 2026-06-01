---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_role_assignment"
sidebar_current: "docs-tencentcloud-resource-identity_center_role_assignment"
description: |-
  Provides a resource to create a Organization identity center role assignment
---

# tencentcloud_identity_center_role_assignment

Provides a resource to create a Organization identity center role assignment

## Example Usage

```hcl
resource "tencentcloud_identity_center_role_assignment" "example" {
  zone_id               = "z-1os7c9znogct"
  principal_id          = "u-lyfm8b7qoi5l"
  principal_type        = "User"
  target_uin            = "100043911945"
  target_type           = "MemberUin"
  role_configuration_id = "rc-ihogrs0e6ceg"
}
```

## Argument Reference

The following arguments are supported:

* `principal_id` - (Required, String, ForceNew) Identity ID for the CAM user synchronization. Valid values:
When the PrincipalType value is Group, it is the CIC user group ID (g-********).
When the PrincipalType value is User, it is the CIC user ID (u-********).
* `principal_type` - (Required, String, ForceNew) Identity type for the CAM user synchronization. Valid values:

User: indicates that the identity for the CAM user synchronization is a CIC user.
Group: indicates that the identity for the CAM user synchronization is a CIC user group.
* `role_configuration_id` - (Required, String, ForceNew) Permission configuration ID.
* `target_type` - (Required, String, ForceNew) Type of the synchronized target account of the Tencent Cloud Organization. ManagerUin: admin account; MemberUin: member account.
* `target_uin` - (Required, Int, ForceNew) UIN of the synchronized target account of the Tencent Cloud Organization.
* `zone_id` - (Required, String, ForceNew) Space ID.
* `deprovision_strategy` - (Optional, String, ForceNew) When you remove the last authorization configured with a certain privilege on a group account target account, whether to cancel the privilege configuration deployment at the same time. Value: DeprovisionForLastRoleAssignmentOnAccount: Remove privileges to configure deployment. None (default): Configure deployment without delegating privileges.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `principal_name` - Principal name.
* `role_configuration_name` - Role configuration name.
* `target_name` - Target name.
* `update_time` - Update time.


## Import

Organization identity center role assignment can be imported using the {zoneId}#{roleConfigurationId}#{targetType}#{targetUinString}#{principalType}, e.g.

```
terraform import tencentcloud_identity_center_role_assignment.example z-1os7c9znogct#rc-ihogrs0e6ceg#MemberUin#100043911945#User#u-lyfm8b7qoi5l
```

