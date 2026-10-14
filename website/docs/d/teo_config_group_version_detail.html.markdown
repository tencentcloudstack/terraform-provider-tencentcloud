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
data "tencentcloud_teo_config_group_version_detail" "teo_config_group_version_detail" {
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
* `content` - Version file content. It is returned in JSON format.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeConfigGroupVersionDetail](https://cloud.tencent.com/document/api/1552/101865) - This API is used to obtain detailed information about a version in version management mode. The response includes the version ID, description, status, creation time, configuration group information, and the content of the version ...

