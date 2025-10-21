---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_machine_group_configs"
sidebar_current: "docs-tencentcloud-datasource-cls_machine_group_configs"
description: |-
  Use this data source to query detailed information of cls machine_group_configs
---

# tencentcloud_cls_machine_group_configs

Use this data source to query detailed information of cls machine_group_configs

## Example Usage

```hcl
resource "tencentcloud_cls_machine_group" "group" {
  group_name        = "tf-describe-mg-config-test"
  service_logging   = true
  auto_update       = true
  update_end_time   = "19:05:00"
  update_start_time = "17:05:00"

  machine_group_type {
    type = "ip"
    values = [
      "192.168.1.1",
      "192.168.1.2",
    ]
  }
}

data "tencentcloud_cls_machine_group_configs" "machine_group_configs" {
  group_id = tencentcloud_cls_machine_group.group.id
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) group id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `configs` - scrape config list.
  * `config_id` - scrape config id.
  * `create_time` - create time.
  * `exclude_paths` - Collection path blocklist.
    * `type` - Type. Valid values: File, Path.
    * `value` - Specific content corresponding to Type.
  * `extract_rule` - Extraction rule. If ExtractRule is set, LogType must be set.
    * `address` - syslog system log collection specifies the address and port that the collector listens to.
    * `backtracking` - Size of the data to be rewound in incremental collection mode. Default value: -1 (full collection).
    * `begin_regex` - First-Line matching rule, which is valid only if log_type is multiline_log or fullregex_log.
    * `delimiter` - Delimiter for delimited log, which is valid only if log_type is delimiter_log.
    * `filter_key_regex` - Log keys to be filtered and the corresponding regex.
      * `key` - Log key to be filtered.
      * `regex` - Filter rule regex corresponding to key.
    * `is_gbk` - GBK encoding. Default 0.
    * `json_standard` - standard json. Default 0.
    * `keys` - Key name of each extracted field. An empty key indicates to discard the field. This parameter is valid only if log_type is delimiter_log. json_log logs use the key of JSON itself.
    * `log_regex` - Full log matching rule, which is valid only if log_type is fullregex_log.
    * `meta_tags` - metadata tags.
      * `key` - tag key.
      * `value` - tag value.
    * `metadata_type` - metadata type.
    * `parse_protocol` - parse protocol.
    * `path_regex` - metadata path regex.
    * `protocol` - syslog protocol, tcp or udp.
    * `time_format` - Time field format. For more information, please see the output parameters of the time format description of the strftime function in C language.
    * `time_key` - Time field key name. time_key and time_format must appear in pair.
    * `un_match_log_key` - Unmatched log key.
    * `un_match_up_load_switch` - Whether to upload the logs that failed to be parsed. Valid values: true: yes; false: no.
  * `log_format` - style of log format.
  * `log_type` - log type.
  * `name` - scrape config name.
  * `output` - topicid.
  * `path` - scrape log path.
  * `update_time` - update time.
  * `user_define_rule` - user define rule.


