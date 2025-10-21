---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_role_configuration_permission_custom_policy_attachment"
sidebar_current: "docs-tencentcloud-resource-identity_center_role_configuration_permission_custom_policy_attachment"
description: |-
  Provides a resource to create a organization identity_center_role_configuration_permission_custom_policy_attachment
---

# tencentcloud_identity_center_role_configuration_permission_custom_policy_attachment

Provides a resource to create a organization identity_center_role_configuration_permission_custom_policy_attachment

## Example Usage

```hcl
resource "tencentcloud_identity_center_role_configuration_permission_custom_policy_attachment" "identity_center_role_configuration_permission_custom_policy_attachment" {
  zone_id               = "z-xxxxxx"
  role_configuration_id = "rc-xxxxxx"
  role_policy_name      = "CustomPolicy"
  role_policy_document  = <<-EOF
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
```

## Argument Reference

The following arguments are supported:

* `role_configuration_id` - (Required, String, ForceNew) Permission configuration ID.
* `role_policy_document` - (Required, String, ForceNew) Role policy document.
* `role_policy_name` - (Required, String, ForceNew) Role policy name.
* `zone_id` - (Required, String, ForceNew) Space ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `add_time` - Role policy add time.
* `role_policy_type` - Role policy type.


## Import

organization identity_center_role_configuration_permission_custom_policy_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_role_configuration_permission_custom_policy_attachment.identity_center_role_configuration_permission_custom_policy_attachment ${zoneId}#${roleConfigurationId}#${rolePolicyName}
```

