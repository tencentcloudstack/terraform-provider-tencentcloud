# Terraform Configuration for TEO Export Zone Config

This example demonstrates how to use the `tencentcloud_teo_export_zone_config` resource to export zone configuration.

## Usage

### Basic Usage

Export all configuration types for a zone:

```hcl
terraform {
  required_providers {
    tencentcloud = {
      source = "tencentcloudstack/tencentcloud"
    }
  }
}

provider "tencentcloud" {
  region = "ap-guangzhou"
}

resource "tencentcloud_teo_zone" "example" {
  type      = "partial"
  area      = "overseas"
  plan_id   = "edgeone-2kfv1h391n6w"
  zone_name = "example.com"
}

resource "tencentcloud_teo_export_zone_config" "example" {
  zone_id = tencentcloud_teo_zone.example.id

  depends_on = [tencentcloud_teo_zone.example]
}
```

### Export Specific Configuration Types

Export only specific configuration types by specifying the `types` parameter:

```hcl
resource "tencentcloud_teo_export_zone_config" "example" {
  zone_id = tencentcloud_teo_zone.example.id
  types   = ["L7AccelerationConfig"]

  depends_on = [tencentcloud_teo_zone.example]
}
```

## Arguments Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Zone ID. Example: zone-2zpqp7qztest
* `types` - (Optional) List of configuration types to export. If not specified, all configuration types will be exported.
  * `L7AccelerationConfig` - Export L7 acceleration configuration, corresponding to "Site Acceleration - Global Acceleration Configuration" and "Site Acceleration - Rule Engine" in the console.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `content` - The specific content of the exported configuration. Returned in JSON format and encoded in UTF-8.

## Import

Teo export zone config can be imported using the zone_id, e.g.

```sh
terraform import tencentcloud_teo_export_zone_config.example zone-2zpqp7qztest
```
