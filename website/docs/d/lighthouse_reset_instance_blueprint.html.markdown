---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_reset_instance_blueprint"
sidebar_current: "docs-tencentcloud-datasource-lighthouse_reset_instance_blueprint"
description: |-
  Use this data source to query detailed information of lighthouse reset_instance_blueprint
---

# tencentcloud_lighthouse_reset_instance_blueprint

Use this data source to query detailed information of lighthouse reset_instance_blueprint

## Example Usage

```hcl
data "tencentcloud_lighthouse_reset_instance_blueprint" "reset_instance_blueprint" {
  instance_id = "lhins-123456"
  offset      = 0
  limit       = 20
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `filters` - (Optional, List) Filter listblueprint-idFilter by image ID.Type: StringRequired: noblueprint-typeFilter by image type.Valid values: APP_OS: application image; PURE_OS: system image; PRIVATE: custom imageType: StringRequired: noplatform-typeFilter by image platform type.Valid values: LINUX_UNIX: Linux or Unix; WINDOWS: WindowsType: StringRequired: noblueprint-nameFilter by image name.Type: StringRequired: noblueprint-stateFilter by image status.Type: StringRequired: noEach request can contain up to 10 Filters and 5 Filter.Values. BlueprintIds and Filters cannot be specified at the same time.
* `limit` - (Optional, Int) Number of returned results. Default value is 20. Maximum value is 100.
* `offset` - (Optional, Int) Offset. Default value is 0.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Filter value of field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `reset_instance_blueprint_set` - List of scene info.
  * `blueprint_info` - Mirror details.
    * `blueprint_id` - Image ID, which is the unique identity of Blueprint.
    * `blueprint_name` - Mirror name.
    * `blueprint_state` - Mirror status.
    * `blueprint_type` - Image type, such as APP_OS, PURE_OS, PRIVATE.
    * `community_url` - The official website Url.
    * `created_time` - Creation time. Expressed according to the ISO8601 standard, and using UTC time. The format is YYYY-MM-DDThh:mm:ssZ.
    * `description` - Mirror description information.
    * `display_title` - The mirror image shows the title to the public.
    * `display_version` - The image shows the version to the public.
    * `docker_version` - Docker version number.
    * `guide_url` - Guide article Url.
    * `image_id` - CVM image ID after sharing the CVM image to the lightweight application server.
    * `image_url` - Mirror image URL.
    * `os_name` - Operating system name.
    * `platform_type` - Operating system platform type, such as LINUX_UNIX, WINDOWS.
    * `platform` - Operating system platform.
    * `required_memory_size` - Memory required for mirroring (in GB).
    * `required_system_disk_size` - The size of the system disk required for image (in GB).
    * `scene_id_set` - The mirror association uses the scene Id list.
    * `support_automation_tools` - Whether the image supports automation helper.
  * `is_resettable` - Whether the instance image can be reset to the target image.
  * `non_resettable_message` - The information cannot be reset. when the mirror can be reset ''.


