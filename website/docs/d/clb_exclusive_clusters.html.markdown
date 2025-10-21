---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_exclusive_clusters"
sidebar_current: "docs-tencentcloud-datasource-clb_exclusive_clusters"
description: |-
  Use this data source to query detailed information of clb exclusive_clusters
---

# tencentcloud_clb_exclusive_clusters

Use this data source to query detailed information of clb exclusive_clusters

## Example Usage

```hcl
data "tencentcloud_clb_exclusive_clusters" "exclusive_clusters" {
  filters {
    name   = "zone"
    values = ["ap-guangzhou-1"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter to query the list of AZ resources as detailed below: cluster-type - String - Required: No - (Filter condition) Filter by cluster type, such as TGW. cluster-id - String - Required: No - (Filter condition) Filter by cluster ID, such as tgw-xxxxxxxx. cluster-name - String - Required: No - (Filter condition) Filter by cluster name, such as test-xxxxxx. cluster-tag - String - Required: No - (Filter condition) Filter by cluster tag, such as TAG-xxxxx. vip - String - Required: No - (Filter condition) Filter by vip in the cluster, such as x.x.x.x. network - String - Required: No - (Filter condition) Filter by cluster network type, such as Public or Private. zone - String - Required: No - (Filter condition) Filter by cluster zone, such as ap-guangzhou-1. isp - String - Required: No - (Filter condition) Filter by TGW cluster isp type, such as BGP. loadblancer-id - String - Required: No - (Filter condition) Filter by loadblancer-id in the cluste, such as lb-xxxxxxxx.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter name.
* `values` - (Required, Set) Filter value array.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_set` - cluster list.
  * `cluster_id` - cluster ID.
  * `cluster_name` - cluster name.
  * `cluster_tag` - Dedicated layer-7 tag. Note: this field may return null, indicating that no valid values can be obtained.
  * `cluster_type` - cluster type: TGW, STGW, VPCGW.
  * `clusters_version` - clusters version.
  * `clusters_zone` - Availability zone where the cluster is located.
    * `master_zone` - Availability master zone where the cluster is located.
    * `slave_zone` - Availability slave zone where the cluster is located.
  * `disaster_recovery_type` - Cluster disaster recovery type:SINGLE-ZONE, DISASTER-RECOVERY, MUTUAL-DISASTER-RECOVERY.
  * `http_max_new_conn` - Maximum number of new http connections.
  * `http_qps` - Http Qps.
  * `https_max_new_conn` - Maximum number of new https connections.
  * `https_qps` - Https Qps.
  * `idle_resource_count` - The total number of free resources in the cluster.
  * `isp` - Isp: BGP, CMCC,CUCC,CTCC,INTERNAL.
  * `load_balance_director_count` - Total number of forwarders in the cluster.
  * `max_conn` - Maximum number of connections.
  * `max_in_flow` - Maximum incoming Bandwidth.
  * `max_in_pkg` - Maximum incoming packet.
  * `max_new_conn` - Maximum number of new connections.
  * `max_out_flow` - Maximum output bandwidth.
  * `max_out_pkg` - Maximum output packet.
  * `network` - cluster network type.
  * `resource_count` - The total number of resources in the cluster.
  * `zone` - .


