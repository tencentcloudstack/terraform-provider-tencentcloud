#examples for MANAGED_CLUSTER  cluster
resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = "${var.vpc}"
  cluster_cidr            = "10.1.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32

  worker_config {
    count                      = 2
    availability_zone          = "${var.availability_zone}"
    instance_type              = "${var.default_instance_type}"
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = "${var.subnet}"

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
}

#examples for INDEPENDENT_CLUSTER  cluster
resource "tencentcloud_kubernetes_cluster" "independing_cluster" {
  vpc_id                  = "${var.vpc}"
  cluster_cidr            = "10.1.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32

  master_config {
    count                      = 3
    availability_zone          = "${var.availability_zone}"
    instance_type              = "${var.default_instance_type}"
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = "${var.subnet}"

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
    availability_zone          = "${var.availability_zone}"
    instance_type              = "${var.default_instance_type}"
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = "${var.subnet}"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "INDEPENDENT_CLUSTER"
}

#examples for scale  worker
resource tencentcloud_kubernetes_scale_worker test_scale {
  cluster_id = "${tencentcloud_kubernetes_cluster.managed_cluster.id}"

  worker_config {
    count                      = 3
    availability_zone          = "${var.availability_zone}"
    instance_type              = "${var.scale_instance_type}"
    subnet_id                  = "${var.subnet}"
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

  cluster_id = "${tencentcloud_kubernetes_cluster.managed_cluster.id}"

  auto_scaling_group {
    scaling_group_name   = "tf-guagua-as-group"
    max_size             = "5"
    min_size             = "0"
    vpc_id               = "${var.vpc}"
    subnet_ids           = ["${var.subnet}"]
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
    instance_type      = "${var.scale_instance_type}"
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
}