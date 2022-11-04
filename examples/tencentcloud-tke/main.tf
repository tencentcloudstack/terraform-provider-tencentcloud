#examples for MANAGED_CLUSTER  cluster
resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = var.vpc
  cluster_cidr            = "10.1.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32

  worker_config {
    count                      = 2
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
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"

  tags = {
    "test" = "test"
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
}

#examples for MANAGED_CLUSTER  cluster with add-on
resource "tencentcloud_kubernetes_cluster" "cluster_with_addon" {
  vpc_id                                     = var.vpc
  cluster_cidr                               = "10.1.0.0/16"
  cluster_max_pod_num                        = 32
  cluster_name                               = "test"
  cluster_desc                               = "test cluster desc"
  cluster_max_service_num                    = 32
  cluster_internet                           = true
  managed_cluster_internet_security_policies = ["3.3.3.3", "1.1.1.1"]
  cluster_deploy_type                        = "MANAGED_CLUSTER"

  worker_config {
    count                      = 1
    availability_zone          = "ap-guangzhou-3"
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = var.subnet
    img_id                     = "img-rkiynh11"
    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  extension_addon {
    name  = "CBS",
    param = jsonencode({
      "kind" : "App", "spec" : {
        "chart" : { "chartName" : "cbs", "chartVersion" : "1.0.7" },
        "values" : { "values" : [], "rawValues" : "e30=", "rawValuesType" : "json" }
      }
    })
  }
  extension_addon {
    name  = "SecurityGroupPolicy",
    param = jsonencode({
      "kind" : "App", "spec" : { "chart" : { "chartName" : "securitygrouppolicy", "chartVersion" : "0.1.0" } }
    })
  }
  extension_addon {
    name  = "OOMGuard",
    param = jsonencode({
      "kind" : "App", "spec" : { "chart" : { "chartName" : "oomguard", "chartVersion" : "1.0.1" } }
    })
  }
  extension_addon {
    name  = "OLM",
    param = jsonencode({
      "kind" : "App", "spec" : { "chart" : { "chartName" : "olm", "chartVersion" : "1.0.0" } }
    })
  }
}

#examples for MANAGED_CLUSTER VPC-CNI network type cluster with customized master params
resource "tencentcloud_kubernetes_cluster" "managed_vpc_cni_cluster" {
  cluster_version         = "1.14.3"
  vpc_id                  = var.vpc
  cluster_max_pod_num     = 32
  cluster_name            = "testvpccni"
  cluster_desc            = "test vpc-cni cluster desc"
  cluster_max_service_num = 32
  service_cidr            = "192.168.128.0/24"
  eni_subnet_ids          = ["subnet-hmmlszs7", "subnet-4o0v4e7j"]
  claim_expired_seconds   = 300
  network_type            = "VPC-CNI"
  is_non_static_ip_mode   = true
  cluster_extra_args {
    kube_apiserver          = ["max-requests-inflight=450"]
    kube_controller_manager = ["kube-api-burst=500", "kube-api-qps=200"]
    kube_scheduler          = ["kube-api-burst=500", "kube-api-qps=200"]
  }

  worker_config {
    count                      = 2
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
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"

  tags = {
    "test" = "test"
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
}

#examples for INDEPENDENT_CLUSTER  cluster
resource "tencentcloud_kubernetes_cluster" "independing_cluster" {
  vpc_id                  = var.vpc
  cluster_cidr            = "10.1.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32

  master_config {
    count                      = 3
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
    password                  = "MMMZZXXccvv1212"
  }

  worker_config {
    count                      = 2
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
    password                  = "ZZXXccvv1212"
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }

  cluster_deploy_type = "INDEPENDENT_CLUSTER"
}

#examples for scale  worker
resource tencentcloud_kubernetes_scale_worker test_scale {
  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id

  worker_config {
    count                      = 3
    availability_zone          = var.availability_zone
    instance_type              = var.scale_instance_type
    subnet_id                  = var.subnet
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 50
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "AABBccdd1122"
  }
}


#examples for auto scaling group
resource "tencentcloud_kubernetes_as_scaling_group" "test" {

  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id

  auto_scaling_group {
    scaling_group_name   = "tf-guagua-as-group"
    max_size             = "5"
    min_size             = "0"
    vpc_id               = var.vpc
    subnet_ids           = [var.subnet]
    project_id           = 0
    default_cooldown     = 400
    desired_capacity     = "0"
    termination_policies = ["NEWEST_INSTANCE"]
    retry_policy         = "INCREMENTAL_INTERVALS"

    tags = {
      "test" = "test"
    }

  }


  auto_scaling_config {
    configuration_name = "tf-guagua-as-config"
    instance_type      = var.scale_instance_type
    project_id         = 0
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"

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

    instance_tags = {
      tag = "as"
    }

  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
}

#example for node pool global config
resource "tencentcloud_kubernetes_cluster" "test_node_pool_global_config" {
  vpc_id                                     = var.vpc
  cluster_cidr                               = "10.1.0.0/16"
  cluster_max_pod_num                        = 32
  cluster_name                               = "test"
  cluster_desc                               = "test cluster desc"
  cluster_max_service_num                    = 32
  cluster_internet                           = true
  managed_cluster_internet_security_policies = ["3.3.3.3", "1.1.1.1"]
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
    password                  = "ZZXXccvv1212"
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