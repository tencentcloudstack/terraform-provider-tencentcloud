Provides a resource to create a teo zone

Example Usage

Basic Usage

```hcl
resource "tencentcloud_teo_zone" "zone" {
  zone_name       = "tf-teo.com"
  type            = "partial"
  area            = "overseas"
  alias_zone_name = "teo-test"
  paused          = false
  plan_id         = "edgeone-2kfv1h391n6w"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Enable Version Control Mode

```hcl
resource "tencentcloud_teo_zone" "zone_with_version_control" {
  zone_name       = "tf-teo-version.com"
  type            = "partial"
  area            = "overseas"
  alias_zone_name = "teo-version-test"
  paused          = false
  plan_id         = "edgeone-2kfv1h391n6w"
  
  work_mode_infos {
    config_group_type = "l7_acceleration"
    work_mode         = "version_control"
  }
  
  tags = {
    "createdBy" = "terraform"
  }
}
```


Import

teo zone can be imported using the id, e.g.
```
terraform import tencentcloud_teo_zone.zone zone_id
```