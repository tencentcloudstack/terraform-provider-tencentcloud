Provides a resource to create a CWP license order

Example Usage

```hcl
resource "tencentcloud_cwp_license_order" "example" {
  alias        = "tf_example"
  license_type = 0
  license_num  = 1
  region_id    = 1
  project_id   = 0
  tags = {
    createdBy = "terraform"
  }
}
```

Import

CWP license order can be imported using the resourceId#regionId, e.g.

```
terraform import tencentcloud_cwp_license_order.example cwplic-130715d2#1
```