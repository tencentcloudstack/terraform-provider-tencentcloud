---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_blueprints"
sidebar_current: "docs-tencentcloud-datasource-lighthouse_blueprints"
description: |-
  Provides a list of Lighthouse blueprints (images).
---

# tencentcloud_lighthouse_blueprints

Provides a list of Lighthouse blueprints (images).

Use this data source to query available blueprints for Lighthouse instances.

## Example Usage

### Query all blueprints:

```hcl
data "tencentcloud_lighthouse_blueprints" "all" {
}

output "blueprints" {
  value = data.tencentcloud_lighthouse_blueprints.all.blueprint_set
}
```

### Filter by platform type:

```hcl
data "tencentcloud_lighthouse_blueprints" "linux" {
  filters {
    name   = "platform-type"
    values = ["LINUX_UNIX"]
  }
}
```

### Filter by blueprint type:

```hcl
data "tencentcloud_lighthouse_blueprints" "app_os" {
  filters {
    name   = "blueprint-type"
    values = ["APP_OS"]
  }
}
```

### Query specific blueprints by ID:

```hcl
data "tencentcloud_lighthouse_blueprints" "specific" {
  blueprint_ids = ["lhbp-xxx", "lhbp-yyy"]
}
```

## Argument Reference

The following arguments are supported:

* `blueprint_ids` - (Optional, Set: [`String`]) Blueprint ID list.
* `filters` - (Optional, List) Filter list.
- `blueprint-id`: Filter by blueprint ID.
- `blueprint-type`: Filter by blueprint type. Values: `APP_OS`, `PURE_OS`, `DOCKER`, `PRIVATE`, `SHARED`.
- `platform-type`: Filter by platform type. Values: `LINUX_UNIX`, `WINDOWS`.
- `blueprint-name`: Filter by blueprint name.
- `blueprint-state`: Filter by blueprint state.
- `scene-id`: Filter by scene ID.
NOTE: The upper limit of Filters per request is 10. The upper limit of Filter.Values is 100. Parameter does not support specifying both BlueprintIds and Filters.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Filter value of field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `blueprint_set` - List of blueprint details.
  * `blueprint_id` - Blueprint ID, which is the unique identifier of Blueprint.
  * `blueprint_name` - Blueprint name.
  * `blueprint_state` - Blueprint state.
  * `blueprint_type` - Blueprint type, such as APP_OS, PURE_OS, DOCKER, PRIVATE, and SHARED.
  * `community_url` - URL of official website of the open-source project.
  * `created_time` - Creation time according to ISO 8601 standard. UTC time is used. Format is YYYY-MM-DDThh:mm:ssZ.
  * `description` - Blueprint description.
  * `display_title` - Blueprint display title.
  * `display_version` - Blueprint display version.
  * `docker_version` - Docker version. Note: This field may return null, indicating that no valid values can be obtained.
  * `guide_url` - Guide documentation URL.
  * `image_id` - ID of the Lighthouse blueprint shared from a CVM image. Note: this field may return null, indicating that no valid values can be obtained.
  * `image_url` - Blueprint image URL.
  * `os_name` - Operating system name.
  * `platform_type` - Platform type, such as LINUX_UNIX and WINDOWS.
  * `platform` - Operating system platform.
  * `required_memory_size` - Memory size required by blueprint in GB.
  * `required_system_disk_size` - System disk size required by blueprint in GB.
  * `scene_id_set` - Array of IDs of scenes associated with a blueprint. Note: This field may return null, indicating that no valid values can be obtained.
  * `support_automation_tools` - Whether the blueprint supports automation tools.


