---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud-container-cluster"
sidebar_current: "docs-tencentcloud-container-cluster-x"
description: |-
  Provides a TencentCloud Container Cluster resource.
---

# tencentcloud_container_cluster

Provides a Container Cluster resource.

## Example Usage

Basic Usage

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

* `cluster_name` - (Required) The name of the cluster. 
* `cpu` - (Required) The cpu of the node. 
* `mem` - (Required) The memory of the node. 
* `os_name` - (Required) The system os name of the node. 
* `bandwidth` - (Required) The network bandwidth of the node. 
* `bandwidth_type` - (Required) The network type of the node. 
* `subnet_id` - (Required) The subnet id which the node stays in. 
* `is_vpc_gateway` - (Required) Describe whether the node enable the gateway capability. 
* `storage_size` - (Required) The size of the data volumn.
* `storage_type` - (Optional) The type of the data volumn. see more from CVM.
* `root_size` - (Required) The size of the root volumn.
* `root_type` - (Optional) The type of the root volumn. see more from CVM.
* `goods_num` - (Required) The node number is going to create in the cluster. 
* `vpc_id` - (Required) Specify vpc which the node(s) stay in. 
* `cluster_cidr` - (Required) The CIDR which the cluster is going to use. 
* `cluster_desc` - (Optional) The description of the cluster. 
* `cvm_type` - (Optional) The type of node needed by cvm. 
* `period` - (Optional) The puchase duration of the node needed by cvm. 
* `zone_id` - (Required) The zone which the node stays in. 
* `instance_type` - (Optional) The instance type of the node needed by cvm. 
* `sg_id` - (Optional) The safe-group id. 
* `mount_target` - (Optional) The path which volumn is going to be mounted. 
* `docker_graph_path` - (Optional) The docker graph path is going to mounted. 
* `instance_name` - (Optional) The name ot node. 
* `cluster_version` - (Optional) The kubernetes version of the cluster. 
* `password` - (Optional) The password of each node. 
* `key_id` - (Optional) The key_id of each node(if using key pair to access).
* `require_wan_ip` - (Optional) Indicate whether wan ip is needed.
* `user_script` - (Optional) User defined script in a base64-format. The script runs after the kubernetes component is ready on node. see more from CCS api documents.

## Attributes Reference

The following attributes are exported:

* `kubernetes_version` - The kubernetes version of the cluster
* `nodes_num` - The node number of the cluster
* `nodes_status` - The node status of the cluster
* `total_cpu` - The total cpu of the cluster
* `total_mem` - The total memory of the cluster
