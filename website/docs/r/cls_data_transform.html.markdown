---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_data_transform"
sidebar_current: "docs-tencentcloud-resource-cls_data_transform"
description: |-
  Provides a resource to create a CLS data transform
---

# tencentcloud_cls_data_transform

Provides a resource to create a CLS data transform

## Example Usage

```hcl
resource "tencentcloud_cls_logset" "logset_src" {
  logset_name = "tf-example-src"
  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_cls_topic" "topic_src" {
  topic_name           = "tf-example_src"
  logset_id            = tencentcloud_cls_logset.logset_src.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_cls_logset" "logset_dst" {
  logset_name = "tf-example-dst"
  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_cls_topic" "topic_dst" {
  topic_name           = "tf-example-dst"
  logset_id            = tencentcloud_cls_logset.logset_dst.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_cls_data_transform" "example" {
  func_type    = 1
  src_topic_id = tencentcloud_cls_topic.topic_src.id
  name         = "tf-example"
  etl_content  = "ext_sep(\"content\", \"f1, f2, f3\", sep=\",\", quote=\"\", restrict=False, mode=\"overwrite\")fields_drop(\"content\")"
  task_type    = 3
  enable_flag  = 1
  dst_resources {
    topic_id = tencentcloud_cls_topic.topic_dst.id
    alias    = "iac-test-dst"
  }
}
```

## Argument Reference

The following arguments are supported:

* `etl_content` - (Required, String) Data transform content. If `func_type` is `2`, must use `log_auto_output`.
* `func_type` - (Required, Int) Task type. `1`: Specify the theme; `2`: Dynamic creation.
* `name` - (Required, String) Task name.
* `src_topic_id` - (Required, String) Source topic ID.
* `task_type` - (Required, Int) Task type. `1`: Use random data from the source log theme for processing preview; `2`: Use user-defined test data for processing preview; `3`: Create real machining tasks.
* `backup_give_up_data` - (Optional, Bool) When `func_type` is `2`, whether to discard data when the number of dynamically created logsets and topics exceeds the product specification limit. Default is `false`. `false`: Create backup logset and topic and write logs to the backup topic; `true`: Discard log data.
* `data_transform_sql_data_sources` - (Optional, List) Associated data source information.
* `data_transform_type` - (Optional, Int) Data transform type. `0`: Standard data transform task; `1`: Pre-processing data transform task (process collected logs before writing to the log topic).
* `dst_resources` - (Optional, List) Data transform des resources. If `func_type` is `1`, this parameter is required. If `func_type` is `2`, this parameter does not need to be filled in.
* `enable_flag` - (Optional, Int) Task enable flag. `1`: enable, `2`: disable, Default is `1`.
* `env_infos` - (Optional, List) Set environment variables.
* `failure_log_key` - (Optional, String) Field name for failure logs.
* `has_services_log` - (Optional, Int) Whether to enable service log delivery. `1`: disable; `2`: enable.
* `keep_failure_log` - (Optional, Int) Keep failure log status. `1`: do not keep (default); `2`: keep.
* `process_from_timestamp` - (Optional, Int) Specify the start time of processing data, in seconds-level timestamp. Any time range within the log topic lifecycle. If it exceeds the lifecycle, only the part with data within the lifecycle is processed.
* `process_to_timestamp` - (Optional, Int) Specify the end time of processing data, in seconds-level timestamp. Cannot specify a future time. If not filled, it means continuous execution.

The `data_transform_sql_data_sources` object supports the following:

* `alias_name` - (Required, String) Alias. Used in data transform statements.
* `data_source` - (Required, Int) Data source type. `1`: MySQL; `2`: Self-built MySQL; `3`: PostgreSQL.
* `instance_id` - (Required, String) Instance ID. When DataSource is `1`, it represents the cloud database MySQL instance ID, such as: cdb-zxcvbnm.
* `password` - (Required, String) MySQL access password.
* `region` - (Required, String) InstanceId region. For example: ap-guangzhou.
* `user` - (Required, String) MySQL access username.

The `dst_resources` object supports the following:

* `alias` - (Required, String) Alias.
* `topic_id` - (Required, String) Dst topic ID.

The `env_infos` object supports the following:

* `key` - (Required, String) Environment variable name.
* `value` - (Required, String) Environment variable value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CLS data transform can be imported using the id, e.g.

```
terraform import tencentcloud_cls_data_transform.example 7b4bcb05-9154-4cdc-a479-f6b5743846e5
```

