Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) accelerate area.

Example Usage

Create accelerate area with BGP

```hcl
resource "tencentcloud_ga2_global_accelerator" "example" {
  name                 = "tf-example"
  instance_charge_type = "POSTPAID"
  description          = "tf example global accelerator"

  tags = {
    createdBy = "Terraform"
  }
}

resource "tencentcloud_ga2_accelerate_area" "example" {
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
  accelerate_region     = "ap-guangzhou"
  bandwidth             = 10
  isp_type              = "BGP"
  ip_version            = "IPv4"
}
```

Create accelerate area with QUALITY_BGP

```hcl
resource "tencentcloud_ga2_accelerate_area" "example" {
  global_accelerator_id = "ga-jg9gepn0"
  accelerate_region     = "ap-hongkong"
  bandwidth             = 50
  isp_type              = "QUALITY_BGP"
  ip_version            = "IPv4"
  timeouts {
    create = "15m"
    delete = "15m"
  }
}
```

Import

GA2 accelerate area can be imported using the composite id `<global_accelerator_id>#<accelerator_area_id>`, e.g.

```
terraform import tencentcloud_ga2_accelerate_area.example ga-jg9gepn0#area-jrsub43y
```
