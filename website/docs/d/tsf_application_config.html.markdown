---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_application_config"
sidebar_current: "docs-tencentcloud-datasource-tsf_application_config"
description: |-
  Use this data source to query detailed information of tsf application_config
---

# tencentcloud_tsf_application_config

Use this data source to query detailed information of tsf application_config

## Example Usage

```hcl
data "tencentcloud_tsf_application_config" "application_config" {
  application_id = "app-123456"
  config_id      = "config-123456"
  config_id_list =
  config_name    = "test-config"
  config_version = "1.0"
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Optional, String) Application ID, query all when not provided.
* `config_id_list` - (Optional, Set: [`String`]) Configuration ID list, query all with lower priority when not provided.
* `config_id` - (Optional, String) Configuration ID, query all with higher priority when not provided.
* `config_name` - (Optional, String) Configuration name, precise query, query all when not provided.
* `config_version` - (Optional, String) Configuration version, precise query, query all when not provided.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Paginated configuration item list. Note: This field may return null, indicating that no valid values can be obtained.
  * `content` - Configuration item list.
    * `application_id` - application Id. Note: This field may return null, indicating that no valid values can be obtained.
    * `application_name` - application Id. Note: This field may return null, indicating that no valid values can be obtained.
    * `config_id` - Configuration ID. Note: This field may return null, indicating that no valid values can be obtained.
    * `config_name` - Configuration name Note: This field may return null, indicating that no valid values can be obtained.
    * `config_type` - Configuration type. Note: This field may return null, indicating that no valid values can be obtained.
    * `config_value` - Configuration value. Note: This field may return null, indicating that no valid values can be obtained.
    * `config_version_count` - config version count.  Note: This field may return null, indicating that no valid values can be obtained.
    * `config_version_desc` - Configuration version description. Note: This field may return null, indicating that no valid values can be obtained.
    * `config_version` - Configuration version. Note: This field may return null, indicating that no valid values can be obtained.
    * `creation_time` - CreationTime. Note: This field may return null, indicating that no valid values can be obtained.
    * `delete_flag` - delete flag, true: allow delete; false: delete prohibit.
    * `last_update_time` - last update time.  Note: This field may return null, indicating that no valid values can be obtained.
  * `total_count` - TsfPageConfig.


