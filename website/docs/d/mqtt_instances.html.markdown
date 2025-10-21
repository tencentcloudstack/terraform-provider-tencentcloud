---
subcategory: "TDMQ for MQTT(MQTT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mqtt_instances"
sidebar_current: "docs-tencentcloud-datasource-mqtt_instances"
description: |-
  Use this data source to query detailed information of MQTT instances
---

# tencentcloud_mqtt_instances

Use this data source to query detailed information of MQTT instances

## Example Usage

### Query all mqtt instances

```hcl
data "tencentcloud_mqtt_instances" "example" {}
```

### Query mqtt instances by filters

```hcl
data "tencentcloud_mqtt_instances" "example" {
  filters {
    name   = "InstanceId"
    values = ["mqtt-kngmpg9p"]
  }

  filters {
    name   = "InstanceName"
    values = ["tf-example"]
  }

  filters {
    name   = "InstanceStatus"
    values = ["RUNNING"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Query criteria list, supporting the following fields: InstanceName: cluster name, fuzzy search, InstanceId: cluster ID, precise search, InstanceStatus: cluster status search (RUNNING - Running, CREATING - Creating, MODIFYING - Changing, DELETING - Deleting).
* `result_output_file` - (Optional, String) Used to save results.
* `tag_filters` - (Optional, List) Tag filters.

The `filters` object supports the following:

* `name` - (Required, String) Filter name.
* `values` - (Required, Set) Values.

The `tag_filters` object supports the following:

* `tag_key` - (Optional, String) Tag key.
* `tag_values` - (Optional, Set) Tag values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Instance list.
  * `authorization_policy_limit` - Limit on the number of authorization rules.
  * `client_num_limit` - Number of client connections online.
  * `create_time` - Creation time, millisecond timestamp.
  * `destroy_time` - Pre destruction time, millisecond timestamp.
  * `expiry_time` - Expiration time, millisecond level timestamp.
  * `instance_id` - Instacen ID.
  * `instance_name` - Instacen name.
  * `instance_status` - Instance status. RUNNING- In operation; MAINTAINING- Under Maintenance; ABNORMAL- abnormal; OVERDUE- Arrears of fees; DESTROYED- Deleted; CREATING- Creating in progress; MODIFYING- In the process of transformation; CREATE_FAILURE- Creation failed; MODIFY_FAILURE- Transformation failed; DELETING- deleting.
  * `instance_type` - Instance type. BASIC- Basic Edition; PRO- professional edition; PLATINUM- Platinum version.
  * `max_ca_num` - Maximum CA quota.
  * `max_subscription_per_client` - Maximum number of subscriptions per client.
  * `max_subscription` - Maximum number of subscriptions.
  * `pay_mode` - Billing mode, POSTPAID, pay as you go PREPAID, annual and monthly package.
  * `remark` - Remark.
  * `renew_flag` - Whether to renew automatically. Only the annual and monthly package cluster is effective. 1: Automatic renewal; 0: Non automatic renewal.
  * `sku_code` - Product specifications.
  * `topic_num_limit` - Maximum number of instance topics.
  * `topic_num` - Topic num.
  * `tps_limit` - Elastic TPS current limit value.
  * `version` - Instacen version.


