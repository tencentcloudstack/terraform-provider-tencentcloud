Provides a resource to export TEO zone configuration

~> **NOTE:** This resource is used to export zone configuration for backup or migration purposes. It does not actually create or modify any resources in the cloud.

Example Usage

```hcl
resource "tencentcloud_teo_export_zone_config" "example" {
  zone_id     = "zone-39quuimqg8r6"
  export_type = "all"
}
```

Export basic configuration only

```hcl
resource "tencentcloud_teo_export_zone_config" "basic" {
  zone_id     = "zone-39quuimqg8r6"
  export_type = "basic"
}
```

Export HTTPS configuration

```hcl
resource "tencentcloud_teo_export_zone_config" "https" {
  zone_id     = "zone-39quuimqg8r6"
  export_type = "https"
}
```

Export cache configuration

```hcl
resource "tencentcloud_teo_export_zone_config" "cache" {
  zone_id     = "zone-39quuimqg8r6"
  export_type = "cache"
}
```

Export WAF configuration

```hcl
resource "tencentcloud_teo_export_zone_config" "waf" {
  zone_id     = "zone-39quuimqg8r6"
  export_type = "waf"
}
```

Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) ID of the site.
* `export_type` - (Required, ForceNew) Export type. Valid values: `all`: all configurations; `basic`: basic configurations; `cache`: cache configurations; `https`: HTTPS configurations; `origin`: origin configurations; `waf`: WAF configurations; `rate_limit`: rate limit configurations; `rule_engine`: rule engine configurations.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `config_content` - Exported configuration content in JSON format.

Import

TEO export zone config can be imported using id, e.g.

```
terraform import tencentcloud_teo_export_zone_config.example zone-39quuimqg8r6#all
```
