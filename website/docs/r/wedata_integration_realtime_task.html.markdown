---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_integration_realtime_task"
sidebar_current: "docs-tencentcloud-resource-wedata_integration_realtime_task"
description: |-
  Provides a resource to create a wedata integration_realtime_task
---

# tencentcloud_wedata_integration_realtime_task

Provides a resource to create a wedata integration_realtime_task

## Example Usage

```hcl
resource "tencentcloud_wedata_integration_realtime_task" "example" {
  project_id  = "1455251608631480391"
  task_name   = "tf_example"
  description = "description."
  sync_type   = 1
  task_info {
    incharge    = "100028439226"
    executor_id = "20230704142425553913"
    config {
      name  = "concurrency"
      value = "1"
    }
    config {
      name  = "TaskManager"
      value = "1"
    }
    config {
      name  = "JobManager"
      value = "1"
    }
    config {
      name  = "TolerateDirtyData"
      value = "0"
    }
    config {
      name  = "CheckpointingInterval"
      value = "1"
    }
    config {
      name  = "CheckpointingIntervalUnit"
      value = "min"
    }
    config {
      name  = "RestartStrategyFixedDelayAttempts"
      value = "-1"
    }
    config {
      name  = "ResourceAllocationType"
      value = "0"
    }
    config {
      name  = "TaskAlarmRegularList"
      value = ""
    }
    mappings {
      source_id = "2"
      sink_id   = "1"
    }
    nodes {
      id               = "1"
      name             = "gf_poc"
      node_type        = "INPUT"
      data_source_type = "MYSQL"
      datasource_id    = "5737"
      config {
        name  = "StartupMode"
        value = "INIT"
      }
      config {
        name  = "Encode"
        value = "utf-8"
      }
      config {
        name  = "Database"
        value = "UNKNOW"
      }
      config {
        name  = "SourceRule"
        value = "all"
      }
      config {
        name  = "FilterOper"
        value = "update"
      }
      config {
        name  = "ServerTimeZone"
        value = "Asia/Shanghai"
      }
      config {
        name  = "GhostChange"
        value = "false"
      }
      config {
        name  = "TableNames"
        value = "gf_db.*,hx_db.*,information_schema.*,mysql.*,performance_schema.*,run_time.*,sys.*,test01.*"
      }
      config {
        name  = "FirstDataSource"
        value = "5737"
      }
      config {
        name  = "MultipleDataSources"
        value = "5737"
      }
      config {
        name  = "SiblingNodes"
        value = "[]"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `sync_type` - (Required, Int) Synchronization type: 1. Whole database synchronization, 2. Single table synchronization.
* `task_info` - (Required, List) Task Information.
* `task_name` - (Required, String) Task name.
* `description` - (Optional, String) Description information.

The `config` object supports the following:

* `name` - (Optional, String) Configuration name.
* `value` - (Optional, String) Configuration value.

The `execute_context` object supports the following:

* `name` - (Optional, String) Configuration name.
* `value` - (Optional, String) Configuration value.

The `ext_config` object supports the following:

* `name` - (Optional, String) Configuration name.
* `value` - (Optional, String) Configuration value.

The `mappings` object supports the following:

* `ext_config` - (Optional, List) Node extension configuration information.
* `schema_mappings` - (Optional, List) Schema mapping information.
* `sink_id` - (Optional, String) Sink node ID.
* `source_id` - (Optional, String) Source node ID.
* `source_schema` - (Optional, List) Source node schema information.

The `node_mapping` object supports the following:

* `ext_config` - (Optional, List) Node extension configuration information.
* `schema_mappings` - (Optional, List) Schema mapping information.
* `sink_id` - (Optional, String) Sink node ID.
* `source_id` - (Optional, String) Source node ID.
* `source_schema` - (Optional, List) Source node schema information.

The `nodes` object supports the following:

* `app_id` - (Optional, String) User App Id.
* `config` - (Optional, List) Node configuration information.
* `create_time` - (Optional, String) Create time.
* `creator_uin` - (Optional, String) Creator User ID.
* `data_source_type` - (Optional, String) Data source type: MYSQL, POSTGRE, ORACLE, SQLSERVER, FTP, HIVE, HDFS, ICEBERG, KAFKA, HBASE, SPARK, TBASE, DB2, DM, GAUSSDB, GBASE, IMPALA, ES, S3_DATAINSIGHT, GREENPLUM, PHOENIX, SAP_HANA, SFTP, OCEANBASE, CLICKHOUSE, KUDU, VERTICA, REDIS, COS, DLC, DORIS, CKAFKA, DTS_KAFKA, S3, CDW, TDSQLC, TDSQL, MONGODB, SYBASE, REST_API, StarRocks, TCHOUSE_X.
* `datasource_id` - (Optional, String) Datasource ID.
* `description` - (Optional, String) Node Description.
* `ext_config` - (Optional, List) Node extension configuration information.
* `id` - (Optional, String) Node ID.
* `name` - (Optional, String) Node Name.
* `node_mapping` - (Optional, List) Node mapping.
* `node_type` - (Optional, String) Node type: INPUT,OUTPUT,JOIN,FILTER,TRANSFORM.
* `operator_uin` - (Optional, String) Operator User ID.
* `owner_uin` - (Optional, String) Owner User ID.
* `project_id` - (Optional, String) Project ID.
* `schema` - (Optional, List) Schema information.
* `task_id` - (Optional, String) The task id to which the node belongs.
* `update_time` - (Optional, String) Update time.

The `properties` object supports the following:

* `name` - (Optional, String) Attributes name.
* `value` - (Optional, String) Attributes value.

The `schema_mappings` object supports the following:

* `sink_schema_id` - (Required, String) Schema ID from sink node.
* `source_schema_id` - (Required, String) Schema ID from source node.

The `schema` object supports the following:

* `id` - (Required, String) Schema ID.
* `name` - (Required, String) Schema name.
* `type` - (Required, String) Schema type.
* `alias` - (Optional, String) Schema alias.
* `comment` - (Optional, String) Schema comment.
* `properties` - (Optional, List) Schema extended attributes.
* `value` - (Optional, String) Schema value.

The `source_schema` object supports the following:

* `id` - (Required, String) Schema ID.
* `name` - (Required, String) Schema name.
* `type` - (Required, String) Schema type.
* `alias` - (Optional, String) Schema alias.
* `comment` - (Optional, String) Schema comment.
* `properties` - (Optional, List) Schema extended attributes.
* `value` - (Optional, String) Schema value.

The `task_info` object supports the following:

* `app_id` - (Optional, String) User App Id.
* `config` - (Optional, List) Task configuration.
* `create_time` - (Optional, String) Create time.
* `creator_uin` - (Optional, String) Creator User ID.
* `data_proxy_url` - (Optional, Set) Data proxy url.
* `execute_context` - (Optional, List) Execute context.
* `executor_group_name` - (Optional, String) Executor group name.
* `executor_id` - (Optional, String) Executor resource ID.
* `ext_config` - (Optional, List) Node extension configuration information.
* `has_version` - (Optional, Bool) Whether the task been submitted.
* `in_long_manager_url` - (Optional, String) InLong manager url.
* `in_long_manager_version` - (Optional, String) InLong manager version.
* `in_long_stream_id` - (Optional, String) InLong stream id.
* `incharge` - (Optional, String) Incharge user.
* `input_datasource_type` - (Optional, String) Input datasource type.
* `instance_version` - (Optional, Int) Instance version.
* `last_run_time` - (Optional, String) The last time the task was run.
* `locked` - (Optional, Bool) Whether the task been locked.
* `locker` - (Optional, String) User locked task.
* `mappings` - (Optional, List) Node mapping.
* `nodes` - (Optional, List) Task Node Information.
* `num_records_in` - (Optional, Int) Number of reads.
* `num_records_out` - (Optional, Int) Number of writes.
* `num_restarts` - (Optional, Int) Times of restarts.
* `operator_uin` - (Optional, String) Operator User ID.
* `output_datasource_type` - (Optional, String) Output datasource type.
* `owner_uin` - (Optional, String) Owner User ID.
* `read_phase` - (Optional, Int) Reading stage, 0: full amount, 1: partial full amount, 2: all incremental.
* `reader_delay` - (Optional, Float64) Read latency.
* `running_cu` - (Optional, Float64) The amount of resources consumed by real-time task.
* `schedule_task_id` - (Optional, String) Task scheduling id (job id such as oceanus or us).
* `status` - (Optional, Int) Task status 1. Not started | Task initialization, 2. Task starting, 3. Running, 4. Paused, 5. Task stopping, 6. Stopped, 7. Execution failed, 8. deleted, 9. Locked, 404. unknown status.
* `stop_time` - (Optional, String) The time the task was stopped.
* `submit` - (Optional, Bool) Whether the task version has been submitted for operation and maintenance.
* `switch_resource` - (Optional, Int) Resource tiering status, 0: in progress, 1: successful, 2: failed.
* `task_alarm_regular_list` - (Optional, Set) Task alarm regular.
* `task_group_id` - (Optional, String) Inlong Task Group ID.
* `task_mode` - (Optional, String) Task display mode, 0: canvas mode, 1: form mode.
* `update_time` - (Optional, String) Update time.
* `workflow_id` - (Optional, String) The workflow id to which the task belongs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `task_id` - Task ID.


## Import

wedata integration_realtime_task can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_integration_realtime_task.example 1776563389209296896#h9d39630a-ae45-4460-90b2-0b093cbfef5d
```

