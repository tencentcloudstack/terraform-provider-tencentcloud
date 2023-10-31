---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_clusters"
sidebar_current: "docs-tencentcloud-datasource-oceanus_clusters"
description: |-
  Use this data source to query detailed information of oceanus clusters
---

# tencentcloud_oceanus_clusters

Use this data source to query detailed information of oceanus clusters

## Example Usage

### Query all clusters

```hcl
data "tencentcloud_oceanus_clusters" "example" {}
```

### Query the specified cluster

```hcl
data "tencentcloud_oceanus_clusters" "example" {
  cluster_ids = ["cluster-5c42n3a5"]
  order_type  = 1
  filters {
    name   = "name"
    values = ["tf_example"]
  }
  work_space_id = "space-2idq8wbr"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_ids` - (Optional, Set: [`String`]) Query one or more clusters by their ID. The maximum number of clusters that can be queried at once is 100.
* `filters` - (Optional, List) The filtering rules.
* `order_type` - (Optional, Int) The sorting rule of the cluster information results. Possible values are 1 (sort by time in descending order), 2 (sort by time in ascending order), and 3 (sort by status).
* `result_output_file` - (Optional, String) Used to save results.
* `work_space_id` - (Optional, String) Workspace SerialId.

The `filters` object supports the following:

* `name` - (Required, String) The field to be filtered.
* `values` - (Required, Set) The filtering values of the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_set` - Cluster list.
  * `app_id` - The user AppID.
  * `arch_generation` - V3 version = 2.Note: This field may return null, indicating that no valid values can be obtained.
  * `auto_renew_flag` - The auto-renewal flag. 0 indicates the default state (the user has not set it, which is the initial state; if the user has enabled the prepaid non-stop privilege, the cluster will be automatically renewed), 1 indicates automatic renewal, and 2 indicates no automatic renewal (set by the user).Note: This field may return null, indicating that no valid values can be obtained.
  * `ccns` - The network.
    * `ccn_id` - The ID of the Cloud Connect Network (CCN), such as ccn-rahigzjd.
    * `subnet_id` - The ID of the subnet.
    * `vpc_id` - The ID of the VPC.
  * `cls_log_name` - The name of the CLS logset of the cluster.Note: This field may return null, indicating that no valid values can be obtained.
  * `cls_log_set` - The CLS logset of the cluster.Note: This field may return null, indicating that no valid values can be obtained.
  * `cls_topic_id` - The CLS topic ID of the cluster.Note: This field may return null, indicating that no valid values can be obtained.
  * `cls_topic_name` - The name of the CLS topic of the cluster.Note: This field may return null, indicating that no valid values can be obtained.
  * `cluster_id` - The ID of the cluster.
  * `cluster_sessions` - Session cluster information.Note: This field may return null, indicating that no valid values can be obtained.
  * `cluster_type` - 0: TKE, 1: EKS.Note: This field may return null, indicating that no valid values can be obtained.
  * `correlations` - Space information.Note: This field may return null, indicating that no valid values can be obtained.
    * `cluster_group_id` - Cluster ID.
    * `cluster_group_serial_id` - Cluster SerialId.
    * `cluster_name` - Cluster name.
    * `project_id_str` - Project ID in string format.Note: This field may return null, indicating that no valid values can be obtained.
    * `project_id` - Project ID.
    * `status` - Binding status. 2 - bound, 1 - unbound.
    * `work_space_id` - Workspace SerialId.
    * `work_space_name` - Workspace name.
  * `create_time` - The time when the cluster was created.
  * `creator_uin` - The creator UIN.
  * `cu_mem` - The memory specification of the CU.
  * `cu_num` - The number of CUs.
  * `customized_dns_enabled` - Value: 0 - not set, 1 - set, 2 - not allowed to set.Note: This field may return null, indicating that no valid values can be obtained.
  * `default_cos_bucket` - The default COS bucket of the cluster.Note: This field may return null, indicating that no valid values can be obtained.
  * `default_log_collect_conf` - The default log collection configuration of the cluster.Note: This field may return null, indicating that no valid values can be obtained.
  * `expire_time` - The expiration time of the cluster. If the cluster does not have an expiration time, this field will be -.Note: This field may return null, indicating that no valid values can be obtained.
  * `free_cu_num` - The number of free CUs.
  * `free_cu` - The number of free CUs at the granularity level.Note: This field may return null, indicating that no valid values can be obtained.
  * `is_need_manage_node` - Front-end distinguishes whether the cluster needs 2CU logic, because historical clusters do not need to be changed. Default is 1. All new clusters need to be changed.Note: This field may return null, indicating that no valid values can be obtained.
  * `isolated_time` - The time when the cluster was isolated. If the cluster has not been isolated, this field will be -.Note: This field may return null, indicating that no valid values can be obtained.
  * `name` - The name of the cluster.
  * `net_environment_type` - The network.
  * `orders` - Order information.Note: This field may return null, indicating that no valid values can be obtained.
    * `auto_renew_flag` - 1 - auto-renewal.Note: This field may return null, indicating that no valid values can be obtained.
    * `compute_cu` - The number of CUs in the final cluster.Note: This field may return null, indicating that no valid values can be obtained.
    * `operate_uin` - UIN of the operator.Note: This field may return null, indicating that no valid values can be obtained.
    * `order_time` - The time of the order.Note: This field may return null, indicating that no valid values can be obtained.
    * `type` - 1 - create, 2 - renew, 3 - scale.Note: This field may return null, indicating that no valid values can be obtained.
  * `owner_uin` - The main account UIN.
  * `pay_mode` - 0 - postpaid, 1 - prepaid.Note: This field may return null, indicating that no valid values can be obtained.
  * `region` - The region where the cluster is located.
  * `remark` - A description of the cluster.
  * `running_cu` - Running CU.Note: This field may return null, indicating that no valid values can be obtained.
  * `seconds_until_expiry` - The number of seconds until the cluster expires. If the cluster does not have an expiration time, this field will be -.Note: This field may return null, indicating that no valid values can be obtained.
  * `sql_gateways` - Gateway information.Note: This field may return null, indicating that no valid values can be obtained.
    * `create_time` - Creation time.Note: This field may return null, indicating that no valid values can be obtained.
    * `creator_uin` - Creator.Note: This field may return null, indicating that no valid values can be obtained.
    * `cu_spec` - CU specification.Note: This field may return null, indicating that no valid values can be obtained.
    * `flink_version` - Flink kernel version.Note: This field may return null, indicating that no valid values can be obtained.
    * `properties` - Configuration parameters.Note: This field may return null, indicating that no valid values can be obtained.
      * `key` - Key of the system configuration.
      * `value` - Value of the system configuration.
    * `resource_refs` - Reference resources.Note: This field may return null, indicating that no valid values can be obtained.
      * `resource_id` - Unique identifier of the resource.
      * `type` - Reference type. 0: user resource.Note: This field may return null, indicating that no valid values can be obtained.
      * `version` - Version number.
      * `workspace_id` - Unique identifier of the space.
    * `serial_id` - Unique identifier.Note: This field may return null, indicating that no valid values can be obtained.
    * `status` - Status. 1 - stopped, 2 - starting, 3 - started, 4 - start failed, 5 - stopping.Note: This field may return null, indicating that no valid values can be obtained.
    * `update_time` - Update time.Note: This field may return null, indicating that no valid values can be obtained.
  * `status_desc` - The status description.
  * `status` - The status of the cluster. Possible values are 1 (uninitialized), 3 (initializing), and 2 (running).
  * `tags` - The tags bound to the cluster.Note: This field may return null, indicating that no valid values can be obtained.
    * `tag_key` - The tag key.Note: This field may return null, indicating that no valid values can be obtained.
    * `tag_value` - The tag value.Note: This field may return null, indicating that no valid values can be obtained.
  * `update_time` - The time of the last operation on the cluster.
  * `version` - The version information of the cluster.Note: This field may return null, indicating that no valid values can be obtained.
    * `flink` - The Flink version of the cluster.Note: This field may return null, indicating that no valid values can be obtained.
    * `supported_flink` - The Flink versions supported by the cluster.Note: This field may return null, indicating that no valid values can be obtained.
  * `zone` - The availability zone.


