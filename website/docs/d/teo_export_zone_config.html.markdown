---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_export_zone_config"
sidebar_current: "docs-tencentcloud-datasource-teo_export_zone_config"
description: |-
  Use this data source to export TEO (EdgeOne) site configuration
---

# tencentcloud_teo_export_zone_config

Use this data source to export TEO (EdgeOne) site configuration

## Example Usage

### Export all types of zone configuration

```hcl
data "tencentcloud_teo_export_zone_config" "example" {
  zone_id = "zone-3fkff38fyw8s"
}
```

### Export specific types of zone configuration

```hcl
data "tencentcloud_teo_export_zone_config" "example" {
  zone_id = "zone-3fkff38fyw8s"
  types   = ["L7AccelerationConfig"]
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Specifies the site ID.
* `result_output_file` - (Optional, String) Used to save results.
* `types` - (Optional, List: [`String`]) Types of configuration to export. If not specified, all types of configuration will be exported. Valid values: `L7AccelerationConfig`, `WebSecurity`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `content` - Exported zone configuration content in JSON format.


