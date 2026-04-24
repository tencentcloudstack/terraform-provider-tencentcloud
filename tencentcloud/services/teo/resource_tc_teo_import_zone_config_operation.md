Provides a resource to import TEO zone configuration.

Example Usage

```hcl
data "tencentcloud_teo_export_zone_config" "example" {
  zone_id = "zone-id1"
  types   = ["L7AccelerationConfig"]
}

resource "tencentcloud_teo_import_zone_config_operation" "example" {
  zone_id = "zone-id2"
  content = data.tencentcloud_teo_export_zone_config.example.content
}
```
