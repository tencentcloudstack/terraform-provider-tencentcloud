---
subcategory: "MapReduce(EMR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_emr_service_node_infos"
sidebar_current: "docs-tencentcloud-datasource-emr_service_node_infos"
description: |-
  Use this data source to query detailed information of emr emr_service_node_infos
---

# tencentcloud_emr_service_node_infos

Use this data source to query detailed information of emr emr_service_node_infos

## Example Usage

```hcl
data "tencentcloud_emr_service_node_infos" "emr_service_node_infos" {
  instance_id              = "emr-rzrochgp"
  offset                   = 1
  limit                    = 10
  search_text              = ""
  conf_status              = 2
  maintain_state_id        = 2
  operator_state_id        = 1
  health_state_id          = "2"
  service_name             = "YARN"
  node_type_name           = "master"
  data_node_maintenance_id = 0
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) EMR Instance ID.
* `conf_status` - (Optional, Int) Configuration Status, -2: Configuration Failed, -1: Configuration Expired, 1: Synchronized, -99 All.
* `data_node_maintenance_id` - (Optional, Int) Filter Condition: Whether DN is in Maintenance Mode - 0 represents all statuses, 1 represents in maintenance mode.
* `health_state_id` - (Optional, String) Filter Conditions: Health Status, 0 represents unavailable, 1 represents good, -2 represents unknown, -99 represents all, -3 represents potential risks, -4 represents not detected.
* `limit` - (Optional, Int) Number of Items per Page.
* `maintain_state_id` - (Optional, Int) Filter Condition: Maintenance Status - 0 represents all statuses, 1 represents normal mode, 2 represents maintenance mode.
* `node_type_name` - (Optional, String) Node Names: master, core, task, common, router, all.
* `offset` - (Optional, Int) Page Number.
* `operator_state_id` - (Optional, Int) Filter Condition: Operation Status - 0 represents all statuses, 1 represents started, 2 represents stopped.
* `result_output_file` - (Optional, String) Used to save results.
* `search_fields` - (Optional, List) Search Fields.
* `search_text` - (Optional, String) Search Field.
* `service_name` - (Optional, String) Service Component Name, all in uppercase, e.g., YARN.

The `search_fields` object supports the following:

* `search_type` - (Required, String) Types Supported for Search.
* `search_value` - (Required, String) Values Supported for Search.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `alias_info` - Serialization of Aliases for All Nodes in the Cluster.
* `service_node_list` - Service Node Detail Information.
  * `conf_group_id` - Configuration Group ID.
  * `conf_group_name` - Configuration Group Name.
  * `conf_status` - Configuration Status.
  * `data_node_maintenance_state` - Data Node Maintenance State.
  * `flag` - Flag.
  * `ha_state` - HA State.
  * `health_status` - Process Health Status.
    * `code` - Health Status Code.
    * `desc` - Health Status Description.
    * `text` - Health Status Description.
  * `ip` - The IP address of the node where the process resides.
  * `is_federation` - Whether Federation is Supported.
  * `is_support_role_monitor` - Whether Monitoring is Supported.
  * `last_restart_time` - Most Recent Restart Time.
  * `monitor_status` - Monitor Status.
  * `name_service` - Name Service.
  * `node_flag_filter` - Node Flag Filter.
  * `node_name` - Node Name.
  * `node_type` - Node Type.
  * `ports_info` - Process Port Information.
  * `service_detection_info` - Process Detection Information.
    * `detect_alert` - Detection Alert Level.
    * `detect_function_key` - Detection Function Description.
    * `detect_function_value` - Detection Function Result.
    * `detect_time` - Detection Time.
  * `service_status` - Service Status.
  * `status` - Status.
  * `stop_policies` - Stop Policy.
    * `batch_size_range` - Batch  Node Count Optional Range.
    * `describe` - Policy Description.
    * `display_name` - Policy Display Name.
    * `is_default` - Whether it is the Default Policy.
    * `name` - Policy Name.
* `support_node_flag_filter_list` - Supported FlagNode List.
* `total_cnt` - Total Count.


