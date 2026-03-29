## tencentcloud_teo_export_zone_config

Provides a data source to export the complete configuration of a TEO (TencentCloud EdgeOne) zone.

This data source allows you to query and retrieve detailed configuration information for a TEO zone, including basic information, acceleration settings, security rules, and origin server settings.

## Example Usage

### Query by Zone ID

```hcl
data "tencentcloud_teo_export_zone_config" "config" {
  zone_id = "zone-xxxxx"
}

output "zone_config" {
  value = data.tencentcloud_teo_export_zone_config.config
}
```

### Query by Zone Name

```hcl
data "tencentcloud_teo_export_zone_config" "config" {
  zone_name = "example.com"
}

output "zone_config" {
  value = data.tencentcloud_teo_export_zone_config.config
}
```

### Query with Both Zone ID and Zone Name (Zone ID takes precedence)

```hcl
data "tencentcloud_teo_export_zone_config" "config" {
  zone_id   = "zone-xxxxx"
  zone_name = "example.com"
}

output "zone_config" {
  value = data.tencentcloud_teo_export_zone_config.config
}
```

### Export Zone Configuration to File

```hcl
data "tencentcloud_teo_export_zone_config" "config" {
  zone_id              = "zone-xxxxx"
  result_output_file   = "zone_config.json"
}

output "basic_info" {
  value = {
    zone_id    = data.tencentcloud_teo_export_zone_config.config.zone_id_output
    zone_name  = data.tencentcloud_teo_export_zone_config.config.zone_name_output
    area       = data.tencentcloud_teo_export_zone_config.config.area
    type       = data.tencentcloud_teo_export_zone_config.config.type
    status     = data.tencentcloud_teo_export_zone_config.config.status
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Optional) Zone ID. Specify either zone_id or zone_name.
* `zone_name` - (Optional) Zone name. Specify either zone_id or zone_name.
* `result_output_file` - (Optional) Used to save results.

**Note:** At least one of `zone_id` or `zone_name` must be specified. If both are provided, `zone_id` takes precedence.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `zone_id_output` - Zone ID.
* `zone_name_output` - Zone name.
* `area` - Acceleration area. Values: `global`: Global. `mainland`: Chinese mainland. `overseas`: Outside the Chinese mainland.
* `type` - Site access method. Valid values: `full`: NS access; `partial`: CNAME access; `noDomainAccess`: access with no domain name.
* `status` - The site status. Values: `active`: The name server is switched to EdgeOne. `pending`: The name server is not switched. `moved`: The name server is changed to other service providers. `deactivated`: The site is blocked. `initializing`: The site is not bound with any plan.
* `created_on` - The creation time of the site.
* `modified_on` - The modification date of the site.
* `paused` - Whether the site is disabled.
* `active_status` - Status of the proxy. Values: `active`: Enabled; `inactive`: Not activated; `paused`: Disabled.
* `cache_settings` - Cache configuration settings.
* `security_settings` - Security configuration settings.
* `origin_settings` - Origin server configuration settings.

## Import

Export zone config does not support import.
