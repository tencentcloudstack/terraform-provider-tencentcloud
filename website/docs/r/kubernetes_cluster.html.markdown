---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster"
sidebar_current: "docs-tencentcloud-resource-kubernetes_cluster"
description: |-
  Provide a resource to create a kubernetes cluster.
---

# tencentcloud_kubernetes_cluster

Provide a resource to create a kubernetes cluster.

## Example Usage

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

#examples for MANAGED_CLUSTER cluster
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
}

#examples for INDEPENDENT_CLUSTER cluster
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
```

## Argument Reference

The following arguments are supported:

* `cluster_cidr` - (Required, ForceNew) A network address block of the cluster. Different from vpc cidr and cidr of other clusters within this vpc. Must be in  10./192.168/172.[16-31] segments.
* `vpc_id` - (Required, ForceNew) Vpc Id of the cluster.
* `cluster_deploy_type` - (Optional, ForceNew) Deployment type of the cluster, the available values include: 'MANAGED_CLUSTER' and 'INDEPENDENT_CLUSTER', Default is 'MANAGED_CLUSTER'.
* `cluster_desc` - (Optional, ForceNew) Description of the cluster.
* `cluster_ipvs` - (Optional, ForceNew) Indicates whether ipvs is enabled. Default is true.
* `cluster_max_pod_num` - (Optional, ForceNew) The maximum number of Pods per node in the cluster. Default is 256. Must be a multiple of 16 and large than 32.
* `cluster_max_service_num` - (Optional, ForceNew) The maximum number of services in the cluster. Default is 256. Must be a multiple of 16.
* `cluster_name` - (Optional, ForceNew) Name of the cluster.
* `cluster_os` - (Optional, ForceNew) Operating system of the cluster, the available values include: 'centos7.2x86_64' and 'ubuntu16.04.1 LTSx86_64'. Default is 'ubuntu16.04.1 LTSx86_64'.
* `cluster_version` - (Optional, ForceNew) Version of the cluster, Default is '1.10.5'.
* `container_runtime` - (Optional, ForceNew) Runtime type of the cluster, the available values include: 'docker' and 'containerd'. Default is 'docker'.
* `ignore_cluster_cidr_conflict` - (Optional, ForceNew) Indicates whether to ignore the cluster cidr conflict error. Default is false.
* `master_config` - (Optional, ForceNew) Deploy the machine configuration information of the 'MASTER_ETCD' service, and create <=7 units for common users.
* `project_id` - (Optional, ForceNew) Project ID, default value is 0.
* `tags` - (Optional) The tags of the cluster.
* `worker_config` - (Optional, ForceNew) Deploy the machine configuration information of the 'WORKER' service, and create <=20 units for common users. The other 'WORK' service are added by 'tencentcloud_kubernetes_worker'.

The `data_disk` object supports the following:

* `disk_size` - (Optional, ForceNew) Volume of disk in GB. Default is 0.
* `disk_type` - (Optional, ForceNew) Types of disk, available values: CLOUD_PREMIUM and CLOUD_SSD.
* `snapshot_id` - (Optional, ForceNew) Data disk snapshot ID.

The `master_config` object supports the following:

* `instance_type` - (Required, ForceNew) Specified types of CVM instance.
* `subnet_id` - (Required, ForceNew) Private network ID.
* `availability_zone` - (Optional, ForceNew) Indicates which availability zone will be used.
* `count` - (Optional, ForceNew) Number of cvm.
* `data_disk` - (Optional, ForceNew) Configurations of data disk.
* `enhanced_monitor_service` - (Optional, ForceNew) To specify whether to enable cloud monitor service. Default is TRUE.
* `enhanced_security_service` - (Optional, ForceNew) To specify whether to enable cloud security service. Default is TRUE.
* `instance_name` - (Optional, ForceNew) Name of the CVMs.
* `internet_charge_type` - (Optional, ForceNew) Charge types for network traffic. Available values include TRAFFIC_POSTPAID_BY_HOUR.
* `internet_max_bandwidth_out` - (Optional, ForceNew) Max bandwidth of Internet access in Mbps. Default is 0.
* `key_ids` - (Optional, ForceNew) ID list of keys, should be set if `password` not set.
* `password` - (Optional, ForceNew) Password to access, should be set if `key_ids` not set.
* `public_ip_assigned` - (Optional, ForceNew) Specify whether to assign an Internet IP address.
* `security_group_ids` - (Optional, ForceNew) Security groups to which a CVM instance belongs.
* `system_disk_size` - (Optional, ForceNew) Volume of system disk in GB. Default is 50.
* `system_disk_type` - (Optional, ForceNew) Type of a CVM disk, and available values include CLOUD_PREMIUM and CLOUD_SSD. Default is CLOUD_PREMIUM.
* `user_data` - (Optional, ForceNew) ase64-encoded User Data text, the length limit is 16KB.

The `worker_config` object supports the following:

* `instance_type` - (Required, ForceNew) Specified types of CVM instance.
* `subnet_id` - (Required, ForceNew) Private network ID.
* `availability_zone` - (Optional, ForceNew) Indicates which availability zone will be used.
* `count` - (Optional, ForceNew) Number of cvm.
* `data_disk` - (Optional, ForceNew) Configurations of data disk.
* `enhanced_monitor_service` - (Optional, ForceNew) To specify whether to enable cloud monitor service. Default is TRUE.
* `enhanced_security_service` - (Optional, ForceNew) To specify whether to enable cloud security service. Default is TRUE.
* `instance_name` - (Optional, ForceNew) Name of the CVMs.
* `internet_charge_type` - (Optional, ForceNew) Charge types for network traffic. Available values include TRAFFIC_POSTPAID_BY_HOUR.
* `internet_max_bandwidth_out` - (Optional, ForceNew) Max bandwidth of Internet access in Mbps. Default is 0.
* `key_ids` - (Optional, ForceNew) ID list of keys, should be set if `password` not set.
* `password` - (Optional, ForceNew) Password to access, should be set if `key_ids` not set.
* `public_ip_assigned` - (Optional, ForceNew) Specify whether to assign an Internet IP address.
* `security_group_ids` - (Optional, ForceNew) Security groups to which a CVM instance belongs.
* `system_disk_size` - (Optional, ForceNew) Volume of system disk in GB. Default is 50.
* `system_disk_type` - (Optional, ForceNew) Type of a CVM disk, and available values include CLOUD_PREMIUM and CLOUD_SSD. Default is CLOUD_PREMIUM.
* `user_data` - (Optional, ForceNew) ase64-encoded User Data text, the length limit is 16KB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `certification_authority` - The certificate used for access.
* `cluster_external_endpoint` - External network address to access.
* `cluster_node_num` - Number of nodes in the cluster.
* `domain` - Domain name for access.
* `password` - Password of account.
* `pgw_endpoint` - The Intranet address used for access.
* `security_policy` - Access policy.
* `user_name` - User name of account.
* `worker_instances_list` - An information list of cvm within the 'WORKER' clusters. Each element contains the following attributes:
  * `failed_reason` - Information of the cvm when it is failed.
  * `instance_id` - ID of the cvm.
  * `instance_role` - Role of the cvm.
  * `instance_state` - State of the cvm.


