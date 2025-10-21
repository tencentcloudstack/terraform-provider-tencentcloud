---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_role_configurations"
sidebar_current: "docs-tencentcloud-datasource-identity_center_role_configurations"
description: |-
  Use this data source to query detailed information of identity center role configurations
---

# tencentcloud_identity_center_role_configurations

Use this data source to query detailed information of identity center role configurations

## Example Usage

```hcl
data "tencentcloud_identity_center_role_configurations" "identity_center_role_configurations" {
  zone_id = "z-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Space ID.
* `filter_targets` - (Optional, Set: [`Int`]) Check whether the member account has been configured with permissions. If configured, return IsSelected: true; otherwise, return false.
* `filter` - (Optional, String) Filter criteria, which are case insensitive. Currently, only RoleConfigurationName is supported and only eq (Equals) and sw (Start With) are supported. Example: Filter = "RoleConfigurationName, only sw test" means querying all permission configurations starting with test. Filter = "RoleConfigurationName, only eq TestRoleConfiguration" means querying the permission configuration named TestRoleConfiguration.
* `principal_id` - (Optional, String) UserId of the authorized user or GroupId of the authorized user group, which must be set together with the input parameter FilterTargets.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `role_configurations` - Permission configuration list.


