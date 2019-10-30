---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_scale_worker"
sidebar_current: "docs-tencentcloud-resource-kubernetes_scale_worker"
description: |-
  Provide a resource to increase instance to cluster
---

# tencentcloud_kubernetes_scale_worker

Provide a resource to increase instance to cluster

## Example Usage

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
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) ID of the cluster.
* `worker_config` - (Required, ForceNew) Deploy the machine configuration information of the 'WORK' service, and create <=20 units for common users.

The `data_disk` object supports the following:

* `disk_size` - (Optional, ForceNew) Volume of disk in GB. Default is 0.
* `disk_type` - (Optional, ForceNew) Types of disk, available values: CLOUD_PREMIUM and CLOUD_SSD.
* `snapshot_id` - (Optional, ForceNew) Data disk snapshot ID.

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
* `key_ids` - (Optional, ForceNew) ID list of keys.
* `password` - (Optional, ForceNew) Password to access.
* `public_ip_assigned` - (Optional, ForceNew) Specify whether to assign an Internet IP address.
* `security_group_ids` - (Optional, ForceNew) Security groups to which a CVM instance belongs.
* `system_disk_size` - (Optional, ForceNew) Volume of system disk in GB. Default is 50.
* `system_disk_type` - (Optional, ForceNew) Type of a CVM disk, and available values include CLOUD_PREMIUM and CLOUD_SSD. Default is CLOUD_PREMIUM.
* `user_data` - (Optional, ForceNew) ase64-encoded User Data text, the length limit is 16KB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `worker_instances_list` - An information list of kubernetes cluster 'WORKER'. Each element contains the following attributes:
  * `failed_reason` - Information of the cvm when it is failed.
  * `instance_id` - ID of the cvm.
  * `instance_role` - Role of the cvm.
  * `instance_state` - State of the cvm.


