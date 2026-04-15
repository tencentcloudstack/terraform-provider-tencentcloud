Use this data source to export TEO zone configuration

Example Usage

```hcl
# Export all zone configuration
data "tencentcloud_teo_export_zone_config" "example" {
  zone_id = "zone-2qtuhspy7cr6"
}

# Export specific configuration types
data "tencentcloud_teo_export_zone_config" "specific_types" {
  zone_id = "zone-2qtuhspy7cr6"
  types   = ["L7AccelerationConfig"]
}
```

Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Zone ID.
* `types` - (Optional, List of String) List of configuration types to export. If not specified, all configuration types will be exported. Valid values: L7AccelerationConfig (L7 acceleration configuration).

Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `content` - Exported configuration content in JSON format.
