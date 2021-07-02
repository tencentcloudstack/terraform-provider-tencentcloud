---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster"
sidebar_current: "docs-tencentcloud-resource-kubernetes_cluster"
description: |-
  Provide a resource to create a kubernetes cluster.
---

# tencentcloud_kubernetes_cluster

Provide a resource to create a kubernetes cluster.

~> **NOTE:** To use the custom Kubernetes component startup parameter function (parameter `extra_args`), you need to submit a ticket for application.
~> **NOTE:**  We recommend the usage of one cluster without worker config + node pool to manage cluster and nodes. It's a more flexible way than manage worker config with tencentcloud_kubernetes_cluster, tencentcloud_kubernetes_scale_worker or exist node management of `tencentcloud_kubernetes_attachment`. Cause some unchangeable parameters of `worker_config` may cause the whole cluster resource `force new`.

## Example Usage

```hcl
variable "availability_zone_first" {
  default = "ap-guangzhou-3"
}

variable "availability_zone_second" {
  default = "ap-guangzhou-4"
}

variable "cluster_cidr" {
  default = "10.31.0.0/16"
}

variable "default_instance_type" {
  default = "SA2.2XLARGE16"
}

data "tencentcloud_vpc_subnets" "vpc_first" {
  is_default        = true
  availability_zone = var.availability_zone_first
}

data "tencentcloud_vpc_subnets" "vpc_second" {
  is_default        = true
  availability_zone = var.availability_zone_second
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                                     = data.tencentcloud_vpc_subnets.vpc_first.instance_list.0.vpc_id
  cluster_cidr                               = var.cluster_cidr
  cluster_max_pod_num                        = 32
  cluster_name                               = "test"
  cluster_desc                               = "test cluster desc"
  cluster_max_service_num                    = 32
  cluster_internet                           = true
  managed_cluster_internet_security_policies = ["3.3.3.3", "1.1.1.1"]
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
    img_id                     = "img-rkiynh11"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
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
    subnet_id                  = data.tencentcloud_vpc_subnets.vpc_second.instance_list.0.subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
    cam_role_name             = "CVM_QcsRole"
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
}
```

Use Kubelet

```hcl
variable "availability_zone_first" {
  default = "ap-guangzhou-3"
}

variable "availability_zone_second" {
  default = "ap-guangzhou-4"
}

variable "cluster_cidr" {
  default = "10.31.0.0/16"
}

variable "default_instance_type" {
  default = "SA2.2XLARGE16"
}

data "tencentcloud_vpc_subnets" "vpc_first" {
  is_default        = true
  availability_zone = var.availability_zone_first
}

data "tencentcloud_vpc_subnets" "vpc_second" {
  is_default        = true
  availability_zone = var.availability_zone_second
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                                     = data.tencentcloud_vpc_subnets.vpc_first.instance_list.0.vpc_id
  cluster_cidr                               = var.cluster_cidr
  cluster_max_pod_num                        = 32
  cluster_name                               = "test"
  cluster_desc                               = "test cluster desc"
  cluster_max_service_num                    = 32
  cluster_internet                           = true
  managed_cluster_internet_security_policies = ["3.3.3.3", "1.1.1.1"]
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

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
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
    subnet_id                  = data.tencentcloud_vpc_subnets.vpc_second.instance_list.0.subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
    cam_role_name             = "CVM_QcsRole"
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
    is_scale_in_enabled            = true
    expander                       = "random"
    ignore_daemon_sets_utilization = true
    max_concurrent_scale_in        = 5
    scale_in_delay                 = 15
    scale_in_unneeded_time         = 15
    scale_in_utilization_threshold = 30
    skip_nodes_with_local_storage  = false
    skip_nodes_with_system_pods    = true
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) Vpc Id of the cluster.
* `base_pod_num` - (Optional, ForceNew) The number of basic pods. valid when enable_customized_pod_cidr=true.
* `claim_expired_seconds` - (Optional) Claim expired seconds to recycle ENI. This field can only set when field `network_type` is 'VPC-CNI'. `claim_expired_seconds` must greater or equal than 300 and less than 15768000.
* `cluster_as_enabled` - (Optional, ForceNew) Indicates whether to enable cluster node auto scaler. Default is false.
* `cluster_cidr` - (Optional, ForceNew) A network address block of the cluster. Different from vpc cidr and cidr of other clusters within this vpc. Must be in  10./192.168/172.[16-31] segments.
* `cluster_deploy_type` - (Optional, ForceNew) Deployment type of the cluster, the available values include: 'MANAGED_CLUSTER' and 'INDEPENDENT_CLUSTER'. Default is 'MANAGED_CLUSTER'.
* `cluster_desc` - (Optional) Description of the cluster.
* `cluster_extra_args` - (Optional, ForceNew) Customized parameters for master component,such as kube-apiserver, kube-controller-manager, kube-scheduler.
* `cluster_internet` - (Optional) Open internet access or not.
* `cluster_intranet_subnet_id` - (Optional) Subnet id who can access this independent cluster, this field must and can only set  when `cluster_intranet` is true. `cluster_intranet_subnet_id` can not modify once be set.
* `cluster_intranet` - (Optional) Open intranet access or not.
* `cluster_ipvs` - (Optional, ForceNew) Indicates whether `ipvs` is enabled. Default is true. False means `iptables` is enabled.
* `cluster_max_pod_num` - (Optional, ForceNew) The maximum number of Pods per node in the cluster. Default is 256. Must be a multiple of 16 and large than 32.
* `cluster_max_service_num` - (Optional, ForceNew) The maximum number of services in the cluster. Default is 256. Must be a multiple of 16.
* `cluster_name` - (Optional) Name of the cluster.
* `cluster_os_type` - (Optional, ForceNew) Image type of the cluster os, the available values include: 'GENERAL'. Default is 'GENERAL'.
* `cluster_os` - (Optional, ForceNew) Operating system of the cluster, the available values include: 'centos7.6.0_x64','ubuntu18.04.1x86_64','tlinux2.4x86_64'. Default is 'tlinux2.4x86_64'.
* `cluster_version` - (Optional) Version of the cluster, Default is '1.10.5'.
* `container_runtime` - (Optional, ForceNew) Runtime type of the cluster, the available values include: 'docker' and 'containerd'. Default is 'docker'.
* `deletion_protection` - (Optional) Indicates whether cluster deletion protection is enabled. Default is false.
* `docker_graph_path` - (Optional, ForceNew) Docker graph path. Default is `/var/lib/docker`.
* `enable_customized_pod_cidr` - (Optional) Whether to enable the custom mode of node podCIDR size. Default is false.
* `eni_subnet_ids` - (Optional) Subnet Ids for cluster with VPC-CNI network mode. This field can only set when field `network_type` is 'VPC-CNI'. `eni_subnet_ids` can not empty once be set.
* `exist_instance` - (Optional, ForceNew) create tke cluster by existed instances.
* `extra_args` - (Optional, ForceNew) Custom parameter information related to the node.
* `globe_desired_pod_num` - (Optional, ForceNew) Indicate to set desired pod number in node. valid when enable_customized_pod_cidr=true, and it takes effect for all nodes.
* `ignore_cluster_cidr_conflict` - (Optional, ForceNew) Indicates whether to ignore the cluster cidr conflict error. Default is false.
* `is_non_static_ip_mode` - (Optional, ForceNew) Indicates whether non-static ip mode is enabled. Default is false.
* `kube_proxy_mode` - (Optional) Cluster kube-proxy mode, the available values include: 'kube-proxy-bpf'. Default is not set.When set to kube-proxy-bpf, cluster version greater than 1.14 and with Tencent Linux 2.4 is required.
* `labels` - (Optional, ForceNew) Labels of tke cluster nodes.
* `managed_cluster_internet_security_policies` - (Optional) Security policies for managed cluster internet, like:'192.168.1.0/24' or '113.116.51.27', '0.0.0.0/0' means all. This field can only set when field `cluster_deploy_type` is 'MANAGED_CLUSTER' and `cluster_internet` is true. `managed_cluster_internet_security_policies` can not delete or empty once be set.
* `master_config` - (Optional, ForceNew) Deploy the machine configuration information of the 'MASTER_ETCD' service, and create <=7 units for common users.
* `mount_target` - (Optional, ForceNew) Mount target. Default is not mounting.
* `network_type` - (Optional, ForceNew) Cluster network type, GR or VPC-CNI. Default is GR.
* `node_name_type` - (Optional, ForceNew) Node name type of Cluster, the available values include: 'lan-ip' and 'hostname', Default is 'lan-ip'.
* `node_pool_global_config` - (Optional) Global config effective for all node pools.
* `project_id` - (Optional) Project ID, default value is 0.
* `service_cidr` - (Optional, ForceNew) A network address block of the service. Different from vpc cidr and cidr of other clusters within this vpc. Must be in  10./192.168/172.[16-31] segments.
* `tags` - (Optional) The tags of the cluster.
* `unschedulable` - (Optional, ForceNew) Sets whether the joining node participates in the schedule. Default is '0'. Participate in scheduling.
* `upgrade_instances_follow_cluster` - (Optional) Indicates whether upgrade all instances when cluster_version change. Default is false.
* `worker_config` - (Optional, ForceNew) Deploy the machine configuration information of the 'WORKER' service, and create <=20 units for common users. The other 'WORK' service are added by 'tencentcloud_kubernetes_worker'.

The `cluster_extra_args` object supports the following:

* `kube_apiserver` - (Optional, ForceNew) The customized parameters for kube-apiserver.
* `kube_controller_manager` - (Optional, ForceNew) The customized parameters for kube-controller-manager.
* `kube_scheduler` - (Optional, ForceNew) The customized parameters for kube-scheduler.

The `data_disk` object supports the following:

* `disk_size` - (Optional, ForceNew) Volume of disk in GB. Default is `0`.
* `disk_type` - (Optional, ForceNew) Types of disk, available values: `CLOUD_PREMIUM` and `CLOUD_SSD` and `CLOUD_HSSD`.
* `snapshot_id` - (Optional, ForceNew) Data disk snapshot ID.

The `exist_instance` object supports the following:

* `desired_pod_numbers` - (Optional, ForceNew) Custom mode cluster, you can specify the number of pods for each node. corresponding to the existed_instances_para.instance_ids parameter.
* `instances_para` - (Optional, ForceNew) Reinstallation parameters of an existing instance.
* `node_role` - (Optional, ForceNew) Role of existed node. value:MASTER_ETCD or WORKER.

The `instances_para` object supports the following:

* `instance_ids` - (Required, ForceNew) Cluster IDs.

The `master_config` object supports the following:

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

The `node_pool_global_config` object supports the following:

* `expander` - (Optional) Indicates which scale-out method will be used when there are multiple scaling groups. Valid values: `random` - select a random scaling group, `most-pods` - select the scaling group that can schedule the most pods, `least-waste` - select the scaling group that can ensure the fewest remaining resources after Pod scheduling.
* `ignore_daemon_sets_utilization` - (Optional) Whether to ignore DaemonSet pods by default when calculating resource usage.
* `is_scale_in_enabled` - (Optional) Indicates whether to enable scale-in.
* `max_concurrent_scale_in` - (Optional) Max concurrent scale-in volume.
* `scale_in_delay` - (Optional) Number of minutes after cluster scale-out when the system starts judging whether to perform scale-in.
* `scale_in_unneeded_time` - (Optional) Number of consecutive minutes of idleness after which the node is subject to scale-in.
* `scale_in_utilization_threshold` - (Optional) Percentage of node resource usage below which the node is considered to be idle.
* `skip_nodes_with_local_storage` - (Optional) During scale-in, ignore nodes with local storage pods.
* `skip_nodes_with_system_pods` - (Optional) During scale-in, ignore nodes with pods in the kube-system namespace that are not managed by DaemonSet.

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
* `certification_authority` - The certificate used for access.
* `cluster_external_endpoint` - External network address to access.
* `cluster_node_num` - Number of nodes in the cluster.
* `domain` - Domain name for access.
* `kube_config` - kubernetes config.
* `password` - Password of account.
* `pgw_endpoint` - The Intranet address used for access.
* `security_policy` - Access policy.
* `user_name` - User name of account.
* `worker_instances_list` - An information list of cvm within the 'WORKER' clusters. Each element contains the following attributes:
  * `failed_reason` - Information of the cvm when it is failed.
  * `instance_id` - ID of the cvm.
  * `instance_role` - Role of the cvm.
  * `instance_state` - State of the cvm.
  * `lan_ip` - LAN IP of the cvm.


