# examples for node pool based on a empty cluster
locals {
  first_vpc_id    = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.vpc_id
  first_subnet_id = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.subnet_id
  second_vpc_id    = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.vpc_id
  second_subnet_id = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.subnet_id
  sg_id           = tencentcloud_security_group.sg.id
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
    security_group_ids = [local.sg_id]

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

# examples for node pool based on a empty cluster, and open the network through endpoint
resource "tencentcloud_kubernetes_cluster" "example_np_ep" {
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

resource "tencentcloud_kubernetes_node_pool" "example_np_ep" {
  name                     = "tf_example_node_pool_ep"
  cluster_id               = tencentcloud_kubernetes_cluster.example_np_ep.id
  max_size                 = 6 # set the node scaling range [1,6]
  min_size                 = 1
  vpc_id                   = local.second_vpc_id
  subnet_ids               = [local.second_subnet_id]
  retry_policy             = "INCREMENTAL_INTERVALS"
  desired_capacity         = 4
  enable_auto_scale        = true
  multi_zone_subnet_policy = "EQUALITY"

  auto_scaling_config {
    instance_type      = var.default_instance_type
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"
    security_group_ids = [local.sg_id]

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

resource "tencentcloud_kubernetes_cluster_endpoint" "example_np_ep" {
  cluster_id                      = tencentcloud_kubernetes_cluster.example_np_ep.id
  cluster_internet                = true # open the internet here
  cluster_intranet                = true
  cluster_internet_security_group = local.sg_id
  cluster_intranet_subnet_id      = local.second_subnet_id
  depends_on = [ # wait for the node pool ready
    tencentcloud_kubernetes_node_pool.example_np_ep
  ]
}

