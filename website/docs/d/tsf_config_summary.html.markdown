---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_config_summary"
sidebar_current: "docs-tencentcloud-datasource-tsf_config_summary"
description: |-
  Use this data source to query detailed information of tsf config_summary
---

# tencentcloud_tsf_config_summary

Use this data source to query detailed information of tsf config_summary

## Example Usage

```hcl
data "tencentcloud_tsf_config_summary" "config_summary" {
  application_id             = "application-a24x29xv"
  search_word                = "terraform"
  order_by                   = "last_update_time"
  order_type                 = 0
  disable_program_auth_check = true
  config_id_list             = ["dcfg-y54wzk3a"]
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Optional, String) Application ID. If not passed, the query will be for all.
* `config_id_list` - (Optional, Set: [`String`]) Config Id List.
* `config_tag_list` - (Optional, Set: [`String`]) config tag list.
* `disable_program_auth_check` - (Optional, Bool) Whether to disable dataset authentication.
* `order_by` - (Optional, String) Order term. support Sort by time: creation_time; or Sort by name: config_name.
* `order_type` - (Optional, Int) Pass 0 for ascending order and 1 for descending order.
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) Query keyword, fuzzy query: application name, configuration item name. If not passed, the query will be for all.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - config Page Item.
  * `content` - config list.
    * `application_id` - Application ID.Note: This field may return null, indicating that no valid value was found.
    * `application_name` - Application Name. Note: This field may return null, indicating that no valid value was found.
    * `config_id` - Configuration item ID.Note: This field may return null, indicating that no valid value was found.
    * `config_name` - Configuration name.Note: This field may return null, indicating that no valid value was found.
    * `config_type` - Config type. Note: This field may return null, indicating that no valid value was found.
    * `config_value` - Configuration value.Note: This field may return null, indicating that no valid value was found.
    * `config_version_count` - Configure version count.Note: This field may return null, indicating that no valid value was found.
    * `config_version_desc` - Configuration version description.Note: This field may return null, indicating that no valid value was found.
    * `config_version` - Configuration version. Note: This field may return null, indicating that no valid value was found.
    * `creation_time` - Create time.Note: This field may return null, indicating that no valid value was found.
    * `delete_flag` - Deletion flag, true: deletable; false: not deletable.Note: This field may return null, indicating that no valid value was found.
    * `last_update_time` - Last update time.Note: This field may return null, indicating that no valid value was found.
  * `total_count` - total count.


