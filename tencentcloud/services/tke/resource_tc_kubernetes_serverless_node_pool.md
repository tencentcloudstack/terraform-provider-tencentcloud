Provide a resource to create serverless node pool of cluster.

Example Usage

Add serverless node pool to a cluster

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "example_cluster_cidr" {
  default = "10.31.0.0/16"
}

locals {
  vpc_id    = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
  sg_id     = data.tencentcloud_security_groups.sg.security_groups.0.security_group_id
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

data "tencentcloud_security_groups" "sg" {
  name = "default"
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                  = local.vpc_id
  cluster_cidr            = var.example_cluster_cidr
  cluster_max_pod_num     = 32
  cluster_name            = "tf_example_cluster"
  cluster_desc            = "tf example cluster"
  cluster_max_service_num = 32
  cluster_version         = "1.18.4"
  cluster_deploy_type     = "MANAGED_CLUSTER"
}

resource "tencentcloud_kubernetes_serverless_node_pool" "example" {
  cluster_id = tencentcloud_kubernetes_cluster.example.id
  name       = "tf_example_serverless_node_pool"

  serverless_nodes {
    display_name = "tf_example_serverless_node1"
    subnet_id    = local.subnet_id
  }

  serverless_nodes {
    display_name = "tf_example_serverless_node2"
    subnet_id    = local.subnet_id
  }

  security_group_ids = [local.sg_id]
  labels = {
    "label1" : "value1",
    "label2" : "value2",
  }
}
```

Adding taints to the virtual nodes under this node pool

The pods without appropriate tolerations will not be scheduled on this node. Refer [taint-and-toleration](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/) for more details.

```hcl
resource "tencentcloud_kubernetes_serverless_node_pool" "example" {
  cluster_id = tencentcloud_kubernetes_cluster.example.id
  name       = "tf_example_serverless_node_pool"

  serverless_nodes {
    display_name = "tf_example_serverless_node1"
    subnet_id    = local.subnet_id
  }

  serverless_nodes {
    display_name = "tf_example_serverless_node2"
    subnet_id    = local.subnet_id
  }

  security_group_ids = [local.sg_id]
  labels = {
    "label1" : "value1",
    "label2" : "value2",
  }

  taints{
    key = "key1"
    value = "value1"
    effect = "NoSchedule"
  }

  taints{
    key = "key1"
    value = "value1"
    effect = "NoExecute"
  }
}
```

Import

serverless node pool can be imported, e.g.

```
$ terraform import tencentcloud_kubernetes_serverless_node_pool.test cls-xxx#np-xxx
```