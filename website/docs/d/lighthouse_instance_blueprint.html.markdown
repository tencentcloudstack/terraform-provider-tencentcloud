---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_instance_blueprint"
sidebar_current: "docs-tencentcloud-datasource-lighthouse_instance_blueprint"
description: |-
  Use this data source to query detailed information of lighthouse instance_blueprint
---

# tencentcloud_lighthouse_instance_blueprint

Use this data source to query detailed information of lighthouse instance_blueprint

## Example Usage

```hcl
data "tencentcloud_lighthouse_instance_blueprint" "instance_blueprint" {
  instance_ids = ["lhins-xxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_ids` - (Required, Set: [`String`]) Instance ID list, which currently can contain only one instance.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `blueprint_instance_set` - Blueprint instance list information.
  * `blueprint` - Blueprint instance information.
    * `blueprint_id` - Blueprint ID, which is the unique identifier of Blueprint.
    * `blueprint_name` - Blueprint name.
    * `blueprint_state` - Blueprint status.
    * `blueprint_type` - Blueprint type, such as APP_OS, PURE_OS, and PRIVATE.
    * `community_url` - URL of official website of the open-source project.
    * `created_time` - Creation time according to ISO 8601 standard. UTC time is used. Format is YYYY-MM-DDThh:mm:ssZ.
    * `description` - Image description information.
    * `display_title` - Blueprint title to be displayed.
    * `display_version` - Blueprint version to be displayed.
    * `docker_version` - Docker version.Note: This field may return null, indicating that no valid values can be obtained.
    * `guide_url` - Guide documentation URL.
    * `image_id` - ID of the Lighthouse blueprint shared from a CVM imageNote: this field may return null, indicating that no valid values can be obtained.
    * `image_url` - Blueprint picture URL.
    * `os_name` - OS name.
    * `platform_type` - OS type, such as LINUX_UNIX and WINDOWS.
    * `platform` - OS type.
    * `required_memory_size` - Memory size required by blueprint in GB.
    * `required_system_disk_size` - System disk size required by blueprint in GB.
    * `scene_id_set` - Array of IDs of scenes associated with a blueprintNote: This field may return null, indicating that no valid values can be obtained.
    * `support_automation_tools` - Whether the blueprint supports automation tools.
  * `instance_id` - Instance ID.
  * `software_set` - Software list.
    * `detail_set` - List of software details.
      * `key` - Unique detail key.
      * `title` - Detail title.
      * `value` - Detail value.
    * `image_url` - Software picture URL.
    * `install_dir` - Software installation directory.
    * `name` - Software name.
    * `version` - Software version.


