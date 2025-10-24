---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment"
sidebar_current: "docs-tencentcloud-resource-identity_center_role_configuration_permission_custom_policies_attachment"
description: |-
  Provides a resource to create a organization tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment
---

# tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment

Provides a resource to create a organization tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment

## Example Usage

```hcl
resource "tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment" "identity_center_role_configuration_permission_custom_policies_attachment" {
  zone_id               = "z-xxxxxx"
  role_configuration_id = "rc-xxxxxx"
  policies {
    role_policy_name     = "CustomPolicy1"
    role_policy_document = <<-EOF
{
    "version": "2.0",
    "statement": [
        {
            "effect": "allow",
            "action": [
                "vpc:AcceptAttachCcnInstances"
            ],
            "resource": [
                "*"
            ]
        }
    ]
}
EOF
  }

}
```

## Argument Reference

The following arguments are supported:

* `policies` - (Required, Set, ForceNew) Policies.
* `role_configuration_id` - (Required, String, ForceNew) Permission configuration ID.
* `zone_id` - (Required, String, ForceNew) Space ID.

The `policies` object supports the following:

* `role_policy_document` - (Required, String, ForceNew) Role policy document.
* `role_policy_name` - (Required, String, ForceNew) Role policy name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

organization tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment.identity_center_role_configuration_permission_custom_policies_attachment ${zoneId}#${roleConfigurationId}#${rolePolicyName1},...${rolePolicyNameN}
```

