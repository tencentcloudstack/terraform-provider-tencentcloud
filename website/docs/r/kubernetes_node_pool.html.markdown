---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_node_pool"
sidebar_current: "docs-tencentcloud-resource-kubernetes_node_pool"
description: |-
  Provide a resource to create an auto scaling group for kubernetes cluster.
---

# tencentcloud_kubernetes_node_pool

Provide a resource to create an auto scaling group for kubernetes cluster.

~> **NOTE:**  We recommend the usage of one cluster with essential worker config + node pool to manage cluster and nodes. Its a more flexible way than manage worker config with tencentcloud_kubernetes_cluster, tencentcloud_kubernetes_scale_worker or exist node management of `tencentcloud_kubernetes_attachment`. Cause some unchangeable parameters of `worker_config` may cause the whole cluster resource `force new`.

~> **NOTE:**  In order to ensure the integrity of customer data, if you destroy nodepool instance, it will keep the cvm instance associate with nodepool by default. If you want to destroy together, please set `delete_keep_instance` to `false`.

~> **NOTE:**  In order to ensure the integrity of customer data, if the cvm instance was destroyed due to shrinking, it will keep the cbs associate with cvm by default. If you want to destroy together, please set `delete_with_instance` to `true`.

~> **NOTE:**  There are two parameters `wait_node_ready` and `scale_tolerance` to ensure better management of node pool scaling operations. If this parameter is set when creating a resource, the resource will be marked as `tainted` if the set conditions are not met.

## Example Usage

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

### Using Spot CVM Instance

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

### If instance_type is CBM

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

### Wait for all scaling nodes to be ready with wait_node_ready and scale_tolerance parameters. The default maximum scaling timeout is 30 minutes.

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

### Create Node pool for CDC cluster

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

## Argument Reference

The following arguments are supported:

* `auto_scaling_config` - (Required, List) Auto scaling config parameters.
* `cluster_id` - (Required, String, ForceNew) ID of the cluster.
* `max_size` - (Required, Int) Maximum number of node.
* `min_size` - (Required, Int) Minimum number of node.
* `name` - (Required, String) Name of the node pool. The name does not exceed 25 characters, and only supports Chinese, English, numbers, underscores, separators (`-`) and decimal points.
* `vpc_id` - (Required, String, ForceNew) ID of VPC network.
* `annotations` - (Optional, Set) Node Annotation List.
* `auto_update_instance_tags` - (Optional, Bool, ForceNew) Automatically update instance tags. The default value is false. After configuration, if the scaling group tags are updated, the tags of the running instances in the scaling group will be updated synchronously (synchronous updates only support adding and modifying tags, and do not support deleting tags for the time being). Synchronous updates do not take effect immediately and there is a certain delay.
* `default_cooldown` - (Optional, Int) Seconds of scaling group cool down. Default value is `300`.
* `delete_keep_instance` - (Optional, Bool) Indicate to keep the CVM instance when delete the node pool. Default is `true`.
* `deletion_protection` - (Optional, Bool) Indicates whether the node pool deletion protection is enabled.
* `desired_capacity` - (Optional, Int) Desired capacity of the node. If `enable_auto_scale` is set `true`, this will be a computed parameter.
* `enable_auto_scale` - (Optional, Bool) Indicate whether to enable auto scaling or not.
* `labels` - (Optional, Map) Labels of kubernetes node pool created nodes. The label key name does not exceed 63 characters, only supports English, numbers,'/','-', and does not allow beginning with ('/').
* `multi_zone_subnet_policy` - (Optional, String) Multi-availability zone/subnet policy. Valid values: PRIORITY and EQUALITY. Default value: PRIORITY.
* `node_config` - (Optional, List) Node config.
* `node_os_type` - (Optional, String) The image version of the node. Valida values are `DOCKER_CUSTOMIZE` and `GENERAL`. Default is `GENERAL`. This parameter will only affect new nodes, not including the existing nodes.
* `node_os` - (Optional, String) Node pool operating system (enter the image ID for a custom image, and enter the OS name for a public image). If custom image, please refer to [TencentCloud Documentation](https://www.tencentcloud.com/document/product/457/46750?lang=en&pg=#list-of-public-images-supported-by-tke) for available values. Default is 'tlinux2.4x86_64'. This parameter will only affect new nodes, not including the existing nodes.
* `retry_policy` - (Optional, String, ForceNew) Available values for retry policies include `IMMEDIATE_RETRY` and `INCREMENTAL_INTERVALS`.
* `scale_tolerance` - (Optional, Int) Control how many expectations(`desired_capacity`) can be tolerated successfully. Unit is percentage, Default is `100`. Only can be set if `wait_node_ready` is `true`.
* `scaling_group_name` - (Optional, String) Name of relative scaling group.
* `scaling_group_project_id` - (Optional, Int) Project ID the scaling group belongs to.
* `scaling_mode` - (Optional, String, ForceNew) Auto scaling mode. Valid values are `CLASSIC_SCALING`(scaling by create/destroy instances), `WAKE_UP_STOPPED_SCALING`(Boot priority for expansion. When expanding the capacity, the shutdown operation is given priority to the shutdown of the instance. If the number of instances is still lower than the expected number of instances after the startup, the instance will be created, and the method of destroying the instance will still be used for shrinking).
* `subnet_ids` - (Optional, List: [`String`], ForceNew) ID list of subnet, and for VPC it is required.
* `tags` - (Optional, Map) Node pool tag specifications, will passthroughs to the scaling instances.
* `taints` - (Optional, List) Taints of kubernetes node pool created nodes.
* `termination_policies` - (Optional, List: [`String`]) Policy of scaling group termination. Available values: `["OLDEST_INSTANCE"]`, `["NEWEST_INSTANCE"]`.
* `unschedulable` - (Optional, Int, ForceNew) Sets whether the joining node participates in the schedule. Default is '0'. Participate in scheduling.
* `wait_node_ready` - (Optional, Bool) Whether to wait for all desired nodes to be ready. Default is false. Only can be set if `enable_auto_scale` is `false`.
* `zones` - (Optional, List: [`String`]) List of auto scaling group available zones, for Basic network it is required.

The `annotations` object supports the following:

* `name` - (Required, String) Name in the map table.
* `value` - (Required, String) Value in the map table.

The `auto_scaling_config` object supports the following:

* `instance_type` - (Required, String, ForceNew) Specified types of CVM instance.
* `backup_instance_types` - (Optional, List) Backup CVM instance types if specified instance type sold out or mismatch.
* `bandwidth_package_id` - (Optional, String) bandwidth package id. if user is standard user, then the bandwidth_package_id is needed, or default has bandwidth_package_id.
* `cam_role_name` - (Optional, String, ForceNew) Name of cam role.
* `cdc_id` - (Optional, String, ForceNew) CDC ID.
* `data_disk` - (Optional, List) Configurations of data disk.
* `enhanced_monitor_service` - (Optional, Bool, ForceNew) To specify whether to enable cloud monitor service. Default is TRUE.
* `enhanced_security_service` - (Optional, Bool) To specify whether to enable cloud security service. Default is TRUE.
* `host_name_style` - (Optional, String) The style of the host name of the cloud server, the value range includes ORIGINAL and UNIQUE, and the default is ORIGINAL. For usage, refer to `HostNameSettings` in https://www.tencentcloud.com/document/product/377/31001.
* `host_name` - (Optional, String) The hostname of the cloud server, dot (.) and dash (-) cannot be used as the first and last characters of HostName and cannot be used consecutively. Windows instances are not supported. Examples of other types (Linux, etc.): The character length is [2, 40], multiple periods are allowed, and there is a paragraph between the dots, and each paragraph is allowed to consist of letters (unlimited case), numbers and dashes (-). Pure numbers are not allowed. For usage, refer to `HostNameSettings` in https://www.tencentcloud.com/document/product/377/31001.
* `instance_charge_type_prepaid_period` - (Optional, Int) The tenancy (in month) of the prepaid instance, NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.
* `instance_charge_type_prepaid_renew_flag` - (Optional, String) Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.
* `instance_charge_type` - (Optional, String) Charge type of instance. Valid values are `PREPAID`, `POSTPAID_BY_HOUR`, `SPOTPAID`, `CDCPAID`. The default is `POSTPAID_BY_HOUR`. NOTE: `SPOTPAID` instance must set `spot_instance_type` and `spot_max_price` at the same time.
* `instance_name_style` - (Optional, String) Type of CVM instance name. Valid values: `ORIGINAL` and `UNIQUE`. Default value: `ORIGINAL`. For usage, refer to `InstanceNameSettings` in https://www.tencentcloud.com/document/product/377/31001.
* `instance_name` - (Optional, String) Instance name, no more than 60 characters. For usage, refer to `InstanceNameSettings` in https://www.tencentcloud.com/document/product/377/31001.
* `internet_charge_type` - (Optional, String) Charge types for network traffic. Valid value: `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.
* `internet_max_bandwidth_out` - (Optional, Int) Max bandwidth of Internet access in Mbps. Default is `0`.
* `ipv4_address_type` - (Optional, String) Type of public IP address. WanIP: Ordinary public IP address; HighQualityEIP: High Quality EIP is supported only in Singapore and Hong Kong; AntiDDoSEIP: Anti-DDoS IP is supported only in specific regions. For details, see EIP Product Overview. Specify the type of public IPv4 address to assign a public IPv4 address to the resource. HighQualityEIP and AntiDDoSEIP features are gradually released in select regions. For usage, submit a ticket for consultation.
* `key_ids` - (Optional, List, ForceNew) ID list of keys.
* `orderly_security_group_ids` - (Optional, List) Ordered security groups to which a CVM instance belongs.
* `password` - (Optional, String, ForceNew) Password to access.
* `public_ip_assigned` - (Optional, Bool) Specify whether to assign an Internet IP address.
* `security_group_ids` - (Optional, Set, **Deprecated**) The order of elements in this field cannot be guaranteed. Use `orderly_security_group_ids` instead. Security groups to which a CVM instance belongs.
* `spot_instance_type` - (Optional, String) Type of spot instance, only support `one-time` now. Note: it only works when instance_charge_type is set to `SPOTPAID`.
* `spot_max_price` - (Optional, String) Max price of a spot instance, is the format of decimal string, for example "0.50". Note: it only works when instance_charge_type is set to `SPOTPAID`.
* `system_disk_size` - (Optional, Int) Volume of system disk in GB. Default is `50`.
* `system_disk_type` - (Optional, String) Type of a CVM disk. Valid value: `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_BASIC`, `CLOUD_PREMIUM`, `CLOUD_SSD`, `CLOUD_HSSD`, `CLOUD_TSSD`, `CLOUD_BSSD` and `LOCAL_NVME`. Default is `CLOUD_PREMIUM`.

The `data_disk` object of `auto_scaling_config` supports the following:

* `delete_with_instance` - (Optional, Bool) Indicates whether the disk remove after instance terminated. Default is `false`.
* `disk_size` - (Optional, Int) Volume of disk in GB. Default is `0`.
* `disk_type` - (Optional, String) Types of disk. Valid value: `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_BASIC`, `CLOUD_PREMIUM`, `CLOUD_SSD`, `CLOUD_HSSD`, `CLOUD_TSSD` and `CLOUD_BSSD`.
* `encrypt` - (Optional, Bool) Specify whether to encrypt data disk, default: false. NOTE: Make sure the instance type is offering and the cam role `QcloudKMSAccessForCVMRole` was provided.
* `snapshot_id` - (Optional, String, ForceNew) Data disk snapshot ID.
* `throughput_performance` - (Optional, Int) Add extra performance to the data disk. Only works when disk type is `CLOUD_TSSD` or `CLOUD_HSSD` and `data_size` > 460GB.

The `data_disk` object of `node_config` supports the following:

* `auto_format_and_mount` - (Optional, Bool, ForceNew) Indicate whether to auto format and mount or not. Default is `false`.
* `disk_partition` - (Optional, String, ForceNew) The name of the device or partition to mount. NOTE: this argument doesn't support setting in node pool, or will leads to mount error.
* `disk_size` - (Optional, Int, ForceNew) Volume of disk in GB. Default is `0`.
* `disk_type` - (Optional, String, ForceNew) Types of disk. Valid value: `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_BASIC`, `CLOUD_PREMIUM`, `CLOUD_SSD`, `CLOUD_HSSD`, `CLOUD_TSSD`, `CLOUD_BSSD` and `LOCAL_NVME`.
* `file_system` - (Optional, String, ForceNew) File system, e.g. `ext3/ext4/xfs`.
* `mount_target` - (Optional, String, ForceNew) Mount target.

The `gpu_args` object of `node_config` supports the following:

* `cuda` - (Optional, Map) CUDA  version. Format like: `{ version: String, name: String }`. `version`: Version of GPU driver or CUDA; `name`: Name of GPU driver or CUDA.
* `cudnn` - (Optional, Map) cuDNN version. Format like: `{ version: String, name: String, doc_name: String, dev_name: String }`. `version`: cuDNN version; `name`: cuDNN name; `doc_name`: Doc name of cuDNN; `dev_name`: Dev name of cuDNN.
* `custom_driver` - (Optional, Map) Custom GPU driver. Format like: `{address: String}`. `address`: URL of custom GPU driver address.
* `driver` - (Optional, Map) GPU driver version. Format like: `{ version: String, name: String }`. `version`: Version of GPU driver or CUDA; `name`: Name of GPU driver or CUDA.
* `mig_enable` - (Optional, Bool) Whether to enable MIG.

The `node_config` object supports the following:

* `data_disk` - (Optional, List, ForceNew) Configurations of data disk.
* `desired_pod_num` - (Optional, Int, ForceNew) Indicate to set desired pod number in node. valid when the cluster is podCIDR.
* `docker_graph_path` - (Optional, String, ForceNew) Docker graph path. Default is `/var/lib/docker`.
* `extra_args` - (Optional, List, ForceNew) Custom parameter information related to the node. This is a white-list parameter.
* `gpu_args` - (Optional, List, ForceNew) GPU driver parameters.
* `is_schedule` - (Optional, Bool, ForceNew) Indicate to schedule the adding node or not. Default is true.
* `mount_target` - (Optional, String, ForceNew) Mount target. Default is not mounting.
* `pre_start_user_script` - (Optional, String) Base64-encoded user script, executed before initializing the node, currently only effective for adding existing nodes.
* `user_data` - (Optional, String) Base64-encoded User Data text, the length limit is 16KB.

The `taints` object supports the following:

* `effect` - (Required, String) Effect of the taint. Valid values are: `NoSchedule`, `PreferNoSchedule`, `NoExecute`.
* `key` - (Required, String) Key of the taint. The taint key name does not exceed 63 characters, only supports English, numbers,'/','-', and does not allow beginning with ('/').
* `value` - (Required, String) Value of the taint.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `auto_scaling_group_id` - The auto scaling group ID.
* `autoscaling_added_total` - The total of autoscaling added node.
* `launch_config_id` - The launch config ID.
* `manually_added_total` - The total of manually added node.
* `node_count` - The total node count.
* `status` - Status of the node pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `30m`) Used when creating the resource.
* `update` - (Defaults to `30m`) Used when updating the resource.

## Import

tke node pool can be imported, e.g.

```
$ terraform import tencentcloud_kubernetes_node_pool.example cls-d2xdg3io#np-380ay1o8
```

