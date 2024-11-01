---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster_master_attachment"
sidebar_current: "docs-tencentcloud-resource-kubernetes_cluster_master_attachment"
description: |-
  Provides a resource to create a tke kubernetes cluster master attachment
---

# tencentcloud_kubernetes_cluster_master_attachment

Provides a resource to create a tke kubernetes cluster master attachment

## Example Usage

```hcl
resource "tencentcloud_kubernetes_cluster_master_attachment" "example" {
  cluster_id                  = "cls-fp5o961e"
  instance_id                 = "ins-7d6tpbyg"
  node_role                   = "MASTER_ETCD"
  enhanced_security_service   = true
  enhanced_monitor_service    = true
  enhanced_automation_service = true
  password                    = "Password@123"
  security_group_ids          = ["sg-hjs685q9"]

  master_config {
    mount_target      = "/var/data"
    docker_graph_path = "/var/lib/containerd"
    unschedulable     = 0
    labels {
      name  = "key"
      value = "value"
    }

    data_disk {
      file_system           = "ext4"
      auto_format_and_mount = true
      mount_target          = "/var/data"
      disk_partition        = "/dev/vdb"
    }

    extra_args {
      kubelet = ["root-dir=/root"]
    }

    taints {
      key    = "key"
      value  = "value"
      effect = "NoSchedule"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) ID of the cluster.
* `instance_id` - (Required, String, ForceNew) ID of the CVM instance, this cvm will reinstall the system.
* `node_role` - (Required, String, ForceNew) Node role, values: MASTER_ETCD, WORKER. MASTER_ETCD needs to be specified only when creating an INDEPENDENT_CLUSTER independent cluster. The number of MASTER_ETCD nodes is 3-7, and it is recommended to have an odd number. The minimum configuration for MASTER_ETCD is 4C8G.
* `desired_pod_numbers` - (Optional, List: [`Int`], ForceNew) When the node belongs to the podCIDR size customization mode, the maximum number of pods running on the node can be specified.
* `enhanced_automation_service` - (Optional, Bool, ForceNew) Activate TencentCloud Automation Tools (TAT) service. If this parameter is not specified, the public image will default to enabling the Cloud Automation Assistant service, while other images will default to not enabling the Cloud Automation Assistant service.
* `enhanced_monitor_service` - (Optional, Bool, ForceNew) To specify whether to enable cloud monitor service. Default is TRUE.
* `enhanced_security_service` - (Optional, Bool, ForceNew) To specify whether to enable cloud security service. Default is TRUE.
* `extra_args` - (Optional, List, ForceNew) Custom parameters for cluster master component.
* `host_name` - (Optional, String, ForceNew) When reinstalling the system, you can specify the HostName of the instance to be modified (this parameter must be passed when the cluster is in HostName mode, and the rule name should be consistent with the HostName of the CVM instance creation interface except that uppercase characters are not supported).
* `key_ids` - (Optional, List: [`String`], ForceNew) The key pair to use for the instance, it looks like skey-16jig7tx, it should be set if `password` not set.
* `master_config` - (Optional, List, ForceNew) Advanced Node Settings. commonly used to attach existing instances.
* `password` - (Optional, String, ForceNew) Password to access, should be set if `key_ids` not set.
* `security_group_ids` - (Optional, List: [`String`], ForceNew) The security group to which the instance belongs. This parameter can be obtained by calling the sgId field in the return value of DescribeSecureGroups. If this parameter is not specified, the default security group will be bound.

The `data_disk` object of `master_config` supports the following:

* `auto_format_and_mount` - (Optional, Bool, ForceNew) Indicate whether to auto format and mount or not. Default is `false`.
* `disk_partition` - (Optional, String, ForceNew) The name of the device or partition to mount. NOTE: this argument doesn't support setting in node pool, or will leads to mount error.
* `disk_size` - (Optional, Int, ForceNew) Volume of disk in GB. Default is `0`.
* `disk_type` - (Optional, String, ForceNew) Types of disk. Valid value: `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_BASIC`, `CLOUD_PREMIUM`, `CLOUD_SSD`, `CLOUD_HSSD`, `CLOUD_TSSD` and `CLOUD_BSSD`.
* `file_system` - (Optional, String, ForceNew) File system, e.g. `ext3/ext4/xfs`.
* `mount_target` - (Optional, String, ForceNew) Mount target.

The `extra_args` object of `master_config` supports the following:

* `kubelet` - (Optional, List, ForceNew) Kubelet custom parameter. The parameter format is ["k1=v1", "k1=v2"].

The `extra_args` object supports the following:

* `etcd` - (Optional, Set, ForceNew) etcd custom parameters. Only supports independent clusters.
* `kube_api_server` - (Optional, Set, ForceNew) Kube apiserver custom parameters. The parameter format is ["k1=v1", "k1=v2"].
* `kube_controller_manager` - (Optional, Set, ForceNew) Kube controller manager custom parameters.
* `kube_scheduler` - (Optional, Set, ForceNew) kube scheduler custom parameters.

The `gpu_args` object of `master_config` supports the following:

* `cuda` - (Optional, Map, ForceNew) CUDA  version. Format like: `{ version: String, name: String }`. `version`: Version of GPU driver or CUDA; `name`: Name of GPU driver or CUDA.
* `cudnn` - (Optional, Map, ForceNew) cuDNN version. Format like: `{ version: String, name: String, doc_name: String, dev_name: String }`. `version`: cuDNN version; `name`: cuDNN name; `doc_name`: Doc name of cuDNN; `dev_name`: Dev name of cuDNN.
* `custom_driver` - (Optional, Map, ForceNew) Custom GPU driver. Format like: `{address: String}`. `address`: URL of custom GPU driver address.
* `driver` - (Optional, Map, ForceNew) GPU driver version. Format like: `{ version: String, name: String }`. `version`: Version of GPU driver or CUDA; `name`: Name of GPU driver or CUDA.
* `mig_enable` - (Optional, Bool, ForceNew) Whether to enable MIG.

The `labels` object of `master_config` supports the following:

* `name` - (Required, String, ForceNew) Name of map.
* `value` - (Required, String, ForceNew) Value of map.

The `master_config` object supports the following:

* `data_disk` - (Optional, List, ForceNew) Configurations of data disk.
* `desired_pod_number` - (Optional, Int, ForceNew) Indicate to set desired pod number in node. valid when the cluster is podCIDR.
* `docker_graph_path` - (Optional, String, ForceNew) Docker graph path. Default is `/var/lib/docker`.
* `extra_args` - (Optional, List, ForceNew) Custom parameter information related to the node. This is a white-list parameter.
* `gpu_args` - (Optional, List, ForceNew) GPU driver parameters.
* `labels` - (Optional, List, ForceNew) Node label list.
* `mount_target` - (Optional, String, ForceNew) Mount target. Default is not mounting.
* `taints` - (Optional, List, ForceNew) Node taint.
* `unschedulable` - (Optional, Int, ForceNew) Set whether the joined nodes participate in scheduling, with a default value of 0, indicating participation in scheduling; Non 0 means not participating in scheduling.
* `user_script` - (Optional, String, ForceNew) User script encoded in base64, which will be executed after the k8s component runs. The user needs to ensure the script's reentrant and retry logic. The script and its generated log files can be viewed in the node path /data/ccs_userscript/. If the node needs to be initialized before joining the schedule, it can be used in conjunction with the `unschedulable` parameter. After the final initialization of the userScript is completed, add the command "kubectl uncordon nodename --kubeconfig=/root/.kube/config" to add the node to the schedule.

The `taints` object of `master_config` supports the following:

* `effect` - (Optional, String, ForceNew) Effect of the taint.
* `key` - (Optional, String, ForceNew) Key of the taint.
* `value` - (Optional, String, ForceNew) Value of the taint.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



