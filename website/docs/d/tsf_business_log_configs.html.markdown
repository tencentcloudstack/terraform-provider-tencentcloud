---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_business_log_configs"
sidebar_current: "docs-tencentcloud-datasource-tsf_business_log_configs"
description: |-
  Use this data source to query detailed information of tsf business_log_configs
---

# tencentcloud_tsf_business_log_configs

Use this data source to query detailed information of tsf business_log_configs

## Example Usage

```hcl
data "tencentcloud_tsf_business_log_configs" "business_log_configs" {
  search_word                = "terraform"
  disable_program_auth_check = true
  config_id_list             = ["apm-busi-log-cfg-qv3x3rdv"]
}
```

## Argument Reference

The following arguments are supported:

* `config_id_list` - (Optional, Set: [`String`]) Config Id list.
* `disable_program_auth_check` - (Optional, Bool) Disable Program auth check or not.
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) wild search word.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - List of business log configurations.Note: This field may return null, indicating that no valid values can be obtained.
  * `content` - Log configuration item list. Note: This field may return null, indicating that no valid values can be obtained.
    * `config_associated_groups` - the associate group of Config.Note: This field may return null, indicating that no valid values can be obtained.
      * `application_id` - Application Id of Group. Note: This field may return null, indicating that no valid values can be obtained.
      * `application_name` - Application Name. Note: This field may return null, indicating that no valid values can be obtained.
      * `application_type` - Application Type. Note: This field may return null, indicating that no valid values can be obtained.
      * `associated_time` - Time when the deployment group is associated with the log configuration.Note: This field may return null, indicating that no valid values can be obtained.
      * `cluster_id` - Cluster ID to which the deployment group belongs.Note: This field may return null, indicating that no valid values can be obtained.
      * `cluster_name` - Cluster Name to which the deployment group belongs.Note: This field may return null, indicating that no valid values can be obtained.
      * `cluster_type` - Cluster type to which the deployment group belongs.Note: This field may return null, indicating that no valid values can be obtained.
      * `group_id` - Group Id. Note: This field may return null, indicating that no valid values can be obtained.
      * `group_name` - Group Name. Note: This field may return null, indicating that no valid values can be obtained.
      * `namespace_id` - Namespace ID to which the deployment group belongs.Note: This field may return null, indicating that no valid values can be obtained.
      * `namespace_name` - Namespace Name to which the deployment group belongs.Note: This field may return null, indicating that no valid values can be obtained.
    * `config_create_time` - Create time of configuration item.Note: This field may return null, indicating that no valid values can be obtained.
    * `config_desc` - Description of configuration item.Note: This field may return null, indicating that no valid values can be obtained.
    * `config_id` - ConfigId.
    * `config_name` - ConfigName.
    * `config_path` - Log path of configuration item.Note: This field may return null, indicating that no valid values can be obtained.
    * `config_pipeline` - Pipeline of configuration item.Note: This field may return null, indicating that no valid values can be obtained.
    * `config_schema` - ParserSchema of configuration item.Note: This field may return null, indicating that no valid values can be obtained.
      * `schema_content` - content of schema.
      * `schema_create_time` - Create time of configuration item.Note: This field may return null, indicating that no valid values can be obtained.
      * `schema_date_format` - Schema format.Note: This field may return null, indicating that no valid values can be obtained.
      * `schema_multiline_pattern` - Schema pattern of configuration item.Note: This field may return null, indicating that no valid values can be obtained.
      * `schema_pattern_layout` - User-defined parsing rules.Note: This field may return null, indicating that no valid values can be obtained.
      * `schema_type` - Schema type.
    * `config_tags` - configuration Tag.Note: This field may return null, indicating that no valid values can be obtained.
    * `config_update_time` - Update time of configuration item.Note: This field may return null, indicating that no valid values can be obtained.
  * `total_count` - Total Count.Note: This field may return null, indicating that no valid values can be obtained.


