---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_quality_rule_group"
sidebar_current: "docs-tencentcloud-resource-wedata_quality_rule_group"
description: |-
  Provides a resource to create a wedata quality rule group
---

# tencentcloud_wedata_quality_rule_group

Provides a resource to create a wedata quality rule group

## Example Usage

```hcl
resource "tencentcloud_wedata_quality_rule_group" "quality_rule_group" {
  project_id = "1612982498218618880"

  rule_group_exec_strategy_bo_list {
    monitor_type        = 3
    executor_group_id   = "20220509_114220_111"
    rule_group_name     = "tf-test-quality-rule-group"
    database_name       = "tf_test_db"
    datasource_id       = "20220509_114220_111"
    table_name          = "tf_test_table"
    exec_queue          = "default"
    executor_group_name = "default-group"

    start_time = "2024-01-01 00:00:00"
    end_time   = "2024-12-31 23:59:59"
    cycle_type = "D"
    delay_time = 0
    cycle_step = 1

    description        = "Test quality rule group"
    schedule_time_zone = "UTC+8"
    exec_engine_type   = "SPARK"

    group_config {
      analysis_type = "SNAPSHOT"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `rule_group_exec_strategy_bo_list` - (Required, List) Task parameters.

The `group_config` object of `rule_group_exec_strategy_bo_list` supports the following:

* `analysis_type` - (Optional, String) Analysis type, optional values: `INFERENCE`-inference table; `TIME_SERIES`-time series table; `SNAPSHOT`-snapshot table.
* `base_db` - (Optional, String) Base database.
* `base_table` - (Optional, String) Base table.
* `comparison_column_type` - (Optional, String) Comparison column type.
* `comparison_column` - (Optional, String) Comparison column.
* `feature_column` - (Optional, String) Feature column.
* `granularity_type` - (Optional, String) Metric granularity unit.
* `granularity` - (Optional, Int) Metric granularity.
* `label_column_type` - (Optional, String) Label column type.
* `label_column` - (Optional, String) Label column.
* `model_id_column_type` - (Optional, String) Model ID column type.
* `model_id_column` - (Optional, String) Model ID column.
* `model_monitor_type` - (Optional, String) Model detection type, required when analysis type is inference table (INFERENCE), optional values: `CLAASSIFICATION`-classification; `REGRESSION`-regression.
* `positive_value` - (Optional, String) Positive class value.
* `predict_column_type` - (Optional, String) Prediction column type.
* `predict_column` - (Optional, String) Prediction column.
* `protection_value` - (Optional, String) Protection group.
* `timestamp_column_type` - (Optional, String) Timestamp column type.
* `timestamp_column` - (Optional, String) Timestamp column.

The `rule_group_exec_strategy_bo_list` object supports the following:

* `database_name` - (Required, String) Database name.
* `datasource_id` - (Required, String) Data source ID.
* `executor_group_id` - (Required, String) Execution resource group ID.
* `monitor_type` - (Required, Int) Monitor type `2`. Associated production scheduling, `3`. Offline periodic detection.
* `rule_group_name` - (Required, String) Monitor task name.
* `table_name` - (Required, String) Table name.
* `catalog_name` - (Optional, String) Data catalog name, defaults to DataLakeCatalog if not filled (this parameter is invalid when updating quality monitoring).
* `cycle_step` - (Optional, Int) Interval, required when MonitorType=3, indicates the interval time of periodic tasks; Week/Month/Day tasks can choose: `1`; Minute tasks can choose: `10`, `20`, `30`; Hour tasks can choose: `1`, `2`, `3`, `4`, `6`, `8`, `12`.
* `cycle_type` - (Optional, String) Scheduling cycle type, required when MonitorType=3, specific values: `I`: Schedule by minute; `H`: Schedule by hour; `D`: Schedule by day; `W`: Schedule by week; `M`: Schedule by month.
* `delay_time` - (Optional, Int) Delayed scheduling time, required when MonitorType=3, mainly used for day/week/month tasks, measured in minutes. For example, if a day task needs to be delayed to 02:00, this field value is 120, indicating a delay of 2 hours (120 minutes). For hour/minute tasks, this field is meaningless, fill in fixed value 0, otherwise field validation will fail.
* `description` - (Optional, String) Task description.
* `dlc_group_name` - (Optional, String) When data source is DLC, corresponds to DLC resource group. According to the DLC engine name filled in ExecQueue, select the resource group under the corresponding engine.
* `end_time` - (Optional, String) Cycle end time, required when MonitorType=3.
* `engine_param` - (Optional, String) Engine parameters.
* `exec_engine_type` - (Optional, String) Running execution engine, if not passed, will request the default execution engine under this data source.
* `exec_plan` - (Optional, String) Execution plan.
* `exec_queue` - (Optional, String) Compute queue, required when data source is HIVE, ICEBERG, DLC. When data source is DLC, this field should be filled with DLC data engine name.
* `executor_group_name` - (Optional, String) Execution resource group name.
* `group_config` - (Optional, List) Task monitoring parameters.
* `rule_group_id` - (Optional, Int) Monitor task ID, required when editing and updating monitor tasks.
* `rule_id` - (Optional, Int) Rule ID.
* `rule_name` - (Optional, String) Rule name.
* `schedule_time_zone` - (Optional, String) Time zone, default is UTC+8.
* `schema_name` - (Optional, String) Schema name.
* `start_time` - (Optional, String) Cycle start time, required when MonitorType=3.
* `task_action` - (Optional, String) Time specification, mainly used for week/month scheduling cycle tasks. For week scheduling cycle: means specifying which day of the week to run, multiple options separated by English commas, can fill 1,2...7, representing Sunday, Monday...Saturday respectively, for example fill "1,2", means execute on Sunday and Monday; For month scheduling cycle: means specifying which day of the month to run, multiple options separated by English commas, can fill 1,2,...,31, representing 1st, 2nd...31st respectively, for example fill "1,2", means execute on 1st and 2nd of each month.
* `tasks` - (Optional, List) Associated production scheduling task list, required when MonitorType=2.
* `trigger_types` - (Optional, Set) Trigger type, mainly used for "Associated production scheduling" (MonitorType=2) monitoring tasks, optional values: `CYCLE`: Periodic scheduling; `MAKE_UP`: Backfill; `RERUN`: Rerun.

The `tasks` object of `rule_group_exec_strategy_bo_list` supports the following:

* `task_id` - (Required, String) Production scheduling task ID.
* `task_name` - (Required, String) Production scheduling task name.
* `workflow_id` - (Required, String) Production scheduling task workflow ID.
* `cycle_type` - (Optional, Int) Production scheduling task cycle type.
* `in_charge_id_list` - (Optional, Set) Person in charge ID.
* `in_charge_name_list` - (Optional, Set) Person in charge name.
* `schedule_time_zone` - (Optional, String) Time zone.
* `task_type` - (Optional, String) Production task type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_group_id` - Rule group ID.


