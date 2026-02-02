---
subcategory: "Video on Demand(VOD)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_sub_applications"
sidebar_current: "docs-tencentcloud-datasource-vod_sub_applications"
description: |-
  Use this data source to query detailed information of VOD sub-applications.
---

# tencentcloud_vod_sub_applications

Use this data source to query detailed information of VOD sub-applications.

## Example Usage

### Query all sub-applications

```hcl
data "tencentcloud_vod_sub_applications" "all" {}

output "app_list" {
  value = data.tencentcloud_vod_sub_applications.all.sub_application_info_set
}
```

### Filter by application name

```hcl
data "tencentcloud_vod_sub_applications" "by_name" {
  name = "MyVideoApp"
}

output "app_id" {
  value = data.tencentcloud_vod_sub_applications.by_name.sub_application_info_set[0].sub_app_id
}
```

### Filter by tags

```hcl
data "tencentcloud_vod_sub_applications" "by_tags" {
  tags = {
    Environment = "Production"
    Team        = "VideoTeam"
  }
}

output "production_apps" {
  value = data.tencentcloud_vod_sub_applications.by_tags.sub_application_info_set[*].sub_app_id_name
}
```

### Reference sub-application in other resources

```hcl
data "tencentcloud_vod_sub_applications" "existing" {
  name = "ProductionApp"
}

resource "tencentcloud_vod_super_player_config" "config" {
  sub_app_id = data.tencentcloud_vod_sub_applications.existing.sub_application_info_set[0].sub_app_id
  name       = "player-config"
  drm_switch = false
  # ... other configuration
}
```

### Export results to file

```hcl
data "tencentcloud_vod_sub_applications" "export" {
  result_output_file = "/tmp/vod_apps.json"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Application name for exact match filtering.
* `result_output_file` - (Optional, String) Used to save results in JSON format.
* `tags` - (Optional, Map) Tag key-value pairs for filtering applications. Applications matching all specified tags will be returned.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `sub_application_info_set` - List of sub-application information.
  * `create_time` - Creation time in ISO 8601 format.
  * `description` - Sub-application description.
  * `mode` - Application mode. Valid values: fileid, fileid+path.
  * `name` - Legacy name field (for backward compatibility).
  * `status` - Application status. Valid values: On, Off, Destroying, Destroyed.
  * `storage_regions` - List of enabled storage regions.
  * `sub_app_id_name` - Sub-application name.
  * `sub_app_id` - Sub-application ID.
  * `tags` - Resource tags bound to the sub-application.


