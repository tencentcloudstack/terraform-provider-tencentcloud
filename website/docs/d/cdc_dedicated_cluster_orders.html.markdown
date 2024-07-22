---
subcategory: "CDC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdc_dedicated_cluster_orders"
sidebar_current: "docs-tencentcloud-datasource-cdc_dedicated_cluster_orders"
description: |-
  Use this data source to query detailed information of CDC dedicated cluster orders
---

# tencentcloud_cdc_dedicated_cluster_orders

Use this data source to query detailed information of CDC dedicated cluster orders

## Example Usage

### Query all orders

```hcl
data "tencentcloud_cdc_dedicated_cluster_orders" "orders" {}
```

### Query orders by filter

```hcl
data "tencentcloud_cdc_dedicated_cluster_orders" "orders1" {
  dedicated_cluster_ids = ["cluster-262n63e8"]
}

data "tencentcloud_cdc_dedicated_cluster_orders" "orders3" {
  status      = "PENDING"
  action_type = "CREATE"
}
```

## Argument Reference

The following arguments are supported:

* `action_type` - (Optional, String) Filter by Dedicated Cluster Order Action Type. Allow filter value: CREATE, EXTEND.
* `dedicated_cluster_ids` - (Optional, Set: [`String`]) Filter by Dedicated Cluster ID.
* `result_output_file` - (Optional, String) Used to save results.
* `status` - (Optional, String) Filter by Dedicated Cluster Order Status. Allow filter value: PENDING, INCONSTRUCTION, DELIVERING, DELIVERED, EXPIRED, CANCELLED, OFFLINE.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dedicated_cluster_order_set` - Filter by Dedicated Cluster Order.
  * `action` - Dedicated Cluster Order Action Type.
  * `cpu` - Dedicated Cluster CPU.
  * `create_time` - Dedicated Cluster Order Create time.
  * `dedicated_cluster_id` - Dedicated Cluster ID.
  * `dedicated_cluster_order_id` - Dedicated Cluster Order ID.
  * `dedicated_cluster_order_items` - Dedicated Cluster Order Item List.
    * `compute_format` - Dedicated Cluster Compute Format.
    * `count` - Dedicated Cluster SubOrder Count.
    * `create_time` - Dedicated Cluster Order Create time.
    * `dedicated_cluster_type_id` - Dedicated Cluster ID.
    * `description` - Dedicated Cluster Type Description.
    * `name` - Dedicated Cluster Type Name.
    * `power_draw` - Dedicated Cluster Supported PowerDraw.
    * `sub_order_id` - Dedicated Cluster SubOrder ID.
    * `sub_order_pay_status` - Dedicated Cluster SubOrder Pay Status.
    * `sub_order_status` - Dedicated Cluster Order Status.
    * `supported_instance_family` - Dedicated Cluster Supported Instance Family.
    * `supported_storage_type` - Dedicated Cluster Storage Type.
    * `supported_uplink_speed` - Dedicated Cluster Supported Uplink Speed.
    * `total_cpu` - Dedicated Cluster Total CPU.
    * `total_gpu` - Dedicated Cluster Total GPU.
    * `total_mem` - Dedicated Cluster Total Memory.
    * `type_family` - Dedicated Cluster Type Family.
    * `type_name` - Dedicated Cluster Type Name.
    * `weight` - Dedicated Cluster Supported Weight.
  * `dedicated_cluster_type_id` - Dedicated Cluster Type ID.
  * `gpu` - Dedicated Cluster GPU.
  * `mem` - Dedicated Cluster Memory.
  * `order_status` - Dedicated Cluster Order Status.
  * `order_type` - Dedicated Cluster Order Type.
  * `pay_status` - Dedicated Cluster Order Pay Status.
  * `pay_type` - Dedicated Cluster Order Pay Type.
  * `power_draw` - Dedicated Cluster Supported PowerDraw.
  * `supported_instance_family` - Dedicated Cluster Supported Instance Family.
  * `supported_storage_type` - Dedicated Cluster Storage Type.
  * `supported_uplink_speed` - Dedicated Cluster Supported Uplink Speed.
  * `time_span` - Dedicated Cluster Order Pay Time Span.
  * `time_unit` - Dedicated Cluster Order Pay Time Unit.
  * `weight` - Dedicated Cluster Supported Weight.


