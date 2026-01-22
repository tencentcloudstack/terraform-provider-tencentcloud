Provides a resource to create a organization tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment

Example Usage

```hcl
resource "tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment" "identity_center_role_configuration_permission_custom_policies_attachment" {
    zone_id = "z-xxxxxx"
    role_configuration_id = "rc-xxxxxx"
    policies {
        role_policy_name = "CustomPolicy1"
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

Import

organization tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment.identity_center_role_configuration_permission_custom_policies_attachment ${zoneId}#${roleConfigurationId}#${rolePolicyName1},...${rolePolicyNameN}
```
