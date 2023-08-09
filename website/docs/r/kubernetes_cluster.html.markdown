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

~> **NOTE:** We recommend this usage that uses the `tencentcloud_kubernetes_cluster` resource to create a cluster without any `worker_config`, then adds nodes by the `tencentcloud_kubernetes_node_pool` resource.
It's more flexible than managing worker config directly with `tencentcloud_kubernetes_cluster`, `tencentcloud_kubernetes_scale_worker`, or existing node management of `tencentcloud_kubernetes_attachment`. The reason is that `worker_config` is unchangeable and may cause the whole cluster resource to `ForceNew`.

## Example Usage

### Create a basic cluster with two worker nodes

```hcl
variable "default_instance_type" {
  default = "SA2.2XLARGE16"
}

variable "availability_zone_first" {
  default = "ap-guangzhou-3"
}

variable "availability_zone_second" {
  default = "ap-guangzhou-4"
}

variable "example_cluster_cidr" {
  default = "10.31.0.0/16"
}

locals {
  first_vpc_id     = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.vpc_id
  first_subnet_id  = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.subnet_id
  second_vpc_id    = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.vpc_id
  second_subnet_id = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.subnet_id
  sg_id            = tencentcloud_security_group.sg.id
  image_id         = data.tencentcloud_images.default.image_id
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
  name = "tf-example-sg"
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

data "tencentcloud_images" "default" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                          = local.first_vpc_id
  cluster_cidr                    = var.example_cluster_cidr
  cluster_max_pod_num             = 32
  cluster_name                    = "tf_example_cluster"
  cluster_desc                    = "example for tke cluster"
  cluster_max_service_num         = 32
  cluster_internet                = false
  cluster_internet_security_group = local.sg_id
  cluster_version                 = "1.22.5"
  cluster_deploy_type             = "MANAGED_CLUSTER"

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_first
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.first_subnet_id
    img_id                     = local.image_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    # key_ids                   = ["skey-11112222"]
    password = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
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
    subnet_id                  = local.second_subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    key_ids                   = ["skey-11112222"]
    cam_role_name             = "CVM_QcsRole"
    # password                  = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
}
```

### Use Kubelet

```hcl
# Create a baisc kubernetes cluster with two nodes.
variable "default_instance_type" {
  default = "SA2.2XLARGE16"
}

variable "availability_zone_first" {
  default = "ap-guangzhou-3"
}

variable "availability_zone_second" {
  default = "ap-guangzhou-4"
}

variable "example_cluster_cidr" {
  default = "10.31.0.0/16"
}

locals {
  first_vpc_id     = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.vpc_id
  first_subnet_id  = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.subnet_id
  second_vpc_id    = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.vpc_id
  second_subnet_id = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.subnet_id
  sg_id            = tencentcloud_security_group.sg.id
  image_id         = data.tencentcloud_images.default.image_id
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
  name = "tf-example-sg"
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

data "tencentcloud_images" "default" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                          = local.first_vpc_id
  cluster_cidr                    = var.example_cluster_cidr
  cluster_max_pod_num             = 32
  cluster_name                    = "tf_example_cluster"
  cluster_desc                    = "example for tke cluster"
  cluster_max_service_num         = 32
  cluster_internet                = false
  cluster_internet_security_group = local.sg_id
  cluster_version                 = "1.22.5"
  cluster_deploy_type             = "MANAGED_CLUSTER"

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_first
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.first_subnet_id
    img_id                     = local.image_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
      encrypt   = false
    }

    enhanced_security_service  = false
    enhanced_monitor_service   = false
    user_data                  = "dGVzdA=="
    disaster_recover_group_ids = []
    security_group_ids         = []
    key_ids                    = []
    password                   = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
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
    subnet_id                  = local.second_subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service  = false
    enhanced_monitor_service   = false
    user_data                  = "dGVzdA=="
    disaster_recover_group_ids = []
    security_group_ids         = []
    key_ids                    = []
    cam_role_name              = "CVM_QcsRole"
    password                   = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
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

### Use extension addons

```hcl
variable "availability_zone_first" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "10.31.0.0/16"
}

variable "default_instance_type" {
  default = "S5.SMALL1"
}

data "tencentcloud_vpc_subnets" "vpc_first" {
  is_default        = true
  availability_zone = var.availability_zone_first
}

# fetch latest addon(chart) versions
data "tencentcloud_kubernetes_charts" "charts" {}

locals {
  chartNames    = data.tencentcloud_kubernetes_charts.charts.chart_list.*.name
  chartVersions = data.tencentcloud_kubernetes_charts.charts.chart_list.*.latest_version
  chartMap      = zipmap(local.chartNames, local.chartVersions)
}

resource "tencentcloud_kubernetes_cluster" "cluster_with_addon" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc_first.instance_list.0.vpc_id
  cluster_cidr            = var.cluster_cidr
  cluster_max_pod_num     = 32
  cluster_name            = "test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32
  cluster_internet        = true
  # managed_cluster_internet_security_policies = ["3.3.3.3", "1.1.1.1"]
  cluster_deploy_type = "MANAGED_CLUSTER"

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
    enhanced_security_service  = false
    enhanced_monitor_service   = false
    user_data                  = "dGVzdA=="
    # password                  = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
    key_ids = "skey-11112222"
  }

  extension_addon {
    name = "COS"
    param = jsonencode({
      "kind" : "App", "spec" : {
        "chart" : { "chartName" : "cos", "chartVersion" : local.chartMap["cos"] },
        "values" : { "values" : [], "rawValues" : "e30=", "rawValuesType" : "json" }
      }
    })
  }
  extension_addon {
    name = "SecurityGroupPolicy"
    param = jsonencode({
      "kind" : "App", "spec" : { "chart" : { "chartName" : "securitygrouppolicy", "chartVersion" : local.chartMap["securitygrouppolicy"] } }
    })
  }
  extension_addon {
    name = "OOMGuard"
    param = jsonencode({
      "kind" : "App", "spec" : { "chart" : { "chartName" : "oomguard", "chartVersion" : local.chartMap["oomguard"] } }
    })
  }
  extension_addon {
    name = "OLM"
    param = jsonencode({
      "kind" : "App", "spec" : { "chart" : { "chartName" : "olm", "chartVersion" : local.chartMap["olm"] } }
    })
  }
}
```

### Use node pool global config

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
  vpc_id                  = var.vpc
  cluster_cidr            = "10.1.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32
  cluster_internet        = true
  # managed_cluster_internet_security_policies = ["3.3.3.3", "1.1.1.1"]
  cluster_deploy_type = "MANAGED_CLUSTER"

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
    # password                  = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
    key_ids = "skey-11112222"
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

### Using VPC-CNI network type

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-1"
}

variable "vpc" {
  default = "vpc-r1m1fyx5"
}

variable "default_instance_type" {
  default = "SA2.SMALL2"
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = var.vpc
  cluster_max_pod_num     = 32
  cluster_name            = "test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 256
  cluster_internet        = true
  # managed_cluster_internet_security_policies = ["3.3.3.3", "1.1.1.1"]
  cluster_deploy_type = "MANAGED_CLUSTER"
  network_type        = "VPC-CNI"
  eni_subnet_ids      = ["subnet-bk1etlyu"]
  service_cidr        = "10.1.0.0/24"

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_PREMIUM"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = "subnet-t5dv27rs"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    # password                  = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
    key_ids = "skey-11112222"
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
}
```

### Using ops options

```hcl
resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  # ...your basic fields

  log_agent {
    enabled          = true
    kubelet_root_dir = "" # optional
  }

  event_persistence {
    enabled    = true
    log_set_id = "" # optional
    topic_id   = "" # optional
  }

  cluster_audit {
    enabled    = true
    log_set_id = "" # optional
    topic_id   = "" # optional
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, String, ForceNew) Vpc Id of the cluster.
* `acquire_cluster_admin_role` - (Optional, Bool) If set to true, it will acquire the ClusterRole tke:admin. NOTE: this arguments cannot revoke to `false` after acquired.
* `auth_options` - (Optional, List) Specify cluster authentication configuration. Only available for managed cluster and `cluster_version` >= 1.20.
* `auto_upgrade_cluster_level` - (Optional, Bool) Whether the cluster level auto upgraded, valid for managed cluster.
* `base_pod_num` - (Optional, Int, ForceNew) The number of basic pods. valid when enable_customized_pod_cidr=true.
* `claim_expired_seconds` - (Optional, Int) Claim expired seconds to recycle ENI. This field can only set when field `network_type` is 'VPC-CNI'. `claim_expired_seconds` must greater or equal than 300 and less than 15768000.
* `cluster_as_enabled` - (Optional, Bool, ForceNew, **Deprecated**) This argument is deprecated because the TKE auto-scaling group was no longer available. Indicates whether to enable cluster node auto scaling. Default is false.
* `cluster_audit` - (Optional, List) Specify Cluster Audit config. NOTE: Please make sure your TKE CamRole have permission to access CLS service.
* `cluster_cidr` - (Optional, String, ForceNew) A network address block of the cluster. Different from vpc cidr and cidr of other clusters within this vpc. Must be in  10./192.168/172.[16-31] segments.
* `cluster_deploy_type` - (Optional, String, ForceNew) Deployment type of the cluster, the available values include: 'MANAGED_CLUSTER' and 'INDEPENDENT_CLUSTER'. Default is 'MANAGED_CLUSTER'.
* `cluster_desc` - (Optional, String) Description of the cluster.
* `cluster_extra_args` - (Optional, List, ForceNew) Customized parameters for master component,such as kube-apiserver, kube-controller-manager, kube-scheduler.
* `cluster_internet_domain` - (Optional, String) Domain name for cluster Kube-apiserver internet access. Be careful if you modify value of this parameter, the cluster_external_endpoint value may be changed automatically too.
* `cluster_internet_security_group` - (Optional, String) Specify security group, NOTE: This argument must not be empty if cluster internet enabled.
* `cluster_internet` - (Optional, Bool) Open internet access or not. If this field is set 'true', the field below `worker_config` must be set. Because only cluster with node is allowed enable access endpoint.
* `cluster_intranet_domain` - (Optional, String) Domain name for cluster Kube-apiserver intranet access. Be careful if you modify value of this parameter, the pgw_endpoint value may be changed automatically too.
* `cluster_intranet_subnet_id` - (Optional, String) Subnet id who can access this independent cluster, this field must and can only set  when `cluster_intranet` is true. `cluster_intranet_subnet_id` can not modify once be set.
* `cluster_intranet` - (Optional, Bool) Open intranet access or not. If this field is set 'true', the field below `worker_config` must be set. Because only cluster with node is allowed enable access endpoint.
* `cluster_ipvs` - (Optional, Bool, ForceNew) Indicates whether `ipvs` is enabled. Default is true. False means `iptables` is enabled.
* `cluster_level` - (Optional, String) Specify cluster level, valid for managed cluster, use data source `tencentcloud_kubernetes_cluster_levels` to query available levels. Available value examples `L5`, `L20`, `L50`, `L100`, etc.
* `cluster_max_pod_num` - (Optional, Int, ForceNew) The maximum number of Pods per node in the cluster. Default is 256. The minimum value is 4. When its power unequal to 2, it will round upward to the closest power of 2.
* `cluster_max_service_num` - (Optional, Int, ForceNew) The maximum number of services in the cluster. Default is 256. The range is from 32 to 32768. When its power unequal to 2, it will round upward to the closest power of 2.
* `cluster_name` - (Optional, String) Name of the cluster.
* `cluster_os_type` - (Optional, String, ForceNew) Image type of the cluster os, the available values include: 'GENERAL'. Default is 'GENERAL'.
* `cluster_os` - (Optional, String, ForceNew) Operating system of the cluster, the available values include: 'centos7.6.0_x64','ubuntu18.04.1x86_64','tlinux2.4x86_64'. Default is 'tlinux2.4x86_64'.
* `cluster_version` - (Optional, String) Version of the cluster. Use `tencentcloud_kubernetes_available_cluster_versions` to get the upgradable cluster version.
* `container_runtime` - (Optional, String, ForceNew) Runtime type of the cluster, the available values include: 'docker' and 'containerd'.The Kubernetes v1.24 has removed dockershim, so please use containerd in v1.24 or higher.Default is 'docker'.
* `deletion_protection` - (Optional, Bool) Indicates whether cluster deletion protection is enabled. Default is false.
* `docker_graph_path` - (Optional, String, ForceNew) Docker graph path. Default is `/var/lib/docker`.
* `enable_customized_pod_cidr` - (Optional, Bool) Whether to enable the custom mode of node podCIDR size. Default is false.
* `eni_subnet_ids` - (Optional, List: [`String`]) Subnet Ids for cluster with VPC-CNI network mode. This field can only set when field `network_type` is 'VPC-CNI'. `eni_subnet_ids` can not empty once be set.
* `event_persistence` - (Optional, List) Specify cluster Event Persistence config. NOTE: Please make sure your TKE CamRole have permission to access CLS service.
* `exist_instance` - (Optional, List, ForceNew) create tke cluster by existed instances.
* `extension_addon` - (Optional, List) Information of the add-on to be installed.
* `extra_args` - (Optional, List: [`String`], ForceNew) Custom parameter information related to the node.
* `globe_desired_pod_num` - (Optional, Int, ForceNew) Indicate to set desired pod number in node. valid when enable_customized_pod_cidr=true, and it takes effect for all nodes.
* `ignore_cluster_cidr_conflict` - (Optional, Bool, ForceNew) Indicates whether to ignore the cluster cidr conflict error. Default is false.
* `is_non_static_ip_mode` - (Optional, Bool, ForceNew) Indicates whether non-static ip mode is enabled. Default is false.
* `kube_proxy_mode` - (Optional, String) Cluster kube-proxy mode, the available values include: 'kube-proxy-bpf'. Default is not set.When set to kube-proxy-bpf, cluster version greater than 1.14 and with Tencent Linux 2.4 is required.
* `labels` - (Optional, Map, ForceNew) Labels of tke cluster nodes.
* `log_agent` - (Optional, List) Specify cluster log agent config.
* `managed_cluster_internet_security_policies` - (Optional, List: [`String`], **Deprecated**) this argument was deprecated, use `cluster_internet_security_group` instead. Security policies for managed cluster internet, like:'192.168.1.0/24' or '113.116.51.27', '0.0.0.0/0' means all. This field can only set when field `cluster_deploy_type` is 'MANAGED_CLUSTER' and `cluster_internet` is true. `managed_cluster_internet_security_policies` can not delete or empty once be set.
* `master_config` - (Optional, List, ForceNew) Deploy the machine configuration information of the 'MASTER_ETCD' service, and create <=7 units for common users.
* `mount_target` - (Optional, String, ForceNew) Mount target. Default is not mounting.
* `network_type` - (Optional, String, ForceNew) Cluster network type, GR or VPC-CNI. Default is GR.
* `node_name_type` - (Optional, String, ForceNew) Node name type of Cluster, the available values include: 'lan-ip' and 'hostname', Default is 'lan-ip'.
* `node_pool_global_config` - (Optional, List) Global config effective for all node pools.
* `project_id` - (Optional, Int) Project ID, default value is 0.
* `runtime_version` - (Optional, String) Container Runtime version.
* `service_cidr` - (Optional, String, ForceNew) A network address block of the service. Different from vpc cidr and cidr of other clusters within this vpc. Must be in  10./192.168/172.[16-31] segments.
* `tags` - (Optional, Map) The tags of the cluster.
* `unschedulable` - (Optional, Int, ForceNew) Sets whether the joining node participates in the schedule. Default is '0'. Participate in scheduling.
* `upgrade_instances_follow_cluster` - (Optional, Bool) Indicates whether upgrade all instances when cluster_version change. Default is false.
* `worker_config` - (Optional, List, ForceNew) Deploy the machine configuration information of the 'WORKER' service, and create <=20 units for common users. The other 'WORK' service are added by 'tencentcloud_kubernetes_worker'.

The `auth_options` object supports the following:

* `auto_create_discovery_anonymous_auth` - (Optional, Bool) If set to `true`, the rbac rule will be created automatically which allow anonymous user to access '/.well-known/openid-configuration' and '/openid/v1/jwks'.
* `issuer` - (Optional, String) Specify service-account-issuer. If use_tke_default is set to `true`, please do not set this field, it will be ignored anyway.
* `jwks_uri` - (Optional, String) Specify service-account-jwks-uri. If use_tke_default is set to `true`, please do not set this field, it will be ignored anyway.
* `use_tke_default` - (Optional, Bool) If set to `true`, the issuer and jwks_uri will be generated automatically by tke, please do not set issuer and jwks_uri, and they will be ignored.

The `cluster_audit` object supports the following:

* `enabled` - (Required, Bool) Specify weather the Cluster Audit enabled. NOTE: Enable Cluster Audit will also auto install Log Agent.
* `delete_audit_log_and_topic` - (Optional, Bool) when you want to close the cluster audit log or delete the cluster, you can use this parameter to determine whether the audit log set and topic created by default will be deleted.
* `log_set_id` - (Optional, String) Specify id of existing CLS log set, or auto create a new set by leave it empty.
* `topic_id` - (Optional, String) Specify id of existing CLS log topic, or auto create a new topic by leave it empty.

The `cluster_extra_args` object supports the following:

* `kube_apiserver` - (Optional, List, ForceNew) The customized parameters for kube-apiserver.
* `kube_controller_manager` - (Optional, List, ForceNew) The customized parameters for kube-controller-manager.
* `kube_scheduler` - (Optional, List, ForceNew) The customized parameters for kube-scheduler.

The `data_disk` object supports the following:

* `auto_format_and_mount` - (Optional, Bool, ForceNew) Indicate whether to auto format and mount or not. Default is `false`.
* `disk_partition` - (Optional, String, ForceNew) The name of the device or partition to mount.
* `disk_size` - (Optional, Int, ForceNew) Volume of disk in GB. Default is `0`.
* `disk_type` - (Optional, String, ForceNew) Types of disk, available values: `CLOUD_PREMIUM` and `CLOUD_SSD` and `CLOUD_HSSD` and `CLOUD_TSSD`.
* `encrypt` - (Optional, Bool) Indicates whether to encrypt data disk, default `false`.
* `file_system` - (Optional, String, ForceNew) File system, e.g. `ext3/ext4/xfs`.
* `kms_key_id` - (Optional, String) ID of the custom CMK in the format of UUID or `kms-abcd1234`. This parameter is used to encrypt cloud disks.
* `mount_target` - (Optional, String, ForceNew) Mount target.
* `snapshot_id` - (Optional, String, ForceNew) Data disk snapshot ID.

The `event_persistence` object supports the following:

* `enabled` - (Required, Bool) Specify weather the Event Persistence enabled.
* `delete_event_log_and_topic` - (Optional, Bool) when you want to close the cluster event persistence or delete the cluster, you can use this parameter to determine whether the event persistence log set and topic created by default will be deleted.
* `log_set_id` - (Optional, String) Specify id of existing CLS log set, or auto create a new set by leave it empty.
* `topic_id` - (Optional, String) Specify id of existing CLS log topic, or auto create a new topic by leave it empty.

The `exist_instance` object supports the following:

* `desired_pod_numbers` - (Optional, List, ForceNew) Custom mode cluster, you can specify the number of pods for each node. corresponding to the existed_instances_para.instance_ids parameter.
* `instances_para` - (Optional, List, ForceNew) Reinstallation parameters of an existing instance.
* `node_role` - (Optional, String, ForceNew) Role of existed node. value:MASTER_ETCD or WORKER.

The `extension_addon` object supports the following:

* `name` - (Required, String) Add-on name.
* `param` - (Required, String) Parameter of the add-on resource object in JSON string format, please check the example at the top of page for reference.

The `instances_para` object supports the following:

* `instance_ids` - (Required, List, ForceNew) Cluster IDs.

The `log_agent` object supports the following:

* `enabled` - (Required, Bool) Whether the log agent enabled.
* `kubelet_root_dir` - (Optional, String) Kubelet root directory as the literal.

The `master_config` object supports the following:

* `instance_type` - (Required, String, ForceNew) Specified types of CVM instance.
* `subnet_id` - (Required, String, ForceNew) Private network ID.
* `availability_zone` - (Optional, String, ForceNew) Indicates which availability zone will be used.
* `bandwidth_package_id` - (Optional, String) bandwidth package id. if user is standard user, then the bandwidth_package_id is needed, or default has bandwidth_package_id.
* `cam_role_name` - (Optional, String, ForceNew) CAM role name authorized to access.
* `count` - (Optional, Int, ForceNew) Number of cvm.
* `data_disk` - (Optional, List, ForceNew) Configurations of data disk.
* `desired_pod_num` - (Optional, Int, ForceNew) Indicate to set desired pod number in node. valid when enable_customized_pod_cidr=true, and it override `[globe_]desired_pod_num` for current node. Either all the fields `desired_pod_num` or none.
* `disaster_recover_group_ids` - (Optional, List, ForceNew) Disaster recover groups to which a CVM instance belongs. Only support maximum 1.
* `enhanced_monitor_service` - (Optional, Bool, ForceNew) To specify whether to enable cloud monitor service. Default is TRUE.
* `enhanced_security_service` - (Optional, Bool, ForceNew) To specify whether to enable cloud security service. Default is TRUE.
* `hostname` - (Optional, String, ForceNew) The host name of the attached instance. Dot (.) and dash (-) cannot be used as the first and last characters of HostName and cannot be used consecutively. Windows example: The length of the name character is [2, 15], letters (capitalization is not restricted), numbers and dashes (-) are allowed, dots (.) are not supported, and not all numbers are allowed. Examples of other types (Linux, etc.): The character length is [2, 60], and multiple dots are allowed. There is a segment between the dots. Each segment allows letters (with no limitation on capitalization), numbers and dashes (-).
* `hpc_cluster_id` - (Optional, String) Id of cvm hpc cluster.
* `img_id` - (Optional, String) The valid image id, format of img-xxx.
* `instance_charge_type_prepaid_period` - (Optional, Int, ForceNew) The tenancy (time unit is month) of the prepaid instance. NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.
* `instance_charge_type_prepaid_renew_flag` - (Optional, String, ForceNew) Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.
* `instance_charge_type` - (Optional, String, ForceNew) The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. The default is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`, `PREPAID` instance will not terminated after cluster deleted, and may not allow to delete before expired.
* `instance_name` - (Optional, String, ForceNew) Name of the CVMs.
* `internet_charge_type` - (Optional, String, ForceNew) Charge types for network traffic. Available values include `TRAFFIC_POSTPAID_BY_HOUR`.
* `internet_max_bandwidth_out` - (Optional, Int) Max bandwidth of Internet access in Mbps. Default is 0.
* `key_ids` - (Optional, List, ForceNew) ID list of keys, should be set if `password` not set.
* `password` - (Optional, String, ForceNew) Password to access, should be set if `key_ids` not set.
* `public_ip_assigned` - (Optional, Bool, ForceNew) Specify whether to assign an Internet IP address.
* `security_group_ids` - (Optional, List, ForceNew) Security groups to which a CVM instance belongs.
* `system_disk_size` - (Optional, Int, ForceNew) Volume of system disk in GB. Default is `50`.
* `system_disk_type` - (Optional, String, ForceNew) System disk type. For more information on limits of system disk types, see [Storage Overview](https://intl.cloud.tencent.com/document/product/213/4952). Valid values: `LOCAL_BASIC`: local disk, `LOCAL_SSD`: local SSD disk, `CLOUD_SSD`: SSD, `CLOUD_PREMIUM`: Premium Cloud Storage. NOTE: `CLOUD_BASIC`, `LOCAL_BASIC` and `LOCAL_SSD` are deprecated.
* `user_data` - (Optional, String, ForceNew) ase64-encoded User Data text, the length limit is 16KB.

The `node_pool_global_config` object supports the following:

* `expander` - (Optional, String) Indicates which scale-out method will be used when there are multiple scaling groups. Valid values: `random` - select a random scaling group, `most-pods` - select the scaling group that can schedule the most pods, `least-waste` - select the scaling group that can ensure the fewest remaining resources after Pod scheduling.
* `ignore_daemon_sets_utilization` - (Optional, Bool) Whether to ignore DaemonSet pods by default when calculating resource usage.
* `is_scale_in_enabled` - (Optional, Bool) Indicates whether to enable scale-in.
* `max_concurrent_scale_in` - (Optional, Int) Max concurrent scale-in volume.
* `scale_in_delay` - (Optional, Int) Number of minutes after cluster scale-out when the system starts judging whether to perform scale-in.
* `scale_in_unneeded_time` - (Optional, Int) Number of consecutive minutes of idleness after which the node is subject to scale-in.
* `scale_in_utilization_threshold` - (Optional, Int) Percentage of node resource usage below which the node is considered to be idle.
* `skip_nodes_with_local_storage` - (Optional, Bool) During scale-in, ignore nodes with local storage pods.
* `skip_nodes_with_system_pods` - (Optional, Bool) During scale-in, ignore nodes with pods in the kube-system namespace that are not managed by DaemonSet.

The `worker_config` object supports the following:

* `instance_type` - (Required, String, ForceNew) Specified types of CVM instance.
* `subnet_id` - (Required, String, ForceNew) Private network ID.
* `availability_zone` - (Optional, String, ForceNew) Indicates which availability zone will be used.
* `bandwidth_package_id` - (Optional, String) bandwidth package id. if user is standard user, then the bandwidth_package_id is needed, or default has bandwidth_package_id.
* `cam_role_name` - (Optional, String, ForceNew) CAM role name authorized to access.
* `count` - (Optional, Int, ForceNew) Number of cvm.
* `data_disk` - (Optional, List, ForceNew) Configurations of data disk.
* `desired_pod_num` - (Optional, Int, ForceNew) Indicate to set desired pod number in node. valid when enable_customized_pod_cidr=true, and it override `[globe_]desired_pod_num` for current node. Either all the fields `desired_pod_num` or none.
* `disaster_recover_group_ids` - (Optional, List, ForceNew) Disaster recover groups to which a CVM instance belongs. Only support maximum 1.
* `enhanced_monitor_service` - (Optional, Bool, ForceNew) To specify whether to enable cloud monitor service. Default is TRUE.
* `enhanced_security_service` - (Optional, Bool, ForceNew) To specify whether to enable cloud security service. Default is TRUE.
* `hostname` - (Optional, String, ForceNew) The host name of the attached instance. Dot (.) and dash (-) cannot be used as the first and last characters of HostName and cannot be used consecutively. Windows example: The length of the name character is [2, 15], letters (capitalization is not restricted), numbers and dashes (-) are allowed, dots (.) are not supported, and not all numbers are allowed. Examples of other types (Linux, etc.): The character length is [2, 60], and multiple dots are allowed. There is a segment between the dots. Each segment allows letters (with no limitation on capitalization), numbers and dashes (-).
* `hpc_cluster_id` - (Optional, String) Id of cvm hpc cluster.
* `img_id` - (Optional, String) The valid image id, format of img-xxx.
* `instance_charge_type_prepaid_period` - (Optional, Int, ForceNew) The tenancy (time unit is month) of the prepaid instance. NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.
* `instance_charge_type_prepaid_renew_flag` - (Optional, String, ForceNew) Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.
* `instance_charge_type` - (Optional, String, ForceNew) The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. The default is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`, `PREPAID` instance will not terminated after cluster deleted, and may not allow to delete before expired.
* `instance_name` - (Optional, String, ForceNew) Name of the CVMs.
* `internet_charge_type` - (Optional, String, ForceNew) Charge types for network traffic. Available values include `TRAFFIC_POSTPAID_BY_HOUR`.
* `internet_max_bandwidth_out` - (Optional, Int) Max bandwidth of Internet access in Mbps. Default is 0.
* `key_ids` - (Optional, List, ForceNew) ID list of keys, should be set if `password` not set.
* `password` - (Optional, String, ForceNew) Password to access, should be set if `key_ids` not set.
* `public_ip_assigned` - (Optional, Bool, ForceNew) Specify whether to assign an Internet IP address.
* `security_group_ids` - (Optional, List, ForceNew) Security groups to which a CVM instance belongs.
* `system_disk_size` - (Optional, Int, ForceNew) Volume of system disk in GB. Default is `50`.
* `system_disk_type` - (Optional, String, ForceNew) System disk type. For more information on limits of system disk types, see [Storage Overview](https://intl.cloud.tencent.com/document/product/213/4952). Valid values: `LOCAL_BASIC`: local disk, `LOCAL_SSD`: local SSD disk, `CLOUD_SSD`: SSD, `CLOUD_PREMIUM`: Premium Cloud Storage. NOTE: `CLOUD_BASIC`, `LOCAL_BASIC` and `LOCAL_SSD` are deprecated.
* `user_data` - (Optional, String, ForceNew) ase64-encoded User Data text, the length limit is 16KB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `certification_authority` - The certificate used for access.
* `cluster_external_endpoint` - External network address to access.
* `cluster_node_num` - Number of nodes in the cluster.
* `domain` - Domain name for access.
* `kube_config_intranet` - Kubernetes config of private network.
* `kube_config` - Kubernetes config.
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


