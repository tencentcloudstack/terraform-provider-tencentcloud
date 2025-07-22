Provide a resource to increase instance to cluster

~> **NOTE:** To use the custom Kubernetes component startup parameter function (parameter `extra_args`), you need to submit a ticket for application.

~> **NOTE:** Import Node: Currently, only one node can be imported at a time.

~> **NOTE:** If you need to view error messages during instance creation, you can use parameter `create_result_output_file` to set the file save path

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

resource "tencentcloud_kubernetes_scale_worker" "example" {
  cluster_id      = "cls-godovr32"
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
    password                  = "Password@123"

    tags {
      key   = "createBy"
      value = "Terraform"
    }
  }

  create_result_output_file = "my_output_file_path"
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

resource "tencentcloud_kubernetes_scale_worker" "example" {
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
    password                  = "Password@123"
  }
}
```

Create scale worker for CDC cluster

```hcl
resource "tencentcloud_kubernetes_scale_worker" "example" {
  cluster_id = "cls-0o0dpx1a"

  worker_config {
    count                = 2
    instance_charge_type = "CDCPAID"
    instance_name        = "tke_worker_demo"
    availability_zone    = "ap-guangzhou-4"
    instance_type        = "S5.MEDIUM2"
    subnet_id            = "subnet-muu9a0gk"
    system_disk_type     = "CLOUD_SSD"
    system_disk_size     = 50
    internet_charge_type = "TRAFFIC_POSTPAID_BY_HOUR"
    security_group_ids   = ["sg-4z20n68d"]

    data_disk {
      disk_type = "CLOUD_SSD"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "Password@123"
    cdc_id                    = "cluster-262n63e8"
  }
}
```

Import

tke scale worker can be imported, e.g.

```
$ terraform import tencentcloud_kubernetes_scale_worker.example cls-mij6c2pq#ins-n6esjkdi
```