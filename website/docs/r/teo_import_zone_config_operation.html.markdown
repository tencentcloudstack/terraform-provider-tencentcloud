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

* `content` - (Required, String, ForceNew) The configuration content to import. It must be in JSON format and encoded in UTF-8. You can obtain the configuration content via the tencentcloud_teo_export_zone_config data source.
* `zone_id` - (Required, String, ForceNew) Site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `finish_time` - The end time of the import task.
* `import_time` - The start time of the import task.
* `message` - The status message of the import task. When the configuration import fails, you can view the failure reason through this field.
* `status` - The import task status. Valid values: success, failure, doing.
* `task_id` - The task ID of the import configuration operation.


