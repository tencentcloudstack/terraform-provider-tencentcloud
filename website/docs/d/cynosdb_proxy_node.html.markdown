---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_proxy_node"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_proxy_node"
description: |-
  Use this data source to query detailed information of cynosdb proxy_node
---

# tencentcloud_cynosdb_proxy_node

Use this data source to query detailed information of cynosdb proxy_node

## Example Usage

```hcl
data "tencentcloud_cynosdb_proxy_node" "proxy_node" {
  order_by      = "CREATETIME"
  order_by_type = "DESC"
  filters {
    names       = "ClusterId"
    values      = "cynosdbmysql-cgd2gpwr"
    exact_match = false
    name        = "ClusterId"
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Search criteria, if there are multiple filters, the relationship between the filters is a logical AND relationship.
* `order_by_type` - (Optional, String) Sort type, value range:ASC: ascending sort; DESC: descending sort.
* `order_by` - (Optional, String) Sort field, value range:CREATETIME: creation time; PRIODENDTIME: expiration time.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `names` - (Required, Set) Search String.
* `values` - (Required, Set) Search String.
* `exact_match` - (Optional, Bool) Exact match or not.
* `name` - (Optional, String) Search Fields. Supported: Status, ProxyNodeId, ClusterId.
* `operator` - (Optional, String) Operator.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `proxy_node_infos` - Database Agent Node List.
  * `app_id` - User AppID.
  * `cluster_id` - Cluster ID.
  * `cpu` - Database Agent Node CPU.
  * `mem` - Database Agent Node Memory.
  * `proxy_group_id` - Database Agent Group ID.
  * `proxy_node_connections` - The current number of connections of the node. The DescribeProxyNodes interface does not return a value for this field.
  * `proxy_node_id` - Database Agent Node ID.
  * `region` - region.
  * `status` - Database Agent Node Status.
  * `zone` - Availability Zone.


