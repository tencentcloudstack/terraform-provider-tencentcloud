---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_export_zone_config"
sidebar_current: "docs-tencentcloud-resource-teo_export_zone_config"
description: |-
  Provides a resource to export teo zone configuration
---

# tencentcloud_teo_export_zone_config

Provides a resource to export teo zone configuration

## Example Usage

### Basic Usage

```hcl
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

```hcl
resource "tencentcloud_teo_export_zone_config" "example" {
  zone_id = tencentcloud_teo_zone.example.id
  types   = ["L7AccelerationConfig"]

  depends_on = [tencentcloud_teo_zone.example]
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) Zone ID. Example: zone-2zpqp7qztest
* `types` - (Optional, List) List of configuration types to export. If not specified, all configuration types will be exported.
  - `L7AccelerationConfig`: Export L7 acceleration configuration, corresponding to "Site Acceleration - Global Acceleration Configuration" and "Site Acceleration - Rule Engine" in the console. Note: The supported export types will increase with iteration. When exporting all types, please pay attention to the size of the exported file. It is recommended to specify the configuration types to be exported to control the size of the request response payload.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `content` - The specific content of the exported configuration. Returned in JSON format and encoded in UTF-8.

## Import

Teo export zone config can be imported using the zone_id, e.g.

```
terraform import tencentcloud_teo_export_zone_config.example zone-2zpqp7qztest
```
