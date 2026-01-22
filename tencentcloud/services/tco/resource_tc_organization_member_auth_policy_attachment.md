Provides a resource to create a Organization member auth policy attachment

Example Usage

```hcl
resource "tencentcloud_organization_member_auth_policy_attachment" "example" {
  policy_id           = 252421751
  org_sub_account_uin = 100037718939
}
```

Import

Organization member auth policy attachment can be imported using the id, e.g.

```
terraform import tencentcloud_organization_member_auth_policy_attachment.example 252421751#100037718939
```
