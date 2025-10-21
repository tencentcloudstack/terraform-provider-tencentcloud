---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_user_sync_provisioning"
sidebar_current: "docs-tencentcloud-resource-identity_center_user_sync_provisioning"
description: |-
  Provides a resource to create a organization identity_center_user_sync_provisioning
---

# tencentcloud_identity_center_user_sync_provisioning

Provides a resource to create a organization identity_center_user_sync_provisioning

## Example Usage

```hcl
resource "tencentcloud_identity_center_user_sync_provisioning" "identity_center_user_sync_provisioning" {
  zone_id              = "z-xxxxxx"
  description          = "tf-test"
  deletion_strategy    = "Keep"
  duplication_strategy = "TakeOver"
  principal_id         = "u-xxxxxx"
  principal_type       = "User"
  target_uin           = "xxxxxx"
  target_type          = "MemberUin"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Space ID.
* `deletion_strategy` - (Optional, String) Deletion policy. It indicates the handling policy for CAM users already synchronized when the CAM user synchronization is deleted. Valid values: Delete: Delete the CAM users already synchronized from CIC to CAM when the CAM user synchronization is deleted; Keep: Keep the CAM users already synchronized from CIC to CAM when the CAM user synchronization is deleted.
* `description` - (Optional, String) Description.
* `duplication_strategy` - (Optional, String) Conflict policy. It indicates the handling policy for existence of a user with the same username when CIC users are synchronized to CAM. Valid values: KeepBoth: Keep both, that is, add the _cic suffix to the CIC user's username and then try to create a CAM user with the username when CIC users are synchronized to CAM and a user with the same username already exists in CAM; TakeOver: Replace, that is, directly replace the existing CAM user with the synchronized CIC user when CIC users are synchronized to CAM and a user with the same username already exists in CAM.
* `principal_id` - (Optional, String) Identity ID for the CAM user synchronization. Valid values:
When the PrincipalType value is Group, it is the CIC user group ID (g-********).
When the PrincipalType value is User, it is the CIC user ID (u-********).
* `principal_type` - (Optional, String) Identity type for the CAM user synchronization. Valid values:

User: indicates that the identity for the CAM user synchronization is a CIC user.
Group: indicates that the identity for the CAM user synchronization is a CIC user group.
* `target_type` - (Optional, String) Type of the synchronized target account of the Tencent Cloud Organization. ManagerUin: admin account; MemberUin: member account.
* `target_uin` - (Optional, Int) UIN of the synchronized target account of the Tencent Cloud Organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `principal_name` - The identity name of the CAM user synchronization. Value: When PrincipalType is Group, the value is the CIC user group name; When PrincipalType takes the value to User, the value is the CIC user name.
* `status` - Status of CAM user synchronization. Value:
	* Enabled: CAM user synchronization is enabled;
	* Disabled: CAM user synchronization is not enabled.
* `target_name` - Group account The name of the target account..
* `update_time` - Update time.
* `user_provisioning_id` - User provisioning id.


## Import

organization identity_center_user_sync_provisioning can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_user_sync_provisioning.identity_center_user_sync_provisioning ${zoneId}#${userProvisioningId}
```

