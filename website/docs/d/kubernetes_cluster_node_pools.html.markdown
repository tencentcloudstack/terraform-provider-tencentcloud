---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster_node_pools"
sidebar_current: "docs-tencentcloud-datasource-kubernetes_cluster_node_pools"
description: |-
  Use this data source to query detailed information of kubernetes cluster_node_pools
---

# tencentcloud_kubernetes_cluster_node_pools

Use this data source to query detailed information of kubernetes cluster_node_pools

## Example Usage

```hcl
data "tencentcloud_kubernetes_cluster_node_pools" "cluster_node_pools" {
  cluster_id = "cls-kzilgv5m"
  filters {
    name   = "NodePoolsName"
    values = ["mynodepool_xxxx"]
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
* `filters` - (Optional, List) NodePoolsName, Filter according to the node pool name, type: String, required: no. NodePoolsId, Filter according to the node pool ID, type: String, required: no. tags, Filter according to the label key value pairs, type: String, required: no. tag:tag-key, Filter according to the label key value pairs, type: String, required: no.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) The attribute name, if there are multiple filters, the relationship between the filters is a logical AND relationship.
* `values` - (Required, Set) Attribute values, if there are multiple values in the same filter, the relationship between values under the same filter is a logical OR relationship.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `node_pool_set` - Node Pool List.
  * `autoscaling_group_id` - ID of autoscaling group.
  * `autoscaling_group_status` - Status information.
  * `cluster_instance_id` - ID of the cluster.
  * `data_disks` - Multi disk data disk mounting information.
    * `auto_format_and_mount` - Whether to automate the format disk and mount it.
    * `disk_partition` - Mount device name or partition name.
    * `disk_size` - Cloud disk size(G).
    * `disk_type` - Cloud disk type.
    * `file_system` - File system(ext3/ext4/xfs).
    * `mount_target` - Mount directory.
  * `deletion_protection` - Remove protection switch.
  * `desired_nodes_num` - Expected number of nodes.
  * `desired_pod_num` - When the cluster belongs to the node podCIDR size customization mode, the node pool needs to have the pod number attribute.
  * `docker_graph_path` - Dockerd --graph specified value, default to /var/lib/docker.
  * `extra_args` - Node configuration.
    * `kubelet` - Kubelet custom parameters.
  * `gpu_args` - GPU driver related parameters.
    * `cuda` - CUDA version information.
      * `name` - GPU driver or CUDA name.
      * `version` - GPU driver or CUDA version.
    * `cudnn` - CuDNN version information.
      * `dev_name` - Dev name of cuDNN.
      * `doc_name` - Doc name of cuDNN.
      * `name` - Name of cuDNN.
      * `version` - Version of cuDNN.
    * `custom_driver` - Custom GPU driver information.
      * `address` - Custom GPU driver address link.
    * `driver` - GPU driver version information.
      * `name` - GPU driver or CUDA name.
      * `version` - GPU driver or CUDA version.
    * `mig_enable` - Is the MIG feature enabled.
  * `image_id` - ID of image.
  * `labels` - Labels of the node pool.
    * `name` - Name in the map table.
    * `value` - Value in the map table.
  * `launch_configuration_id` - ID of launch configuration.
  * `life_state` - Life cycle state of the node pool, include: creating, normal, updating, deleting, deleted.
  * `max_nodes_num` - Maximum number of nodes.
  * `min_nodes_num` - Minimum number of nodes.
  * `name` - Name of the node pool.
  * `node_count_summary` - Node List.
    * `autoscaling_added` - Automatically managed nodes.
      * `initializing` - Number of nodes in initialization.
      * `joining` - Number of nodes joining.
      * `normal` - Normal number of nodes.
      * `total` - Total number of nodes.
    * `manually_added` - Manually managed nodes.
      * `initializing` - Number of nodes in initialization.
      * `joining` - Number of nodes joining.
      * `normal` - Normal number of nodes.
      * `total` - Total number of nodes.
  * `node_pool_id` - ID of the node pool.
  * `node_pool_os` - Node Pool OS Name.
  * `os_customize_type` - Mirror version of container.
  * `pre_start_user_script` - User defined script, executed before User Script.
  * `tags` - Resource tags.
    * `key` - Label key.
    * `value` - Label value.
  * `taints` - Labels of the node pool.
    * `effect` - Effect of taints mark.
    * `key` - Key of taints mark.
    * `value` - Value of taints mark.
  * `unschedulable` - Is it not schedulable.
  * `user_script` - User defined scripts.


