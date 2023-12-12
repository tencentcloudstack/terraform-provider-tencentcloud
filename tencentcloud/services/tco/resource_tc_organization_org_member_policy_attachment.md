Provides a resource to create a organization org_member_policy_attachment

Example Usage

```hcl
resource "tencentcloud_organization_org_member_policy_attachment" "org_member_policy_attachment" {
  member_uins = [100033905366,100033905356]
  policy_name = "example-iac"
  identity_id = 1
}
```

Import

organization org_member_policy_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_member_policy_attachment.org_member_policy_attachment org_member_policy_attachment_id
```