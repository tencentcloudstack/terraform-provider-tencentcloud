---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster_native_node_pools"
sidebar_current: "docs-tencentcloud-datasource-kubernetes_cluster_native_node_pools"
description: |-
  Use this data source to query detailed information of tke kubernetes cluster_native_node_pools
---

# tencentcloud_kubernetes_cluster_native_node_pools

Use this data source to query detailed information of tke kubernetes cluster_native_node_pools

## Example Usage

```hcl
data "tencentcloud_kubernetes_cluster_native_node_pools" "kubernetes_cluster_native_node_pools" {
  cluster_id = "cls-eyi0erm0"
  filters {
    name   = "NodePoolsName"
    values = ["native_node_pool"]
  }
  filters {
    name   = "NodePoolsId"
    values = ["np-ngjwhdv4"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) ID of the cluster.
* `filters` - (Optional, List) Query filter conditions: NodePoolsName, Filter according to the node pool name, type: String, required: no. NodePoolsId, Filter according to the node pool ID, type: String, required: no. tags, Filter according to the label key value pairs, type: String, required: no. tag:tag-key, Filter according to the label key value pairs, type: String, required: no.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) The attribute name, if there are multiple filters, the relationship between the filters is a logical AND relationship.
* `values` - (Required, Set) Attribute values, if there are multiple values in the same filter, the relationship between values under the same filter is a logical OR relationship.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `node_pools` - Node pool list.
  * `annotations` - Node Annotation List.
    * `name` - Name in the map table.
    * `value` - Value in the map table.
  * `cluster_id` - ID of the cluster.
  * `created_at` - Creation time.
  * `deletion_protection` - Whether to enable deletion protection.
  * `labels` - Node Labels.
    * `name` - Name in the map table.
    * `value` - Value in map table.
  * `life_state` - Node pool status.
  * `name` - Node pool name.
  * `native` - Native node pool creation parameters.
    * `auto_repair` - Whether to enable self-healing ability.
    * `data_disks` - Native node pool data disk list.
      * `auto_format_and_mount` - Whether to automatically format the disk and mount it.
      * `disk_partition` - Mount device name or partition name.
      * `disk_size` - Cloud disk size (G).
      * `disk_type` - Cloud disk type. Valid values: `CLOUD_PREMIUM`: Premium Cloud Storage, `CLOUD_SSD`: cloud SSD disk, `CLOUD_BSSD`: Basic SSD, `CLOUD_HSSD`: Enhanced SSD, `CLOUD_TSSD`: Tremendous SSD, `LOCAL_NVME`: local NVME disk.
      * `encrypt` - Pass in this parameter to create an encrypted cloud disk. The value is fixed to `ENCRYPT`.
      * `file_system` - File system (ext3/ext4/xfs).
      * `kms_key_id` - Customize the key when purchasing an encrypted disk. When this parameter is passed in, the Encrypt parameter is not empty.
      * `mount_target` - Mount directory.
      * `snapshot_id` - Snapshot ID. If passed in, the cloud disk will be created based on this snapshot. The snapshot type must be a data disk snapshot.
      * `throughput_performance` - Cloud disk performance, unit: MB/s. Use this parameter to purchase additional performance for the cloud disk.
    * `enable_autoscaling` - Whether to enable elastic scaling.
    * `health_check_policy_name` - Fault self-healing rule name.
    * `host_name_pattern` - Native node pool hostName pattern string.
    * `instance_charge_prepaid` - Billing configuration for yearly and monthly models.
      * `period` - Postpaid billing cycle, unit (month): 1, 2, 3, 4, 5,, 6, 7, 8, 9, 10, 11, 12, 24, 36, 48, 60.
      * `renew_flag` - Prepaid renewal method:
  - `NOTIFY_AND_AUTO_RENEW`: Notify users of expiration and automatically renew (default).
  - `NOTIFY_AND_MANUAL_RENEW`: Notify users of expiration, but do not automatically renew.
  - `DISABLE_NOTIFY_AND_MANUAL_RENEW`: Do not notify users of expiration and do not automatically renew.
    * `instance_charge_type` - Node billing type. `PREPAID` is a yearly and monthly subscription, `POSTPAID_BY_HOUR` is a pay-as-you-go plan. The default is `POSTPAID_BY_HOUR`.
    * `instance_types` - Model list.
    * `internet_accessible` - Public network bandwidth settings.
      * `bandwidth_package_id` - Bandwidth package ID. Note: When ChargeType is BANDWIDTH_PACKAG, the value cannot be empty; otherwise, the value must be empty.
      * `charge_type` - Network billing method. Optional value is `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.
      * `max_bandwidth_out` - Maximum bandwidth output. Note: When chargeType is `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_POSTPAID_BY_HOUR`, the valid range is 1~100. When chargeType is `BANDWIDTH_PACKAG`, the valid range is 1~2000.
    * `key_ids` - Node pool ssh public key id array.
    * `kubelet_args` - Kubelet custom parameters.
    * `lifecycle` - Predefined scripts.
      * `post_init` - Custom script after node initialization.
      * `pre_init` - Custom script before node initialization.
    * `management` - Node pool management parameter settings.
      * `hosts` - Hosts configuration.
      * `kernel_args` - Kernel parameter configuration.
      * `nameservers` - Dns configuration.
    * `replicas` - Desired number of nodes.
    * `runtime_root_dir` - Runtime root directory.
    * `scaling` - Node pool scaling configuration.
      * `create_policy` - Node pool expansion strategy. `ZoneEquality`: multiple availability zones are broken up; `ZonePriority`: the preferred availability zone takes precedence.
      * `max_replicas` - Maximum number of replicas in node pool.
      * `min_replicas` - Minimum number of replicas in node pool.
    * `security_group_ids` - Security group list.
    * `subnet_ids` - Subnet list.
    * `system_disk` - System disk configuration.
      * `disk_size` - Cloud disk size (G).
      * `disk_type` - Cloud disk type.
  * `node_pool_id` - ID of the node pool.
  * `tags` - Node tags.
    * `resource_type` - The resource type bound to the label.
    * `tags` - Tag pair list.
      * `key` - Tag Key.
      * `value` - Tag Value.
  * `taints` - node taint.
    * `effect` - Effect of the taint.
    * `key` - Key of the taint.
    * `value` - Value of the taint.
  * `type` - Node pool type. Optional value is `Native`.
  * `unschedulable` - Whether the node is not schedulable by default.


