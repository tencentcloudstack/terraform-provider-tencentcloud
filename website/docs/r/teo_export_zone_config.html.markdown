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

```hcl
resource "tencentcloud_teo_zone" "zone" {
  zone_name = "example.com"
  type      = "full"
}

resource "tencentcloud_teo_export_zone_config" "export_zone_config" {
  zone_id = tencentcloud_teo_zone.zone.id
}
```

```hcl
resource "tencentcloud_teo_export_zone_config" "export_zone_config" {
  zone_id = "zone-297z8rf93cfw"

  types = [
    "L7AccelerationConfig",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) ID of the site.
* `types` - (Optional, List) List of configuration types to export. If not specified, all configuration types will be exported. Currently supported types:
  - `L7AccelerationConfig`: Seven-layer acceleration configuration (corresponds to "Site Acceleration - Global Acceleration Configuration" and "Site Acceleration - Rule Engine" in the console).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource, which is the zone_id.
* `content` - The exported configuration content in JSON format, encoded in UTF-8. This content can be used for importing zone configuration.

## Import

teo_export_zone_config can be imported using the zone_id, e.g.

```
terraform import tencentcloud_teo_export_zone_config.export_zone_config zone-297z8rf93cfw
```
