---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_config_extra"
sidebar_current: "docs-tencentcloud-resource-cls_config_extra"
description: |-
  Provides a resource to create a cls config extra
---

# tencentcloud_cls_config_extra

Provides a resource to create a cls config extra

## Example Usage

```hcl
resource "tencentcloud_cls_config_extra" "extra" {
  name        = "helloworld"
  topic_id    = tencentcloud_cls_topic.topic.id
  type        = "container_file"
  log_type    = "json_log"
  config_flag = "label_k8s"
  logset_id   = tencentcloud_cls_logset.logset.id
  logset_name = tencentcloud_cls_logset.logset.logset_name
  topic_name  = tencentcloud_cls_topic.topic.topic_name
  #  host_file {
  #    log_path = "/var/log/tmep"
  #    file_pattern = "*.log"
  #    custom_labels = ["key1=value1"]
  #  }
  container_file {
    container    = "nginx"
    file_pattern = "log"
    log_path     = "/nginx"
    namespace    = "default"
    workload {
      container = "nginx"
      kind      = "deployment"
      name      = "nginx"
      namespace = "default"
    }
  }
  group_id = "27752a9b-9918-440a-8ee7-9c84a14a47ed"
}
```

## Argument Reference

The following arguments are supported:

* `config_flag` - (Required, String) Collection configuration flag.
* `log_type` - (Required, String) Type of the log to be collected. Valid values: json_log: log in JSON format; delimiter_log: log in delimited format; minimalist_log: minimalist log; multiline_log: log in multi-line format; fullregex_log: log in full regex format. Default value: minimalist_log.
* `logset_id` - (Required, String) Logset Id.
* `logset_name` - (Required, String) Logset Name.
* `name` - (Required, String) Collection configuration name.
* `topic_id` - (Required, String) Log topic ID (TopicId) of collection configuration.
* `topic_name` - (Required, String) Topic Name.
* `type` - (Required, String) Type. Valid values: container_stdout; container_file; host_file.
* `container_file` - (Optional, List) Container file path info.
* `container_stdout` - (Optional, List) Container stdout info.
* `exclude_paths` - (Optional, List) Collection path blocklist.
* `extract_rule` - (Optional, List) Extraction rule. If ExtractRule is set, LogType must be set.
* `group_id` - (Optional, String) Binding group id.
* `group_ids` - (Optional, Set: [`String`], ForceNew) Binding group ids.
* `host_file` - (Optional, List) Node file config info.
* `user_define_rule` - (Optional, String) Custom collection rule, which is a serialized JSON string.

The `container_file` object supports the following:

* `container` - (Required, String) Container name.
* `file_pattern` - (Required, String) log name.
* `log_path` - (Required, String) Log Path.
* `namespace` - (Required, String) Namespace. There can be multiple namespaces, separated by separators, such as A, B.
* `exclude_labels` - (Optional, Set) Pod label to be excluded.
* `exclude_namespace` - (Optional, String) Namespaces to be excluded, separated by separators, such as A, B.
* `include_labels` - (Optional, Set) Pod label info.
* `workload` - (Optional, List) Workload info.

The `container_stdout` object supports the following:

* `all_containers` - (Required, Bool) Is all containers.
* `exclude_labels` - (Optional, Set) Pod label to be excluded.
* `exclude_namespace` - (Optional, String) Namespaces to be excluded, separated by separators, such as A, B.
* `include_labels` - (Optional, Set) Pod label info.
* `namespace` - (Optional, String) Namespace. There can be multiple namespaces, separated by separators, such as A, B.
* `workloads` - (Optional, List) Workload info.

The `exclude_paths` object supports the following:

* `type` - (Optional, String) Type. Valid values: File, Path.
* `value` - (Optional, String) Specific content corresponding to Type.

The `extract_rule` object supports the following:

* `backtracking` - (Optional, Int) Size of the data to be rewound in incremental collection mode. Default value: -1 (full collection).
* `begin_regex` - (Optional, String) First-Line matching rule, which is valid only if log_type is multiline_log or fullregex_log.
* `delimiter` - (Optional, String) Delimiter for delimited log, which is valid only if log_type is delimiter_log.
* `filter_key_regex` - (Optional, List) Log keys to be filtered and the corresponding regex.
* `keys` - (Optional, Set) Key name of each extracted field. An empty key indicates to discard the field. This parameter is valid only if log_type is delimiter_log. json_log logs use the key of JSON itself.
* `log_regex` - (Optional, String) Full log matching rule, which is valid only if log_type is fullregex_log.
* `time_format` - (Optional, String) Time field format. For more information, please see the output parameters of the time format description of the strftime function in C language.
* `time_key` - (Optional, String) Time field key name. time_key and time_format must appear in pair.
* `un_match_log_key` - (Optional, String) Unmatched log key.
* `un_match_up_load_switch` - (Optional, Bool) Whether to upload the logs that failed to be parsed. Valid values: true: yes; false: no.

The `filter_key_regex` object supports the following:

* `key` - (Optional, String) Log key to be filtered.
* `regex` - (Optional, String) Filter rule regex corresponding to key.

The `host_file` object supports the following:

* `file_pattern` - (Required, String) Log file name.
* `log_path` - (Required, String) Log file dir.
* `custom_labels` - (Optional, Set) Metadata info.

The `workload` object supports the following:

* `kind` - (Required, String) workload type.
* `name` - (Required, String) workload name.
* `container` - (Optional, String) container name.
* `namespace` - (Optional, String) namespace.

The `workloads` object supports the following:

* `kind` - (Required, String) workload type.
* `name` - (Required, String) workload name.
* `container` - (Optional, String) container name.
* `namespace` - (Optional, String) namespace.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



