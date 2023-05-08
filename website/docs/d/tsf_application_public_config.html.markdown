---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_application_public_config"
sidebar_current: "docs-tencentcloud-datasource-tsf_application_public_config"
description: |-
  Use this data source to query detailed information of tsf application_public_config
---

# tencentcloud_tsf_application_public_config

Use this data source to query detailed information of tsf application_public_config

## Example Usage

```hcl
data "tencentcloud_tsf_application_public_config" "application_public_config" {
  config_id = "dcfg-p-evjrbgly"
  # config_id_list = [""]
  config_name    = "dsadsa"
  config_version = "123"
}
```

## Argument Reference

The following arguments are supported:

* `config_id_list` - (Optional, Set: [`String`]) Config ID list. Query all items if not passed, low priority.
* `config_id` - (Optional, String) Config ID. Query all items if not passed, high priority.
* `config_name` - (Optional, String) Config name. Exact query. Query all items if not passed.
* `config_version` - (Optional, String) Config version. Exact query. Query all items if not passed.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Paginated global configuration  list. Note: This field may return null, indicating that no valid value can be obtained.
  * `content` - Config list.
    * `application_id` - application Id. Note: This field may return null, indicating that no valid values can be obtained.
    * `application_name` - application Id. Note: This field may return null, indicating that no valid values can be obtained.
    * `config_id` - Config ID. Note: This field may return null, indicating that no valid value can be obtained.
    * `config_name` - Config name. Note: This field may return null, indicating that no valid value can be obtained.
    * `config_type` - Config type. Note: This field may return null, indicating that no valid value can be obtained.
    * `config_value` - Config value. Note: This field may return null, indicating that no valid value can be obtained.
    * `config_version_count` - config version count.  Note: This field may return null, indicating that no valid values can be obtained.
    * `config_version_desc` - Config version description. Note: This field may return null, indicating that no valid value can be obtained.
    * `config_version` - Config version. Note: This field may return null, indicating that no valid value can be obtained.
    * `creation_time` - CreationTime. Note: This field may return null, indicating that no valid values can be obtained.
    * `delete_flag` - delete flag, true: allow delete; false: delete prohibit.
    * `last_update_time` - last update time.  Note: This field may return null, indicating that no valid values can be obtained.
  * `total_count` - TsfPageConfig.


