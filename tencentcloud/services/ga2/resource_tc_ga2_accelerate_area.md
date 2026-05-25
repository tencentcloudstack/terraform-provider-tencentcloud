Provides a resource to create a GA2 accelerate area

Example Usage

```hcl
resource "tencentcloud_ga2_accelerate_area" "example" {
  global_accelerator_id = "ga-xxxxxxxx"

  accelerator_areas {
    accelerate_region = "ap-guangzhou"
    bandwidth         = 10
    isp_type          = "BGP"
    ip_version        = "IPv4"
  }

  accelerator_areas {
    accelerate_region = "ap-shanghai"
    bandwidth         = 20
    isp_type          = "BGP"
    ip_version        = "IPv4"
  }
}
```

Import

GA2 accelerate area can be imported using the globalAcceleratorId, e.g.

```
terraform import tencentcloud_ga2_accelerate_area.example ga-xxxxxxxx
```
