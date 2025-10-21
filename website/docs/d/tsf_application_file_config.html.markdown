---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_application_file_config"
sidebar_current: "docs-tencentcloud-datasource-tsf_application_file_config"
description: |-
  Use this data source to query detailed information of tsf application_file_config
---

# tencentcloud_tsf_application_file_config

Use this data source to query detailed information of tsf application_file_config

## Example Usage

```hcl
data "tencentcloud_tsf_application_file_config" "application_file_config" {
  config_id = "dcfg-f-4y4ekzqv"
  # config_id_list = [""]
  config_name    = "file-log1"
  application_id = "application-2vzk6n3v"
  config_version = "1.2"
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Optional, String) Application ID.
* `config_id_list` - (Optional, Set: [`String`]) List of configuration item ID.
* `config_id` - (Optional, String) Configuration ID.
* `config_name` - (Optional, String) Configuration item name.
* `config_version` - (Optional, String) Configuration item version.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - File configuration item list. Note: This field may return null, indicating that no valid values can be obtained.
  * `content` - File configuration array. Note: This field may return null, indicating that no valid values can be obtained.
    * `application_id` - application Id. Note: This field may return null, indicating that no valid values can be obtained.
    * `application_name` - application name. Note: This field may return null, indicating that no valid values can be obtained.
    * `config_file_code` - Configuration file code. Note: This field may return null, indicating that no valid values can be obtained.
    * `config_file_name` - Configuration item file name. Note: This field may return null, indicating that no valid values can be obtained.
    * `config_file_path` - file config path. Note: This field may return null, indicating that no valid values can be obtained.
    * `config_file_value_length` - config item content length.  Note: This field may return null, indicating that no valid values can be obtained.
    * `config_file_value` - Configuration file content. Note: This field may return null, indicating that no valid values can be obtained.
    * `config_id` - Config ID. Note: This field may return null, indicating that no valid values can be obtained.
    * `config_name` - Configuration item name. Note: This field may return null, indicating that no valid values can be obtained.
    * `config_post_cmd` - last update time.  Note: This field may return null, indicating that no valid values can be obtained.
    * `config_version_count` - config version count.  Note: This field may return null, indicating that no valid values can be obtained.
    * `config_version_desc` - Configuration item version description. Note: This field may return null, indicating that no valid values can be obtained.
    * `config_version` - Configuration version. Note: This field may return null, indicating that no valid values can be obtained.
    * `creation_time` - CreationTime. Note: This field may return null, indicating that no valid values can be obtained.
    * `delete_flag` - delete flag, true: allow delete; false: delete prohibit.
    * `last_update_time` - last update time.  Note: This field may return null, indicating that no valid values can be obtained.
  * `total_count` - total count.


