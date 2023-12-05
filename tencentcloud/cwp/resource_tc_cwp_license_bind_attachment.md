Provides a resource to create a cwp license_bind_attachment

Example Usage

```hcl
data "tencentcloud_cwp_machines_simple" "example" {
  machine_type   = "CVM"
  machine_region = "ap-guangzhou"
  project_ids    = [1210293]

  filters {
    name        = "Version"
    values      = ["BASIC_VERSION"]
    exact_match = true
  }
}

resource "tencentcloud_cwp_license_order" "example" {
  alias        = "tf_example"
  license_type = 0
  license_num  = 1
  region_id    = 1
  project_id   = 0
  tags         = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_cwp_license_bind_attachment" "example" {
  resource_id  = tencentcloud_cwp_license_order.example.resource_id
  license_id   = tencentcloud_cwp_license_order.example.license_id
  license_type = 0
  quuid        = data.tencentcloud_cwp_machines_simple.example.machines[0].quuid
}
```

Import

cwp license_bind_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cwp_license_bind_attachment.example cwplic-ab3edffa#44#2c7e5cce-1cec-4456-8d18-018f160dd987#0
```