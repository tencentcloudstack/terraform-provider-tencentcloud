Provides a resource to create a Organization share unit member

~> **NOTE:** This resource is deprecated, please use `tencentcloud_organization_org_share_unit_member_v2` instead.

Example Usage

```hcl
resource "tencentcloud_organization_org_share_unit" "example" {
  name        = "tf-example"
  area        = "ap-guangzhou"
  description = "description."
}

resource "tencentcloud_organization_org_share_unit_member" "example" {
  unit_id = tencentcloud_organization_org_share_unit.example.unit_id
  area    = tencentcloud_organization_org_share_unit.example.area
  members {
    share_member_uin = 100035309479
  }
}
```
