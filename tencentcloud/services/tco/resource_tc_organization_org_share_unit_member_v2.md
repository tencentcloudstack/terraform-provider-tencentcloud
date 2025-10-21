Provides a resource to create a Organization share unit member

~> **NOTE:** This resource must exclusive in one share unit, do not declare additional members resources of this share unit elsewhere.

Example Usage

```hcl
resource "tencentcloud_organization_org_share_unit" "example" {
  name        = "tf-example"
  area        = "ap-guangzhou"
  description = "description."
}

resource "tencentcloud_organization_org_share_unit_member_v2" "example" {
  unit_id = tencentcloud_organization_org_share_unit.example.unit_id
  area    = tencentcloud_organization_org_share_unit.example.area
  members {
    share_member_uin = 100042257812
  }

  members {
    share_member_uin = 100043990767
  }

  members {
    share_member_uin = 100042234123
  }
}
```

Import

Organization share unit member can be imported using the unitId#area, e.g.

```
terraform import tencentcloud_organization_org_share_unit_member_v2.example shareUnit-switt8i4s4#ap-guangzhou
```
