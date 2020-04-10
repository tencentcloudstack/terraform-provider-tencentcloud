---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_container_cluster"
sidebar_current: "docs-tencentcloud-resource-container_cluster"
description: |-
  Provides a TencentCloud Container Cluster resource.
---

# tencentcloud_container_cluster

Provides a TencentCloud Container Cluster resource.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_kubernetes_cluster.

## Example Usage

```hcl
resource "tencentcloud_container_cluster" "foo" {
  cluster_name                 = "terraform-acc-test"
  cpu                          = 1
  mem                          = 1
  os_name                      = "ubuntu16.04.1 LTSx86_64"
  bandwidth                    = 1
  bandwidth_type               = "PayByHour"
  require_wan_ip               = 1
  subnet_id                    = "subnet-abcdabc"
  is_vpc_gateway               = 0
  storage_size                 = 0
  root_size                    = 50
  goods_num                    = 1
  password                     = "Admin12345678"
  vpc_id                       = "vpc-abcdabc"
  cluster_cidr                 = "10.0.2.0/24"
  ignore_cluster_cidr_conflict = 0
  cvm_type                     = "PayByHour"
  cluster_desc                 = "foofoofoo"
  period                       = 1
  zone_id                      = 100004
  instance_type                = "S2.SMALL1"
  mount_target                 = ""
  docker_graph_path            = ""
  instance_name                = "bar-vm"
  cluster_version              = "1.7.8"
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth_type` - (Required) The network type of the node.
* `bandwidth` - (Required) The network bandwidth of the node.
* `cluster_cidr` - (Required) The CIDR which the cluster is going to use.
* `cluster_name` - (Required) The name of the cluster.
* `goods_num` - (Required) The node number is going to create in the cluster.
* `instance_type` - (Required) The instance type of the node needed by cvm.
* `is_vpc_gateway` - (Required) Describe whether the node enable the gateway capability.
* `os_name` - (Required) The system os name of the node.
* `root_size` - (Required) The size of the root volume.
* `storage_size` - (Required) The size of the data volume.
* `subnet_id` - (Required) The subnet id which the node stays in.
* `vpc_id` - (Required) Specify vpc which the node(s) stay in.
* `zone_id` - (Required) The zone which the node stays in.
* `cluster_desc` - (Optional) The description of the cluster.
* `cluster_version` - (Optional) The kubernetes version of the cluster.
* `cpu` - (Optional, **Deprecated**) It has been deprecated from version 1.16.0. Set 'instance_type' instead. The cpu of the node.
* `cvm_type` - (Optional) The type of node needed by cvm.
* `docker_graph_path` - (Optional) The docker graph path is going to mounted.
* `instance_name` - (Optional) The name ot node.
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
* `kubernetes_version` - The kubernetes version of the cluster.
* `nodes_num` - The node number of the cluster.
* `nodes_status` - The node status of the cluster.
* `total_cpu` - The total cpu of the cluster.
* `total_mem` - The total memory of the cluster.


