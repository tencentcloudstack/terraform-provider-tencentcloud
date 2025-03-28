---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_config"
sidebar_current: "docs-tencentcloud-resource-cls_config"
description: |-
  Provides a resource to create a cls config
---

# tencentcloud_cls_config

Provides a resource to create a cls config

## Example Usage

```hcl
resource "tencentcloud_cls_config" "config" {
  name     = "config_hello"
  output   = "4d07fba0-b93e-4e0b-9a7f-d58542560bbb"
  path     = "/var/log/kubernetes"
  log_type = "json_log"
  extract_rule {
    filter_key_regex {
      key   = "key1"
      regex = "value1"
    }
    filter_key_regex {
      key   = "key2"
      regex = "value2"
    }
    un_match_up_load_switch = true
    un_match_log_key        = "config"
    backtracking            = -1
  }
  exclude_paths {
    type  = "Path"
    value = "/data"
  }
  exclude_paths {
    type  = "File"
    value = "/file"
  }
  #  user_define_rule = ""
}
```

## Argument Reference

The following arguments are supported:

* `extract_rule` - (Required, List) Extraction rule. If ExtractRule is set, LogType must be set.
* `name` - (Required, String) Collection configuration name.
* `exclude_paths` - (Optional, List) Collection path blocklist.
* `log_type` - (Optional, String) Type of the log to be collected. Valid values: json_log: log in JSON format; delimiter_log: log in delimited format; minimalist_log: minimalist log; multiline_log: log in multi-line format; fullregex_log: log in full regex format. Default value: minimalist_log.
* `output` - (Optional, String) Log topic ID (TopicId) of collection configuration.
* `path` - (Optional, String) Log collection path containing the filename. Required for document collection.
* `user_define_rule` - (Optional, String) Custom collection rule, which is a serialized JSON string. Required when LogType is user_define_log.

The `exclude_paths` object supports the following:

* `type` - (Optional, String) Type. Valid values: File, Path.
* `value` - (Optional, String) Specific content corresponding to Type.

The `extract_rule` object supports the following:

* `address` - (Optional, String) syslog system log collection specifies the address and port that the collector listens to. This parameter is only valid when LogType is service_syslog. It does not need to be filled in for other types.
* `backtracking` - (Optional, Int) Size of the data to be rewound in incremental collection mode. Default value: -1 (full collection).
* `begin_regex` - (Optional, String) First-Line matching rule, which is valid only if log_type is multiline_log or fullregex_log.
* `delimiter` - (Optional, String) Delimiter for delimited log, which is valid only if log_type is delimiter_log.
* `filter_key_regex` - (Optional, List) Log keys to be filtered and the corresponding regex.
* `is_gbk` - (Optional, Int) GBK encoding. Default 0. Note: - Currently, when the value is 0, it means UTF-8 encoding.
* `json_standard` - (Optional, Int) standard json. Default 0.
* `keys` - (Optional, Set) Key name of each extracted field. An empty key indicates to discard the field. This parameter is valid only if log_type is delimiter_log. json_log logs use the key of JSON itself.
* `log_regex` - (Optional, String) Full log matching rule, which is valid only if log_type is fullregex_log.
* `meta_tags` - (Optional, List) metadata tags. Note: - Required when MetadataType is 2. - COS import does not support this field.
* `metadata_type` - (Optional, Int) metadata type. 0: Do not use metadata information; 1: Use machine group metadata; 2: Use user-defined metadata; 3: Use collection configuration path. Note: COS import does not support this field.
* `parse_protocol` - (Optional, String) parse protocol. This parameter is only valid when LogType is service_syslog. It does not need to be filled in for other types.
* `path_regex` - (Optional, String) metadata path regex.
* `protocol` - (Optional, String) syslog protocol, tcp or udp. The value can be tcp or udp. It is effective only when LogType is service_syslog. Other types do not need to be filled in.
* `time_format` - (Optional, String) Time field format. For more information, please see the output parameters of the time format description of the strftime function in C language.
* `time_key` - (Optional, String) Time field key name. time_key and time_format must appear in pair.
* `un_match_log_key` - (Optional, String) Unmatched log key. Required when UnMatchUpLoadSwitch is true.
* `un_match_up_load_switch` - (Optional, Bool) Whether to upload the logs that failed to be parsed. Valid values: true: yes; false: no.

The `filter_key_regex` object of `extract_rule` supports the following:

* `key` - (Optional, String) Log key to be filtered.
* `regex` - (Optional, String) Filter rule regex corresponding to key.

The `meta_tags` object of `extract_rule` supports the following:

* `key` - (Optional, String) tag key.
* `value` - (Optional, String) tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls config can be imported using the id, e.g.

```
terraform import tencentcloud_cls_config.config config_id
```

