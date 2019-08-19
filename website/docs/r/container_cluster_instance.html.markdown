---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud-container-cluster-instance"
sidebar_current: "docs-tencentcloud-container-cluster-instance"
description: |-
  Provides a TencentCloud Container Cluster Instance resource.
---

# tencentcloud_container_cluster_instance

Provides a Container Cluster Instance resource.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_kubernetes_cluster and tencentcloud_kubernetes_scale_worker.
## Example Usage

Basic Usage

```hcl
resource "tencentcloud_container_cluster_instance" "bar_instance" {
  cpu               = 1
  mem               = 1
  bandwidth         = 1
  bandwidth_type    = "PayByHour"
  require_wan_ip    = 1
  is_vpc_gateway    = 0
  storage_size      = 10
  root_size         = 50
  password          = "Admin12345678"
  cvm_type          = "PayByMonth"
  period            = 1
  zone_id           = 100004
  instance_type     = "CVM.S2"
  mount_target      = "/data"
  docker_graph_path = ""
  subnet_id         = "subnet-abcdedf"
  cluster_id        = "cls-abcdef"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The id of the cluster.
* `cpu` - (Required) The cpu of the node.
* `mem` - (Required) The memory of the node.
* `bandwidth` - (Required) The network bandwidth of the node.
* `bandwidth_type` - (Required) The network type of the node.
* `require_wan_ip` - (Optional) Indicate whether wan ip is needed.
* `subnet_id` - (Required) The subnet id which the node stays in.
* `is_vpc_gateway` - (Required) Describe whether the node enable the gateway capability.
* `storage_size` - (Required) The size of the data volumn.
* `storage_type` - (Optional) The type of the data volumn. see more from CVM.
* `root_size` - (Required) The size of the root volumn.
* `root_type` - (Optional) The type of the root volumn. see more from CVM.
* `vpc_id` - (Required) Specify vpc which the node(s) stay in.
* `cvm_type` - (Optional) The type of node needed by cvm.
* `period` - (Optional) The puchase duration of the node needed by cvm.
* `zone_id` - (Required) The zone which the node stays in.
* `instance_type` - (Optional) The instance type of the node needed by cvm.
* `sg_id` - (Optional) The safe-group id.
* `mount_target` - (Optional) The path which volumn is going to be mounted.
* `docker_graph_path` - (Optional) The docker graph path is going to mounted.
* `password` - (Optional) The password of each node.
* `key_id` - (Optional) The key_id of each node(if using key pair to access).
* `unschedulable` - (Optional) Determine whether the node will be schedulable. 0 is the default meaning node will be schedulable. 1 for unschedulable.
* `user_script` - (Optional) User defined script in a base64-format. The script runs after the kubernetes component is ready on node. see more from CCS api documents.

## Attributes Reference

The following attributes are exported:

* `abnormal_reason` - Describe the reason when node is in abnormal state(if it was).
* `instance_id` - An id identify the node, provided by cvm.
* `is_normal` - Describe whether the node is normal.
* `wan_ip` - Describe the wan ip of the node.
* `lan_ip` - Describe the lan ip of the node.
