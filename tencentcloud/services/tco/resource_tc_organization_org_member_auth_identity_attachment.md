Provides a resource to create a organization org_member_auth_identity

Example Usage

```hcl
resource "tencentcloud_organization_org_member_auth_identity_attachment" "org_member_auth_identity" {
  member_uin = 100033704327
  identity_ids = [1657]
}
```

Import

organization org_member_auth_identity can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_member_auth_identity.org_member_auth_identity org_member_auth_identity_id
```