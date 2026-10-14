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

* `content` - (Required, String, ForceNew) Configuration content to be imported. Use JSON format and encode by UTF-8. You can obtain the configuration content through the site configuration export interface (ExportZoneConfig). You can individually import "Site Acceleration - Global Acceleration Configuration" or "Site Acceleration - Rule Engine". Just pass in the corresponding fields. For specific details, see the example below. Note: AccelerationDomain (acceleration domain name configuration) and Origin (origin configuration) exported by ExportZoneConfig are temporary not supported for import through this interface. If the Content contains the above configuration content, it will cause import failure.
* `zone_id` - (Required, String, ForceNew) Site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `finish_time` - The end time of the import task.
* `import_time` - The start time of the import task.
* `message` - The status message of the import task. When the configuration import fails, you can view the failure reason through this field.
* `status` - The import task status. Valid values: success, failure, doing.
* `task_id` - The task ID of the import configuration operation.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeZoneConfigImportResult](https://cloud.tencent.com/document/api/1552/113231) - This API is used to query the results of site configuration import via API (ImportZoneConfig). This feature only supports the sites in the plans of the Standard Edition and the Enterprise Edition.
* Create: [ImportZoneConfig](https://cloud.tencent.com/document/api/1552/113229) - This API is used to quickly import site configuration files. After the import is initiated, the API will return the corresponding task ID.

