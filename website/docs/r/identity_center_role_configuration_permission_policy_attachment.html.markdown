---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_role_configuration_permission_policy_attachment"
sidebar_current: "docs-tencentcloud-resource-identity_center_role_configuration_permission_policy_attachment"
description: |-
  Provides a resource to create a organization identity_center_role_configuration_permission_policy_attachment
---

# tencentcloud_identity_center_role_configuration_permission_policy_attachment

Provides a resource to create a organization identity_center_role_configuration_permission_policy_attachment

## Example Usage

```hcl
resource "tencentcloud_identity_center_role_configuration_permission_policy_attachment" "identity_center_role_configuration_permission_policy_attachment" {
  zone_id               = "z-xxxxxx"
  role_configuration_id = "rc-xxxxxx"
  role_policy_id        = xxxxxx
}
```

## Argument Reference

The following arguments are supported:

* `role_configuration_id` - (Required, String, ForceNew) Permission configuration ID.
* `role_policy_id` - (Required, Int, ForceNew) Role policy id.
* `zone_id` - (Required, String, ForceNew) Space ID.
* `role_policy_name` - (Optional, String) Role policy name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `add_time` - Role policy add time.
* `role_policy_document` - Role policy document.
* `role_policy_type` - Role policy type.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AddPermissionPolicyToRoleConfiguration](https://cloud.tencent.com/document/api/850/109938) - This API is used to add policies to permission configurations.
* Read: [ListPermissionPoliciesInRoleConfiguration](https://cloud.tencent.com/document/api/850/109934) - This API is used to obtain the policy list in permission configurations.
* Delete: [RemovePermissionPolicyFromRoleConfiguration](https://cloud.tencent.com/document/api/850/109931) - This API is used to remove policies from permission configurations.

## Import

organization identity_center_role_configuration_permission_policy_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_role_configuration_permission_policy_attachment.identity_center_role_configuration_permission_policy_attachment ${zoneId}#${roleConfigurationId}#${rolePolicyIdString}
```

