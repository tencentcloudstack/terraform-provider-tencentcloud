# tencentcloud_teo_export_zone_config

## Example Usage

### Basic Usage

```hcl
resource "tencentcloud_teo_export_zone_config" "basic" {
  zone_id = "zone-xxxxxxx"
}
```

### Export Specific Configuration Types

```hcl
resource "tencentcloud_teo_export_zone_config" "specific" {
  zone_id = "zone-xxxxxxx"
  types   = ["L7AccelerationConfig"]
}
```

### Export with Custom Timeouts

```hcl
resource "tencentcloud_teo_export_zone_config" "with_timeout" {
  zone_id = "zone-xxxxxxx"

  timeouts {
    create = "10m"
    read   = "10m"
  }
}
```

## Import

teo export zone config can be imported using the zone_id:

```shell
terraform import tencentcloud_teo_export_zone_config.example zone-xxxxxxx
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Site ID.
* `types` - (Optional, ForceNew, List) List of configuration types to export. If left blank, all configuration types are exported. Supported values include:
  * `L7AccelerationConfig`: Export Layer 7 acceleration configuration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `content` - Exported configuration content in JSON format, encoded in UTF-8.

## Notes

* This resource exports the configuration of a TEO zone. The export operation is synchronous and returns the configuration content as JSON.
* Changes to `zone_id` or `types` require resource recreation due to ForceNew attribute.
* The exported configuration can be used for backup, audit, or other purposes.
* Delete operation on this resource only removes it from Terraform state, as the export itself is not a persistent resource in the TEO API.
