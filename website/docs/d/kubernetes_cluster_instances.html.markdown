---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster_instances"
sidebar_current: "docs-tencentcloud-datasource-kubernetes_cluster_instances"
description: |-
  Use this data source to query detailed information of kubernetes cluster_instances
---

# tencentcloud_kubernetes_cluster_instances

Use this data source to query detailed information of kubernetes cluster_instances

## Example Usage

```hcl
data "tencentcloud_kubernetes_cluster_instances" "cluster_instances" {
  cluster_id    = "cls-ely08ic4"
  instance_ids  = ["ins-kqmx8dm2"]
  instance_role = "WORKER"
  filters {
    name   = "nodepool-id"
    values = ["np-p4e6whqu"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) ID of the cluster.
* `filters` - (Optional, List) List of filter conditions. The optional values of Name are `nodepool-id` and `nodepool-instance-type`. Name is `nodepool-id`, which means filtering machines based on node pool id, and Value is the specific node pool id. Name is `nodepool-instance-type`, which indicates how the node is added to the node pool. Value is MANUALLY_ADDED (manually added to the node pool), AUTOSCALING_ADDED (joined by scaling group expansion method), ALL (manually join the node pool and join the node pool through scaling group expansion).
* `instance_ids` - (Optional, Set: [`String`]) List of node instance IDs to be obtained. If it is empty, it means pulling all node instances in the cluster.
* `instance_role` - (Optional, String) Node role, MASTER, WORKER, ETCD, MASTER_ETCD,ALL, default is WORKER.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) The attribute name, if there are multiple filters, the relationship between the filters is a logical AND relationship.
* `values` - (Required, Set) Attribute values, if there are multiple values in the same filter, the relationship between values under the same filter is a logical OR relationship.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_set` - List of instances in the cluster.
  * `autoscaling_group_id` - Auto scaling group ID.
  * `created_time` - Add time.
  * `drain_status` - Whether the instance is blocked.
  * `failed_reason` - Reasons for instance exception (or being initialized).
  * `instance_advanced_settings` - Node configuration.
    * `data_disks` - Multi-disk data disk mounting information.
      * `auto_format_and_mount` - Whether to automatically format the disk and mount it.
      * `disk_partition` - Mount device name or partition name, required when and only when adding an existing node.
      * `disk_size` - Cloud disk size (G).
      * `disk_type` - Cloud disk type.
      * `file_system` - File system (ext3/ext4/xfs).
      * `mount_target` - Mount directory.
    * `desired_pod_number` - When the node belongs to the podCIDR size customization mode, you can specify the upper limit of the number of pods running on the node.
    * `docker_graph_path` - Dockerd --graph specifies the value, the default is /var/lib/docker.
    * `extra_args` - Node-related custom parameter information.
      * `kubelet` - Kubelet custom parameters.
    * `gpu_args` - GPU driver related parameters, obtain related GPU parameters: https://cloud.tencent.com/document/api/213/15715.
      * `cuda` - CUDA version information.
        * `name` - The name of the GPU driver or CUDA.
        * `version` - GPU driver or CUDA version.
      * `cudnn` - CuDNN version information.
        * `dev_name` - Dev name of cuDNN.
        * `doc_name` - Doc name of cuDNN.
        * `name` - CuDNN name.
        * `version` - Version of cuDNN.
      * `custom_driver` - Custom GPU driver information.
        * `address` - Custom GPU driver address link.
      * `driver` - GPU driver version information.
        * `name` - The name of the GPU driver or CUDA.
        * `version` - GPU driver or CUDA version.
      * `mig_enable` - Whether to enable MIG features.
    * `labels` - Node Label array.
      * `name` - Name in the map table.
      * `value` - Value in map table.
    * `mount_target` - Data disk mount point, the data disk is not mounted by default. Formatted ext3, ext4, xfs file system data disks will be mounted directly. Other file systems or unformatted data disks will be automatically formatted as ext4 (tlinux system formatted as xfs) and mounted. Please pay attention to backing up the data. This setting does not take effect for cloud hosts that have no data disks or multiple data disks.
    * `pre_start_user_script` - Base64 encoded user script, executed before initializing the node, currently only effective for adding existing nodes.
    * `taints` - Node taint.
      * `effect` - Effect of taints mark.
      * `key` - Key of taints mark.
      * `value` - Value of taints mark.
    * `unschedulable` - Set whether the added node participates in scheduling. The default value is 0, which means participating in scheduling; non-0 means not participating in scheduling. After the node initialization is completed, you can execute kubectl uncordon nodename to join the node in scheduling.
    * `user_script` - Base64 encoded userscript.
  * `instance_id` - Instance ID.
  * `instance_role` - Node role, MASTER, WORKER, ETCD, MASTER_ETCD,ALL, default is WORKER.
  * `instance_state` - The status of the instance (running, initializing, failed).
  * `lan_ip` - Node intranet IP.
  * `node_pool_id` - Resource pool ID.


