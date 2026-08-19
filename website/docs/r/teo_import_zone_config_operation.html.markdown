---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_import_zone_config_operation"
sidebar_current: "docs-tencentcloud-resource-teo_import_zone_config_operation"
description: |-
  Provides a resource to import TEO zone configuration.
---

# tencentcloud_teo_import_zone_config_operation

Provides a resource to import TEO zone configuration.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `content` - (Required, String, ForceNew) The configuration content to be imported, which should be in the JSON format and be encoded in the UTF-8 mode. The configuration content can be obtained through the site configuration export API (ExportZoneConfig). You can individually import "Site Acceleration - Global Acceleration Configuration" or "Site Acceleration - Rule Engine" by passing in the corresponding fields. Refer to the example below for details.
* `zone_id` - (Required, String, ForceNew) Site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `finish_time` - The end time of the import task.
* `import_time` - The start time of the import task.
* `message` - The status message of the import task. When the configuration import fails, you can view the failure reason through this field.
* `status` - The import task status. Valid values: success, failure, doing.
* `task_id` - The task ID of the import configuration operation.


