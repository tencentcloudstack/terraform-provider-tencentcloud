---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_config_group_version_detail"
sidebar_current: "docs-tencentcloud-datasource-teo_config_group_version_detail"
description: |-
  Use this data source to query detailed information of teo config group version detail
---

# tencentcloud_teo_config_group_version_detail

Use this data source to query detailed information of teo config group version detail

## Example Usage

```hcl
data "tencentcloud_teo_config_group_version_detail" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  version_id = "sv-2quhspyeq8r6"
}
```

## Argument Reference

The following arguments are supported:

* `version_id` - (Required, String) Version ID.
* `zone_id` - (Required, String) Zone ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `config_group_version_info` - Version information.
  * `create_time` - Version creation time. The time format follows the ISO 8601 standard and is represented in Coordinated Universal Time (UTC).
  * `description` - Version description.
  * `group_id` - Configuraration group ID.
  * `group_type` - Configuration group type. Valid values: l7_acceleration: L7 acceleration configuration group. edge_functions: Edge function configuration group.
  * `source_version` - The source version ID that the config group version was derived from.
  * `status` - Version status. Valid values: creating: Being created.inactive: Not effective.active: Effective.
  * `version_id` - Version ID.
  * `version_number` - Version No.
* `content` - Version file content. It is returned in JSON format.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeConfigGroupVersionDetail](https://www.tencentcloud.com/document/api/1145/57975)

