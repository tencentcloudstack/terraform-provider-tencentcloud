---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_scale_worker"
sidebar_current: "docs-tencentcloud-resource-kubernetes_scale_worker"
description: |-
  Provide a resource to increase instance to cluster
---

# tencentcloud_kubernetes_scale_worker

Provide a resource to increase instance to cluster

~> **NOTE:** To use the custom Kubernetes component startup parameter function (parameter `extra_args`), you need to submit a ticket for application.

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

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) ID of the cluster.
* `worker_config` - (Required, ForceNew) Deploy the machine configuration information of the 'WORK' service, and create <=20 units for common users.
* `data_disk` - (Optional, ForceNew) Configurations of data disk.
* `desired_pod_num` - (Optional, ForceNew) Indicate to set desired pod number in current node. Valid when the cluster enable customized pod cidr.
* `docker_graph_path` - (Optional, ForceNew) Docker graph path. Default is `/var/lib/docker`.
* `extra_args` - (Optional, ForceNew) Custom parameter information related to the node.
* `labels` - (Optional, ForceNew) Labels of kubernetes scale worker created nodes.
* `mount_target` - (Optional, ForceNew) Mount target. Default is not mounting.
* `unschedulable` - (Optional, ForceNew) Sets whether the joining node participates in the schedule. Default is '0'. Participate in scheduling.

The `data_disk` object supports the following:

* `auto_format_and_mount` - (Optional, ForceNew) Indicate whether to auto format and mount or not. Default is `false`.
* `disk_size` - (Optional, ForceNew) Volume of disk in GB. Default is `0`.
* `disk_type` - (Optional, ForceNew) Types of disk, available values: `CLOUD_PREMIUM` and `CLOUD_SSD`.
* `file_system` - (Optional, ForceNew) File system, e.g. `ext3/ext4/xfs`.
* `mount_target` - (Optional, ForceNew) Mount target.

The `data_disk` object supports the following:

* `disk_size` - (Optional, ForceNew) Volume of disk in GB. Default is `0`.
* `disk_type` - (Optional, ForceNew) Types of disk, available values: `CLOUD_PREMIUM` and `CLOUD_SSD`.
* `snapshot_id` - (Optional, ForceNew) Data disk snapshot ID.

The `worker_config` object supports the following:

* `instance_type` - (Required, ForceNew) Specified types of CVM instance.
* `subnet_id` - (Required, ForceNew) Private network ID.
* `availability_zone` - (Optional, ForceNew) Indicates which availability zone will be used.
* `cam_role_name` - (Optional, ForceNew) CAM role name authorized to access.
* `count` - (Optional, ForceNew) Number of cvm.
* `data_disk` - (Optional, ForceNew) Configurations of data disk.
* `desired_pod_num` - (Optional, ForceNew) Indicate to set desired pod number in node. valid when enable_customized_pod_cidr=true, and it override `[globe_]desired_pod_num` for current node. Either all the fields `desired_pod_num` or none.
* `disaster_recover_group_ids` - (Optional, ForceNew) Disaster recover groups to which a CVM instance belongs. Only support maximum 1.
* `enhanced_monitor_service` - (Optional, ForceNew) To specify whether to enable cloud monitor service. Default is TRUE.
* `enhanced_security_service` - (Optional, ForceNew) To specify whether to enable cloud security service. Default is TRUE.
* `hostname` - (Optional, ForceNew) The host name of the attached instance. Dot (.) and dash (-) cannot be used as the first and last characters of HostName and cannot be used consecutively. Windows example: The length of the name character is [2, 15], letters (capitalization is not restricted), numbers and dashes (-) are allowed, dots (.) are not supported, and not all numbers are allowed. Examples of other types (Linux, etc.): The character length is [2, 60], and multiple dots are allowed. There is a segment between the dots. Each segment allows letters (with no limitation on capitalization), numbers and dashes (-).
* `img_id` - (Optional) The valid image id, format of img-xxx.
* `instance_charge_type_prepaid_period` - (Optional, ForceNew) The tenancy (time unit is month) of the prepaid instance. NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.
* `instance_charge_type_prepaid_renew_flag` - (Optional, ForceNew) Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.
* `instance_charge_type` - (Optional, ForceNew) The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. The default is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`, `PREPAID` instance will not terminated after cluster deleted, and may not allow to delete before expired.
* `instance_name` - (Optional, ForceNew) Name of the CVMs.
* `internet_charge_type` - (Optional, ForceNew) Charge types for network traffic. Available values include `TRAFFIC_POSTPAID_BY_HOUR`.
* `internet_max_bandwidth_out` - (Optional, ForceNew) Max bandwidth of Internet access in Mbps. Default is 0.
* `key_ids` - (Optional, ForceNew) ID list of keys, should be set if `password` not set.
* `password` - (Optional, ForceNew) Password to access, should be set if `key_ids` not set.
* `public_ip_assigned` - (Optional, ForceNew) Specify whether to assign an Internet IP address.
* `security_group_ids` - (Optional, ForceNew) Security groups to which a CVM instance belongs.
* `system_disk_size` - (Optional, ForceNew) Volume of system disk in GB. Default is `50`.
* `system_disk_type` - (Optional, ForceNew) System disk type. For more information on limits of system disk types, see [Storage Overview](https://intl.cloud.tencent.com/document/product/213/4952). Valid values: `LOCAL_BASIC`: local disk, `LOCAL_SSD`: local SSD disk, `CLOUD_BASIC`: HDD cloud disk, `CLOUD_SSD`: SSD, `CLOUD_PREMIUM`: Premium Cloud Storage. NOTE: `LOCAL_BASIC` and `LOCAL_SSD` are deprecated.
* `user_data` - (Optional, ForceNew) ase64-encoded User Data text, the length limit is 16KB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `worker_instances_list` - An information list of kubernetes cluster 'WORKER'. Each element contains the following attributes:
  * `failed_reason` - Information of the cvm when it is failed.
  * `instance_id` - ID of the cvm.
  * `instance_role` - Role of the cvm.
  * `instance_state` - State of the cvm.
  * `lan_ip` - LAN IP of the cvm.


