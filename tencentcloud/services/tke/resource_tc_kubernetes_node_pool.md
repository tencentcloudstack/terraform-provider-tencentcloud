Provide a resource to create an auto scaling group for kubernetes cluster.

~> **NOTE:**  We recommend the usage of one cluster with essential worker config + node pool to manage cluster and nodes. Its a more flexible way than manage worker config with tencentcloud_kubernetes_cluster, tencentcloud_kubernetes_scale_worker or exist node management of `tencentcloud_kubernetes_attachment`. Cause some unchangeable parameters of `worker_config` may cause the whole cluster resource `force new`.

~> **NOTE:**  In order to ensure the integrity of customer data, if you destroy nodepool instance, it will keep the cvm instance associate with nodepool by default. If you want to destroy together, please set `delete_keep_instance` to `false`.

~> **NOTE:**  In order to ensure the integrity of customer data, if the cvm instance was destroyed due to shrinking, it will keep the cbs associate with cvm by default. If you want to destroy together, please set `delete_with_instance` to `true`.

~> **NOTE:**  There are two parameters `wait_node_ready` and `scale_tolerance` to ensure better management of node pool scaling operations. If this parameter is set when creating a resource, the resource will be marked as `tainted` if the set conditions are not met.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "172.31.0.0/16"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

variable "default_instance_type" {
  default = "S1.SMALL1"
}

//this is the cluster with empty worker config
resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = var.cluster_cidr
  cluster_max_pod_num     = 32
  cluster_name            = "tf-tke-unit-test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32
  cluster_version         = "1.18.4"
  cluster_deploy_type     = "MANAGED_CLUSTER"
}

//this is one example of managing node using node pool
resource "tencentcloud_kubernetes_node_pool" "example" {
  name                     = "tf-example"
  cluster_id               = tencentcloud_kubernetes_cluster.example.id
  max_size                 = 6
  min_size                 = 1
  vpc_id                   = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_ids               = [data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id]
  retry_policy             = "INCREMENTAL_INTERVALS"
  desired_capacity         = 4
  enable_auto_scale        = true
  multi_zone_subnet_policy = "EQUALITY"
  node_os                  = "img-9qrfy1xt"

  auto_scaling_config {
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_PREMIUM"
    system_disk_size           = "50"
    orderly_security_group_ids = ["sg-24vswocp"]

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "Password@123"
    enhanced_security_service  = false
    enhanced_monitor_service   = false
    host_name                  = "12.123.0.0"
    host_name_style            = "ORIGINAL"
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }

  taints {
    key    = "test_taint"
    value  = "taint_value"
    effect = "PreferNoSchedule"
  }

  taints {
    key    = "test_taint2"
    value  = "taint_value2"
    effect = "PreferNoSchedule"
  }

  node_config {
    docker_graph_path = "/var/lib/docker"
    extra_args = [
      "root-dir=/var/lib/kubelet"
    ]
  }
}
```

Using Spot CVM Instance

```hcl
resource "tencentcloud_kubernetes_node_pool" "example" {
  name                     = "tf-example"
  cluster_id               = tencentcloud_kubernetes_cluster.managed_cluster.id
  max_size                 = 6
  min_size                 = 1
  vpc_id                   = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_ids               = [data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id]
  retry_policy             = "INCREMENTAL_INTERVALS"
  desired_capacity         = 4
  enable_auto_scale        = true
  multi_zone_subnet_policy = "EQUALITY"

  auto_scaling_config {
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_PREMIUM"
    system_disk_size           = "50"
    orderly_security_group_ids = ["sg-24vswocp", "sg-3qntci2v", "sg-7y1t2wax"]
    instance_charge_type       = "SPOTPAID"
    spot_instance_type         = "one-time"
    spot_max_price             = "1000"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "Password@123"
    enhanced_security_service  = false
    enhanced_monitor_service   = false
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2"
  }
}
```

If instance_type is CBM

```hcl
resource "tencentcloud_kubernetes_node_pool" "example" {
  name                     = "tf-example"
  cluster_id               = "cls-23ieal0c"
  max_size                 = 100
  min_size                 = 1
  vpc_id                   = "vpc-i5yyodl9"
  subnet_ids               = ["subnet-d4umunpy"]
  retry_policy             = "INCREMENTAL_INTERVALS"
  enable_auto_scale        = true
  multi_zone_subnet_policy = "EQUALITY"
  node_os                  = "img-eb30mz89"
  delete_keep_instance     = false

  node_config {
    data_disk {
      disk_type    = "LOCAL_NVME"
      disk_size    = 3570
      file_system  = "ext4"
      mount_target = "/var/lib/data1"
    }

    data_disk {
      disk_type    = "LOCAL_NVME"
      disk_size    = 3570
      file_system  = "ext4"
      mount_target = "/var/lib/data2"
    }
  }

  auto_scaling_config {
    instance_type              = "BMI5.24XLARGE384"
    system_disk_type           = "LOCAL_BASIC"
    system_disk_size           = "440"
    orderly_security_group_ids = ["sg-4z20n68d"]

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "Password@123"
    enhanced_security_service  = false
    enhanced_monitor_service   = false
    host_name                  = "example"
    host_name_style            = "ORIGINAL"
  }
}
```

Wait for all scaling nodes to be ready with wait_node_ready and scale_tolerance parameters. The default maximum scaling timeout is 30 minutes.

```hcl
resource "tencentcloud_kubernetes_node_pool" "example" {
  name                     = "tf-example"
  cluster_id               = tencentcloud_kubernetes_cluster.managed_cluster.id
  max_size                 = 100
  min_size                 = 1
  vpc_id                   = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_ids               = [data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id]
  retry_policy             = "INCREMENTAL_INTERVALS"
  desired_capacity         = 50
  enable_auto_scale        = false
  wait_node_ready          = true
  scale_tolerance          = 90
  multi_zone_subnet_policy = "EQUALITY"
  node_os                  = "img-6n21msk1"
  delete_keep_instance     = false

  auto_scaling_config {
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_PREMIUM"
    system_disk_size           = "50"
    orderly_security_group_ids = ["sg-bw28gmso"]

    data_disk {
      disk_type            = "CLOUD_PREMIUM"
      disk_size            = 50
      delete_with_instance = true
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                  = "test123#"
    enhanced_security_service = false
    enhanced_monitor_service  = false
    host_name                 = "12.123.0.0"
    host_name_style           = "ORIGINAL"
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }

  taints {
    key    = "test_taint"
    value  = "taint_value"
    effect = "PreferNoSchedule"
  }

  taints {
    key    = "test_taint2"
    value  = "taint_value2"
    effect = "PreferNoSchedule"
  }

  node_config {
    docker_graph_path = "/var/lib/docker"
    extra_args = [
      "root-dir=/var/lib/kubelet"
    ]
  }

  timeouts {
    create = "30m"
    update = "30m"
  }
}
```

Create Node pool for CDC cluster

```hcl
resource "tencentcloud_kubernetes_node_pool" "example" {
  name                     = "tf-example"
  cluster_id               = "cls-nhhpsdx8"
  default_cooldown         = 400
  max_size                 = 4
  min_size                 = 1
  desired_capacity         = 2
  vpc_id                   = "vpc-pi5u9uth"
  subnet_ids               = ["subnet-muu9a0gk"]
  retry_policy             = "INCREMENTAL_INTERVALS"
  enable_auto_scale        = true
  multi_zone_subnet_policy = "EQUALITY"
  node_os                  = "img-eb30mz89"
  delete_keep_instance     = true

  node_config {
    data_disk {
      disk_type    = "CLOUD_SSD"
      disk_size    = 50
      file_system  = "ext4"
      mount_target = "/var/lib/data1"
    }
  }

  auto_scaling_config {
    instance_type              = "S5.MEDIUM4"
    instance_charge_type       = "CDCPAID"
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = "100"
    orderly_security_group_ids = ["sg-4z20n68d"]

    data_disk {
      disk_type = "CLOUD_SSD"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "Password@123"
    enhanced_security_service  = false
    enhanced_monitor_service   = false
    host_name                  = "example"
    host_name_style            = "ORIGINAL"
    instance_name              = "example"
    instance_name_style        = "ORIGINAL"
    cdc_id                     = "cluster-262n63e8"
  }

  tags = {
    createBy = "Terraform"
  }
}
```

Import

tke node pool can be imported, e.g.

```
$ terraform import tencentcloud_kubernetes_node_pool.example cls-d2xdg3io#np-380ay1o8
```