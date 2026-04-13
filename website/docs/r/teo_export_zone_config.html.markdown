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
resource "tencentcloud_teo_export_zone_config" "export" {
  zone_id = "zone-xxxxxxx"
}
```

### Export Specific Configuration Types

```hcl
resource "tencentcloud_teo_export_zone_config" "export_specific" {
  zone_id = "zone-xxxxxxx"
  types   = ["L7AccelerationConfig"]
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Site ID.
* `types` - (Optional, ForceNew, List) List of configuration types to export. If left blank, all configuration types are exported. Supported values include: `L7AccelerationConfig`: Export Layer 7 acceleration configuration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `content` - Exported configuration content in JSON format, encoded in UTF-8.

## Import

teo export zone config can be imported using the zone_id, e.g.

```
terraform import tencentcloud_teo_export_zone_config.export zone-xxxxxxx
```
