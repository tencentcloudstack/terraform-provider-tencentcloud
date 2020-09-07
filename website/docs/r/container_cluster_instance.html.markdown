---
subcategory: "Container Cluster"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_container_cluster_instance"
sidebar_current: "docs-tencentcloud-resource-container_cluster_instance"
description: |-
  Provides a TencentCloud Container Cluster Instance resource.
---

# tencentcloud_container_cluster_instance

Provides a TencentCloud Container Cluster Instance resource.

~> **NOTE:** It has been deprecated and replaced by  tencentcloud_kubernetes_scale_worker.

## Example Usage

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

* `bandwidth_type` - (Required) The network type of the node.
* `bandwidth` - (Required) The network bandwidth of the node.
* `cluster_id` - (Required) The id of the cluster.
* `is_vpc_gateway` - (Required) Describe whether the node enable the gateway capability.
* `root_size` - (Required) The size of the root volume.
* `storage_size` - (Required) The size of the data volume.
* `subnet_id` - (Required) The subnet id which the node stays in.
* `zone_id` - (Required) The zone which the node stays in.
* `cpu` - (Optional, **Deprecated**) It has been deprecated from version 1.16.0. Set 'instance_type' instead. The cpu of the node.
* `cvm_type` - (Optional) The type of node needed by cvm.
* `docker_graph_path` - (Optional) The docker graph path is going to mounted.
* `instance_name` - (Optional) The name ot node.
* `instance_type` - (Optional) The instance type of the node needed by cvm.
* `key_id` - (Optional) The key_id of each node(if using key pair to access).
* `mem` - (Optional, **Deprecated**) It has been deprecated from version 1.16.0. Set 'instance_type' instead. The memory of the node.
* `mount_target` - (Optional) The path which volume is going to be mounted.
* `password` - (Optional) The password of each node.
* `period` - (Optional) The puchase duration of the node needed by cvm.
* `require_wan_ip` - (Optional) Indicate whether wan ip is needed.
* `root_type` - (Optional) The type of the root volume. see more from CVM.
* `sg_id` - (Optional) The security group id.
* `storage_type` - (Optional) The type of the data volume. see more from CVM.
* `unschedulable` - (Optional) Determine whether the node will be schedulable. 0 is the default meaning node will be schedulable. 1 for unschedulable.
* `user_script` - (Optional) User defined script in a base64-format. The script runs after the kubernetes component is ready on node. see more from CCS api documents.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `abnormal_reason` - Describe the reason when node is in abnormal state(if it was).
* `instance_id` - An id identify the node, provided by cvm.
* `is_normal` - Describe whether the node is normal.
* `lan_ip` - Describe the lan ip of the node.
* `wan_ip` - Describe the wan ip of the node.


