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

Usage with allow_duplicates

```hcl
resource "tencentcloud_teo_zone" "zone_with_duplicates" {
  zone_name        = "tf-teo-duplicates.com"
  type             = "partial"
  area             = "overseas"
  alias_zone_name  = "teo-duplicates-test"
  paused           = false
  plan_id          = "edgeone-2kfv1h391n6w"
  allow_duplicates = true  # Allow duplicate rule configurations
  tags = {
    "createdBy" = "terraform"
  }
}
```

**Important:** The `allow_duplicates` field can only be set during resource creation and cannot be updated afterwards. If you need to change this value, you must recreate the zone.

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
    work_mode         = "immediate_effect"
  }
  work_mode_infos {
    config_group_type = "edge_functions"
    work_mode         = "immediate_effect"
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