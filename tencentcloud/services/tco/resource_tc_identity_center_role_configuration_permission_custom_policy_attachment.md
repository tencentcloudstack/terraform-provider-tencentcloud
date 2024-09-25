Provides a resource to create a organization identity_center_role_configuration_permission_custom_policy_attachment

Example Usage

```hcl
resource "tencentcloud_identity_center_role_configuration_permission_custom_policy_attachment" "identity_center_role_configuration_permission_custom_policy_attachment" {
    zone_id = "z-xxxxxx"
    role_configuration_id = "rc-xxxxxx"
    role_policy_name = "CustomPolicy"
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
```

Import

organization identity_center_role_configuration_permission_custom_policy_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_role_configuration_permission_custom_policy_attachment.identity_center_role_configuration_permission_custom_policy_attachment ${zoneId}#${roleConfigurationId}#${rolePolicyName}
```
