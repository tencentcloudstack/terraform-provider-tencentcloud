Provides a resource to create a organization member email

Example Usage

```hcl
resource "tencentcloud_organization_org_member_email" "example" {
  member_uin   = 100033118139
  email        = "example@tencent.com"
  country_code = "86"
  phone        = "18611111111"
}
```

Import

organization member email can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_member_email.example 100033118139#132
```