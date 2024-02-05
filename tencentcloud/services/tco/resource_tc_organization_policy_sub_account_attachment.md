Provides a resource to create a organization policy_sub_account_attachment

Example Usage

```hcl
resource "tencentcloud_organization_policy_sub_account_attachment" "policy_sub_account_attachment" {
  member_uin               = 100028582828
  org_sub_account_uin      = 100028223737
  policy_id                = 144256499
}
```
