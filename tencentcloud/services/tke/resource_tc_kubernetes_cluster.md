Provide a resource to create a kubernetes cluster.

~> **NOTE:** To use the custom Kubernetes component startup parameter function (parameter `extra_args`), you need to submit a ticket for application.

~> **NOTE:** We recommend this usage that uses the `tencentcloud_kubernetes_cluster` resource to create a cluster without any `worker_config`, then adds nodes by the `tencentcloud_kubernetes_node_pool` resource.
It's more flexible than managing worker config directly with `tencentcloud_kubernetes_cluster`, `tencentcloud_kubernetes_scale_worker`, or existing node management of `tencentcloud_kubernetes_attachment`. The reason is that `worker_config` is unchangeable and may cause the whole cluster resource to `ForceNew`.

~> **NOTE:** Executing `terraform destroy` to destroy the resource will default to deleting the node resource, If it is necessary to preserve node instance resources, Please set `instance_delete_mode` to `retain`.

~> **NOTE:** If you want to set up addon for the tke cluster, it is recommended to use resource `tencentcloud_kubernetes_addon`.

Example Usage

Create a basic cluster with two worker nodes

```hcl
variable "default_instance_type" {
  default = "SA2.2XLARGE16"
}

variable "availability_zone_first" {
  default = "ap-guangzhou-3"
}

variable "availability_zone_second" {
  default = "ap-guangzhou-4"
}

variable "example_cluster_cidr" {
  default = "10.31.0.0/16"
}

locals {
  first_vpc_id     = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.vpc_id
  first_subnet_id  = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.subnet_id
  second_vpc_id    = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.vpc_id
  second_subnet_id = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.subnet_id
  sg_id            = tencentcloud_security_group.sg.id
  image_id         = data.tencentcloud_images.default.image_id
}

data "tencentcloud_vpc_subnets" "vpc_one" {
  is_default        = true
  availability_zone = var.availability_zone_first
}

data "tencentcloud_vpc_subnets" "vpc_two" {
  is_default        = true
  availability_zone = var.availability_zone_second
}

resource "tencentcloud_security_group" "sg" {
  name = "tf-example-sg"
}

resource "tencentcloud_security_group_lite_rule" "sg_rule" {
  security_group_id = tencentcloud_security_group.sg.id

  ingress = [
    "ACCEPT#10.0.0.0/16#ALL#ALL",
    "ACCEPT#172.16.0.0/22#ALL#ALL",
    "DROP#0.0.0.0/0#ALL#ALL",
  ]

  egress = [
    "ACCEPT#172.16.0.0/22#ALL#ALL",
  ]
}

data "tencentcloud_images" "default" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                          = local.first_vpc_id
  cluster_cidr                    = var.example_cluster_cidr
  cluster_max_pod_num             = 32
  cluster_name                    = "tf_example_cluster"
  cluster_desc                    = "example for tke cluster"
  cluster_max_service_num         = 32
  cluster_internet                = false
  cluster_internet_security_group = local.sg_id
  cluster_version                 = "1.22.5"
  cluster_deploy_type             = "MANAGED_CLUSTER"

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_first
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.first_subnet_id
    # img_id                     = local.image_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    # key_ids                   = ["skey-11112222"]
    password = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
  }

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_second
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.second_subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    key_ids                   = ["skey-11112222"]
    cam_role_name             = "CVM_QcsRole"
    # password                  = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
}
```

Create an empty cluster with a node pool

The cluster does not have any nodes, nodes will be added through node pool.

```hcl
variable "default_instance_type" {
  default = "SA2.2XLARGE16"
}

variable "availability_zone_first" {
  default = "ap-guangzhou-3"
}

variable "availability_zone_second" {
  default = "ap-guangzhou-4"
}

variable "example_cluster_cidr" {
  default = "10.31.0.0/16"
}

locals {
  first_vpc_id    = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.vpc_id
  first_subnet_id = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.subnet_id
  sg_id    = tencentcloud_security_group.sg.id
}

data "tencentcloud_vpc_subnets" "vpc_one" {
  is_default        = true
  availability_zone = var.availability_zone_first
}

data "tencentcloud_vpc_subnets" "vpc_two" {
  is_default        = true
  availability_zone = var.availability_zone_second
}

resource "tencentcloud_security_group" "sg" {
  name = "tf-example-np-sg"
}

resource "tencentcloud_security_group_lite_rule" "sg_rule" {
  security_group_id = tencentcloud_security_group.sg.id

  ingress = [
    "ACCEPT#10.0.0.0/16#ALL#ALL",
    "ACCEPT#172.16.0.0/22#ALL#ALL",
    "DROP#0.0.0.0/0#ALL#ALL",
  ]

  egress = [
    "ACCEPT#172.16.0.0/22#ALL#ALL",
  ]
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                  = local.first_vpc_id
  cluster_cidr            = var.example_cluster_cidr
  cluster_max_pod_num     = 32
  cluster_name            = "tf_example_cluster_np"
  cluster_desc            = "example for tke cluster"
  cluster_max_service_num = 32
  cluster_version         = "1.22.5"
  cluster_deploy_type     = "MANAGED_CLUSTER"
  # without any worker config
}

resource "tencentcloud_kubernetes_node_pool" "example" {
  name                     = "tf_example_node_pool"
  cluster_id               = tencentcloud_kubernetes_cluster.example.id
  max_size                 = 6 # set the node scaling range [1,6]
  min_size                 = 1
  vpc_id                   = local.first_vpc_id
  subnet_ids               = [local.first_subnet_id]
  retry_policy             = "INCREMENTAL_INTERVALS"
  desired_capacity         = 4
  enable_auto_scale        = true
  multi_zone_subnet_policy = "EQUALITY"

  auto_scaling_config {
    instance_type      = var.default_instance_type
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"
    orderly_security_group_ids = [local.sg_id]

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "test123#"
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
    extra_args = [
      "root-dir=/var/lib/kubelet"
    ]
  }
}
````

Create a cluster with a node pool and open the network access with cluster endpoint

The cluster's internet and intranet access will be opened after nodes are added through node pool.

```hcl
variable "default_instance_type" {
  default = "SA2.2XLARGE16"
}

variable "availability_zone_first" {
  default = "ap-guangzhou-3"
}

variable "availability_zone_second" {
  default = "ap-guangzhou-4"
}

variable "example_cluster_cidr" {
  default = "10.31.0.0/16"
}

locals {
  first_vpc_id    = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.vpc_id
  first_subnet_id = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.subnet_id
  sg_id    = tencentcloud_security_group.sg.id
}

data "tencentcloud_vpc_subnets" "vpc_one" {
  is_default        = true
  availability_zone = var.availability_zone_first
}

data "tencentcloud_vpc_subnets" "vpc_two" {
  is_default        = true
  availability_zone = var.availability_zone_second
}

resource "tencentcloud_security_group" "sg" {
  name = "tf-example-np-ep-sg"
}

resource "tencentcloud_security_group_lite_rule" "sg_rule" {
  security_group_id = tencentcloud_security_group.sg.id

  ingress = [
    "ACCEPT#10.0.0.0/16#ALL#ALL",
    "ACCEPT#172.16.0.0/22#ALL#ALL",
    "DROP#0.0.0.0/0#ALL#ALL",
  ]

  egress = [
    "ACCEPT#172.16.0.0/22#ALL#ALL",
  ]
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                  = local.first_vpc_id
  cluster_cidr            = var.example_cluster_cidr
  cluster_max_pod_num     = 32
  cluster_name            = "tf_example_cluster"
  cluster_desc            = "example for tke cluster"
  cluster_max_service_num = 32
  cluster_internet        = false # (can be ignored) open it after the nodes added
  cluster_version         = "1.22.5"
  cluster_deploy_type     = "MANAGED_CLUSTER"
  # without any worker config
}

resource "tencentcloud_kubernetes_node_pool" "example" {
  name                     = "tf_example_node_pool"
  cluster_id               = tencentcloud_kubernetes_cluster.example.id
  max_size                 = 6 # set the node scaling range [1,6]
  min_size                 = 1
  vpc_id                   = local.first_vpc_id
  subnet_ids               = [local.first_subnet_id]
  retry_policy             = "INCREMENTAL_INTERVALS"
  desired_capacity         = 4
  enable_auto_scale        = true
  multi_zone_subnet_policy = "EQUALITY"

  auto_scaling_config {
    instance_type      = var.default_instance_type
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"
    orderly_security_group_ids = [local.sg_id]

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "test123#"
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
    extra_args = [
      "root-dir=/var/lib/kubelet"
    ]
  }
}

resource "tencentcloud_kubernetes_cluster_endpoint" "example" {
  cluster_id                      = tencentcloud_kubernetes_cluster.example.id
  cluster_internet                = true # open the internet here
  cluster_intranet                = true
  cluster_internet_security_group = local.sg_id
  cluster_intranet_subnet_id      = local.first_subnet_id
  depends_on = [ # wait for the node pool ready
    tencentcloud_kubernetes_node_pool.example
  ]
}
````

Use Kubelet

```hcl
# Create a baisc kubernetes cluster with two nodes.
variable "default_instance_type" {
  default = "SA2.2XLARGE16"
}

variable "availability_zone_first" {
  default = "ap-guangzhou-3"
}

variable "availability_zone_second" {
  default = "ap-guangzhou-4"
}

variable "example_cluster_cidr" {
  default = "10.31.0.0/16"
}

locals {
  first_vpc_id     = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.vpc_id
  first_subnet_id  = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.subnet_id
  second_vpc_id    = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.vpc_id
  second_subnet_id = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.subnet_id
  sg_id            = tencentcloud_security_group.sg.id
  image_id         = data.tencentcloud_images.default.image_id
}

data "tencentcloud_vpc_subnets" "vpc_one" {
  is_default        = true
  availability_zone = var.availability_zone_first
}

data "tencentcloud_vpc_subnets" "vpc_two" {
  is_default        = true
  availability_zone = var.availability_zone_second
}

resource "tencentcloud_security_group" "sg" {
  name = "tf-example-sg"
}

resource "tencentcloud_security_group_lite_rule" "sg_rule" {
  security_group_id = tencentcloud_security_group.sg.id

  ingress = [
    "ACCEPT#10.0.0.0/16#ALL#ALL",
    "ACCEPT#172.16.0.0/22#ALL#ALL",
    "DROP#0.0.0.0/0#ALL#ALL",
  ]

  egress = [
    "ACCEPT#172.16.0.0/22#ALL#ALL",
  ]
}

data "tencentcloud_images" "default" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                          = local.first_vpc_id
  cluster_cidr                    = var.example_cluster_cidr
  cluster_max_pod_num             = 32
  cluster_name                    = "tf_example_cluster"
  cluster_desc                    = "example for tke cluster"
  cluster_max_service_num         = 32
  cluster_internet                = false
  cluster_internet_security_group = local.sg_id
  cluster_version                 = "1.22.5"
  cluster_deploy_type             = "MANAGED_CLUSTER"

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_first
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.first_subnet_id
    # img_id                     = local.image_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
      encrypt   = false
    }

    enhanced_security_service  = false
    enhanced_monitor_service   = false
    user_data                  = "dGVzdA=="
    disaster_recover_group_ids = []
    security_group_ids         = []
    key_ids                    = []
    password                   = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
  }

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_second
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.second_subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service  = false
    enhanced_monitor_service   = false
    user_data                  = "dGVzdA=="
    disaster_recover_group_ids = []
    security_group_ids         = []
    key_ids                    = []
    cam_role_name              = "CVM_QcsRole"
    password                   = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }

  extra_args = [
    "root-dir=/var/lib/kubelet"
  ]
}
```

Use extension addons

```hcl
variable "availability_zone_first" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "10.31.0.0/16"
}

variable "default_instance_type" {
  default = "S5.SMALL1"
}

data "tencentcloud_vpc_subnets" "vpc_first" {
  is_default        = true
  availability_zone = var.availability_zone_first
}

# fetch latest addon(chart) versions
data "tencentcloud_kubernetes_charts" "charts" {}

locals {
  chartNames = data.tencentcloud_kubernetes_charts.charts.chart_list.*.name
  chartVersions = data.tencentcloud_kubernetes_charts.charts.chart_list.*.latest_version
  chartMap = zipmap(local.chartNames, local.chartVersions)
}

resource "tencentcloud_kubernetes_cluster" "cluster_with_addon" {
  vpc_id                                     = data.tencentcloud_vpc_subnets.vpc_first.instance_list.0.vpc_id
  cluster_cidr                               = var.cluster_cidr
  cluster_max_pod_num                        = 32
  cluster_name                               = "test"
  cluster_desc                               = "test cluster desc"
  cluster_max_service_num                    = 32
  cluster_internet                           = true
  # managed_cluster_internet_security_policies = ["3.3.3.3", "1.1.1.1"]
  cluster_deploy_type                        = "MANAGED_CLUSTER"

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_first
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = data.tencentcloud_vpc_subnets.vpc_first.instance_list.0.subnet_id
    # img_id                     = "img-rkiynh11"
    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    # password                  = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
    key_ids                   = "skey-11112222"
  }

  extension_addon {
    name  = "COS"
    param = jsonencode({
      "kind" : "App", "spec" : {
        "chart" : { "chartName" : "cos", "chartVersion" : local.chartMap["cos"] },
        "values" : { "values" : [], "rawValues" : "e30=", "rawValuesType" : "json" }
      }
    })
  }
  extension_addon {
    name  = "SecurityGroupPolicy"
    param = jsonencode({
      "kind" : "App", "spec" : { "chart" : { "chartName" : "securitygrouppolicy", "chartVersion" : local.chartMap["securitygrouppolicy"] } }
    })
  }
  extension_addon {
    name  = "OOMGuard"
    param = jsonencode({
      "kind" : "App", "spec" : { "chart" : { "chartName" : "oomguard", "chartVersion" : local.chartMap["oomguard"] } }
    })
  }
  extension_addon {
    name  = "OLM"
    param = jsonencode({
      "kind" : "App", "spec" : { "chart" : { "chartName" : "olm", "chartVersion" : local.chartMap["olm"] } }
    })
  }
}
```

Use node pool global config

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "vpc" {
  default = "vpc-dk8zmwuf"
}

variable "subnet" {
  default = "subnet-pqfek0t8"
}

variable "default_instance_type" {
  default = "SA1.LARGE8"
}

resource "tencentcloud_kubernetes_cluster" "test_node_pool_global_config" {
  vpc_id                                     = var.vpc
  cluster_cidr                               = "10.1.0.0/16"
  cluster_max_pod_num                        = 32
  cluster_name                               = "test"
  cluster_desc                               = "test cluster desc"
  cluster_max_service_num                    = 32
  cluster_internet                           = true
  # managed_cluster_internet_security_policies = ["3.3.3.3", "1.1.1.1"]
  cluster_deploy_type                        = "MANAGED_CLUSTER"

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = var.subnet

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    # password                  = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
    key_ids                   = "skey-11112222"
  }

  node_pool_global_config {
    is_scale_in_enabled = true
    expander = "random"
    ignore_daemon_sets_utilization = true
    max_concurrent_scale_in = 5
    scale_in_delay = 15
    scale_in_unneeded_time = 15
    scale_in_utilization_threshold = 30
    skip_nodes_with_local_storage = false
    skip_nodes_with_system_pods = true
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
}
```

Using VPC-CNI network type

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-1"
}

variable "vpc" {
  default = "vpc-r1m1fyx5"
}

variable "default_instance_type" {
  default = "SA2.SMALL2"
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = var.vpc
  cluster_max_pod_num     = 32
  cluster_name            = "test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 256
  cluster_internet        = true
  cluster_deploy_type     = "MANAGED_CLUSTER"
  network_type            = "VPC-CNI"
  eni_subnet_ids          = ["subnet-bk1etlyu"]
  service_cidr            = "10.1.0.0/24"

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_PREMIUM"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = "subnet-t5dv27rs"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    key_ids                   = "skey-11112222"
    # password                  = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
}
```

Using ops options

```
resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  # ...your basic fields

  log_agent {
    enabled = true
    kubelet_root_dir = "" # optional
  }

  event_persistence {
    enabled = true
    log_set_id = "" # optional
    topic_id = "" # optional
  }

  cluster_audit {
    enabled = true
    log_set_id = "" # optional
    topic_id = "" # optional
  }
}
```

Create a CDC scenario cluster

```
resource "tencentcloud_kubernetes_cluster" "cdc_cluster" {
  cdc_id                  = "cluster-262n63e8"
  vpc_id                  = "vpc-0m6078eb"
  cluster_cidr            = "192.168.0.0/16"
  cluster_max_pod_num     = 64
  cluster_name            = "test-cdc"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 1024
  cluster_version         = "1.30.0"
  cluster_os              = "tlinux3.1x86_64"
  cluster_level           = "L20"
  cluster_deploy_type     = "INDEPENDENT_CLUSTER"
  container_runtime       = "containerd"
  runtime_version         = "1.6.9"
  pre_start_user_script   = "aXB0YWJsZXMgLUEgSU5QVVQgLXAgdGNwIC1zIDE2OS4yNTQuMC4wLzE5IC0tdGNwLWZsYWdzIFNZTixSU1QgU1lOIC1qIFRDUE1TUyAtLXNldC1tc3MgMTE2MAppcHRhYmxlcyAtQSBPVVRQVVQgLXAgdGNwIC1kIDE2OS4yNTQuMC4wLzE5IC0tdGNwLWZsYWdzIFNZTixSU1QgU1lOIC1qIFRDUE1TUyAtLXNldC1tc3MgMTE2MAoKZWNobyAnCmlwdGFibGVzIC1BIElOUFVUIC1wIHRjcCAtcyAxNjkuMjU0LjAuMC8xOSAtLXRjcC1mbGFncyBTWU4sUlNUIFNZTiAtaiBUQ1BNU1MgLS1zZXQtbXNzIDExNjAKaXB0YWJsZXMgLUEgT1VUUFVUIC1wIHRjcCAtZCAxNjkuMjU0LjAuMC8xOSAtLXRjcC1mbGFncyBTWU4sUlNUIFNZTiAtaiBUQ1BNU1MgLS1zZXQtbXNzIDExNjAKJyA+PiAvZXRjL3JjLmQvcmMubG9jYWw="
  instance_delete_mode    = "retain"
  exist_instance {
    node_role = "MASTER_ETCD"
    instances_para {
      instance_ids              = ["ins-mam0c7lw", "ins-quvwayve", "ins-qbffk8iw"]
      enhanced_security_service = true
      enhanced_monitor_service  = true
      password                  = "Password@123"
      security_group_ids        = ["sg-hjs685q9"]
      master_config {
        mount_target      = "/var/data"
        docker_graph_path = "/var/lib/containerd"
        unschedulable     = 0
        labels {
          name  = "key"
          value = "value"
        }
        data_disk {
          file_system           = "ext4"
          auto_format_and_mount = true
          mount_target          = "/var/data"
          disk_partition        = "/dev/vdb"
        }
        extra_args {
          kubelet = ["root-dir=/root"]
        }
        taints {
          key    = "key"
          value  = "value"
          effect = "NoSchedule"
        }
      }
    }
  }
}
```

Use delete options to delete CBS when deleting the Cluster

```hcl
resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                     = local.first_vpc_id
  cluster_cidr               = var.example_cluster_cidr
  cluster_max_pod_num        = 32
  cluster_name               = "example"
  cluster_desc               = "example for tke cluster"
  cluster_max_service_num    = 32
  cluster_level              = "L50"
  auto_upgrade_cluster_level = true
  cluster_internet           = false # (can be ignored) open it after the nodes added
  cluster_version            = "1.30.0"
  cluster_os                 = "tlinux2.2(tkernel3)x86_64"
  cluster_deploy_type        = "MANAGED_CLUSTER"
  container_runtime          = "containerd"
  docker_graph_path          = "/var/lib/containerd"
  # without any worker config
  tags = {
    "demo" = "test"
  }

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_first
    instance_type              = "SA2.MEDIUM2"
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.first_subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service  = false
    enhanced_monitor_service   = false
    user_data                  = "dGVzdA=="
    disaster_recover_group_ids = []
    security_group_ids         = []
    key_ids                    = []
    cam_role_name              = "CVM_QcsRole"
    password                   = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
  }

  resource_delete_options {
    resource_type = "CBS"
    delete_mode   = "terminate"
  }
}
```

Import

tke cluster can be imported, e.g.

```
$ terraform import tencentcloud_kubernetes_cluster.example cls-n2h4jbtk
```