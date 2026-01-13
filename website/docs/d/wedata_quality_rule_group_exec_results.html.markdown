---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_quality_rule_group_exec_results"
sidebar_current: "docs-tencentcloud-datasource-wedata_quality_rule_group_exec_results"
description: |-
  Use this data source to query detailed information of wedata quality rule group exec results
---

# tencentcloud_wedata_quality_rule_group_exec_results

Use this data source to query detailed information of wedata quality rule group exec results

## Example Usage

```hcl
data "tencentcloud_wedata_quality_rule_group_exec_results" "wedata_quality_rule_group_exec_results" {
  project_id = "1840731342293087232"
  filters {
    name   = "Status"
    values = ["3"]
  }
  order_fields {
    name      = "UpdateTime"
    direction = "DESC"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `filters` - (Optional, List) Filter conditions. Supported filters:
1. GroupType - Rule group type: DEFAULT (default type), WORKFLOW_NODE (workflow node type)
2. InstanceId - Rule group execution instance ID
3. ParentInstanceId - Parent instance ID
4. LifeCycleRunNum - Lifecycle run number
5. InstanceStatus - Instance status: Waiting (INITIAL, EVENT_LISTENING, DEPENDENCE, ALLOCATED, LAUNCHED, BEFORE_ASPECT, ISSUED), Running (RUNNING, AFTER_ASPECT, WAITING_AFTER_ASPECT), Failed (FAILED, EXPIRED, KILL, KILLING, PENDING), Success (COMPLETED)
6. DatasourceId - Data source ID
7. DatasourceType - Data source type: 1-MYSQL, 2-HIVE, 3-DLC, 4-GBASE, 5-TCHouse-P/CDW, 6-ICEBERG, 7-DORIS, 8-TCHouse-D, 9-EMR_STARROCKS, 10-TBDS_STARROCKS, 11-TCHouse-X
8. DatabaseName - Database name
9. DatabaseId - Database ID
10. SchemaName - Schema name
11. ReceiverFlag - Whether it is the current user's subscription: true/false
12. TableName - Table name (supports fuzzy matching)
13. RuleGroupName - Rule group name
14. RuleGroupExecId - Rule group execution ID
15. RuleGroupTableId - Rule group table ID
16. Keyword - Keyword search (rule group execution ID, table name, table owner)
17. StartTime - Actual run start time (Unix timestamp)
18. EndTime - Actual run end time (Unix timestamp)
19. ScheduledStartTime - Scheduled start time (Unix timestamp)
20. ScheduledEndTime - Scheduled end time (Unix timestamp)
21. DsJobId - Data source job ID
22. TriggerType - Trigger type: 1-Manual, 2-Schedule, 3-Periodic
23. Status - Rule group execution status: 0-Initial, 1-Submitted, 2-Detecting, 3-Normal, 4-Abnormal, 5-Delivering, 6-Execution error, 7-Not detected
24. TableIds - Table ID collection
25. RuleGroupId - Rule group ID
26. BizCatalogIds - Business catalog ID
27. CatalogName - Data catalog name.
* `order_fields` - (Optional, List) Sort fields. Supported fields: CreateTime (sort by creation time), UpdateTime (sort by update time, default). Sort direction: 1-ASC, 2-DESC.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Optional, String) Filter field name.
* `values` - (Optional, Set) Filter value list.

The `order_fields` object supports the following:

* `direction` - (Required, String) Sort direction: ASC|DESC.
* `name` - (Required, String) Sort field name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Rule group execution result list.


