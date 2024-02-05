Provides a resource to create a organization org_share_unit_member

Example Usage

```hcl
resource "tencentcloud_organization_org_share_unit_member" "org_share_unit_member" {
  unit_id = &lt;nil&gt;
  area = &lt;nil&gt;
  members {
		share_member_uin = &lt;nil&gt;

  }
}
```

Import

organization org_share_unit_member can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_share_unit_member.org_share_unit_member org_share_unit_member_id
```
