---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_cluster_resources"
sidebar_current: "docs-tencentcloud-datasource-clb_cluster_resources"
description: |-
  Use this data source to query detailed information of clb cluster_resources
---

# tencentcloud_clb_cluster_resources

Use this data source to query detailed information of clb cluster_resources

## Example Usage

```hcl
data "tencentcloud_clb_cluster_resources" "cluster_resources" {
  filters {
    name   = "idle"
    values = ["True"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions to query cluster. cluster-id - String - Required: No - (Filter condition) Filter by cluster ID, such as tgw-12345678. vip - String - Required: No - (Filter condition) Filter by loadbalancer vip, such as 192.168.0.1. loadblancer-id - String - Required: No - (Filter condition) Filter by loadblancer ID, such as lbl-12345678. idle - String - Required: No - (Filter condition) Filter by Whether load balancing is idle, such as True, False.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter name.
* `values` - (Required, Set) Filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_resource_set` - Cluster resource set.
  * `cluster_id` - Cluster ID.
  * `cluster_name` - cluster name.
  * `clusters_zone` - clusters zone.
    * `master_zone` - Availability master zone where the cluster is located.
    * `slave_zone` - Availability slave zone where the cluster is located.
  * `idle` - Is it idle.
  * `isp` - Isp.
  * `load_balancer_id` - Loadbalance Id.
  * `vip` - vip.


