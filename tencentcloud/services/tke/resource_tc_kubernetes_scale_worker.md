Provide a resource to increase instance to cluster

~> **NOTE:** To use the custom Kubernetes component startup parameter function (parameter `extra_args`), you need to submit a ticket for application.

~> **NOTE:** Import Node: Currently, only one node can be imported at a time.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "subnet" {
  default = "subnet-pqfek0t8"
}

variable "scale_instance_type" {
  default = "S2.LARGE16"
}

resource tencentcloud_kubernetes_scale_worker test_scale {
  cluster_id = "cls-godovr32"
  desired_pod_num = 16
  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
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
```

Use Kubelet

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "subnet" {
  default = "subnet-pqfek0t8"
}

variable "scale_instance_type" {
  default = "S2.LARGE16"
}

resource tencentcloud_kubernetes_scale_worker test_scale {
  cluster_id = "cls-godovr32"

  extra_args = [
 	"root-dir=/var/lib/kubelet"
  ]

   labels = {
    "test1" = "test1",
    "test2" = "test2",
  }

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
```