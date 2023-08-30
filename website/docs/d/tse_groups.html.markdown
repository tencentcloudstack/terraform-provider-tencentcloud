---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_groups"
sidebar_current: "docs-tencentcloud-datasource-tse_groups"
description: |-
  Use this data source to query detailed information of tse groups
---

# tencentcloud_tse_groups

Use this data source to query detailed information of tse groups

## Example Usage

```hcl
data "tencentcloud_tse_groups" "groups" {
  gateway_id = "gateway-ddbb709b"
  filters {
    name   = "GroupId"
    values = ["group-013c0d8e"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String) gateway ID.
* `filters` - (Optional, List) filter conditions, valid value:Name,GroupId.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) filter name.
* `values` - (Required, Set) filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - groups information.
  * `gateway_group_list` - group list of gateway.
    * `binding_strategy` - associated strategy informationNote: This field may return null, indicating that a valid value is not available.
      * `config` - auto scaling configurationNote: This field may return null, indicating that a valid value is not available.
        * `auto_scaler_id` - auto scaler IDNote: This field may return null, indicating that a valid value is not available.
        * `create_time` - create timeNote: This field may return null, indicating that a valid value is not available.
        * `enabled` - whether to enable metric auto scalingNote: This field may return null, indicating that a valid value is not available.
        * `max_replicas` - maximum number of replicasNote: This field may return null, indicating that a valid value is not available.
        * `metrics` - metric listNote: This field may return null, indicating that a valid value is not available.
          * `resource_name` - metric resource nameNote: This field may return null, indicating that a valid value is not available.
          * `target_type` - metric target typeNote: This field may return null, indicating that a valid value is not available.
          * `target_value` - metric target valueNote: This field may return null, indicating that a valid value is not available.
          * `type` - metric typeNote: This field may return null, indicating that a valid value is not available.
        * `modify_time` - modify timeNote: This field may return null, indicating that a valid value is not available.
        * `strategy_id` - strategy IDNote: This field may return null, indicating that a valid value is not available.
      * `create_time` - create timeNote: This field may return null, indicating that a valid value is not available.
      * `cron_config` - timing scaling configurationNote: This field may return null, indicating that a valid value is not available.
        * `create_time` - create time.
        * `enabled` - whether to enable timing auto scaling.
        * `modify_time` - modify time.
        * `params` - params of timing auto scaling.
          * `crontab` - cron expression.
          * `period` - period of timing auto scaling.
          * `start_at` - start time.
          * `target_replicas` - target replicas.
        * `strategy_id` - strategy ID.
      * `description` - description of strategyNote: This field may return null, indicating that a valid value is not available.
      * `gateway_id` - gateway IDNote: This field may return null, indicating that a valid value is not available.
      * `max_replicas` - maximum number of replicas.
      * `modify_time` - modify timeNote: This field may return null, indicating that a valid value is not available.
      * `strategy_id` - strategy ID.
      * `strategy_name` - strategy nameNote: This field may return null, indicating that a valid value is not available.
    * `create_time` - group create time.
    * `description` - group description.
    * `gateway_id` - gateway ID.
    * `group_id` - group Id.
    * `internet_max_bandwidth_out` - public network outbound traffic bandwidth.
    * `is_first_group` - whether it is the default group- 0: false.- 1: yes.
    * `modify_time` - modify time.
    * `name` - group name.
    * `node_config` - group node configration.
      * `number` - group node number, 2-50.
      * `specification` - group specification, 1c2g|2c4g|4c8g|8c16g.
    * `status` - group status.
    * `subnet_ids` - subnet IDs.
  * `total_count` - total count.


