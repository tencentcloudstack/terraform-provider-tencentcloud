---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_native_node_pool"
sidebar_current: "docs-tencentcloud-resource-kubernetes_native_node_pool"
description: |-
  Provides a resource to create a tke kubernetes_native_node_pool
---

# tencentcloud_kubernetes_native_node_pool

Provides a resource to create a tke kubernetes_native_node_pool

## Example Usage

```hcl
resource "tencentcloud_kubernetes_native_node_pool" "kubernetes_native_node_pool" {
  cluster_id = "cls-eyier120"
  name       = "native-node-pool"
  type       = "Native"

  labels {
    name  = "test11"
    value = "test21"
  }

  taints {
    key    = "product"
    value  = "coderider"
    effect = "NoExecute"
  }

  tags {
    resource_type = "machine"
    tags {
      key   = "keep-test-np1"
      value = "test1"
    }
    tags {
      key   = "keep-test-np3"
      value = "test3"
    }
  }

  deletion_protection = false
  unschedulable       = false

  native {
    scaling {
      min_replicas  = 1
      max_replicas  = 10
      create_policy = "ZoneEquality"
    }
    subnet_ids           = ["subnet-itb6d123"]
    instance_charge_type = "PREPAID"
    system_disk {
      disk_type = "CLOUD_SSD"
      disk_size = 50
    }
    instance_types     = ["SA2.MEDIUM2"]
    security_group_ids = ["sg-7tum9120"]
    auto_repair        = false
    instance_charge_prepaid {
      period     = 1
      renew_flag = "NOTIFY_AND_MANUAL_RENEW"
    }
    management {
      nameservers = ["183.60.83.19", "183.60.82.98"]
      hosts       = ["192.168.2.42 static.fake.com", "192.168.2.42 static.fake.com2"]
      kernel_args = ["kernel.pid_max=65535", "fs.file-max=400000"]
    }
    host_name_pattern = "aaa{R:3}"
    kubelet_args      = ["allowed-unsafe-sysctls=net.core.somaxconn", "root-dir=/var/lib/test"]
    lifecycle {
      pre_init  = "ZWNobyBoZWxsb3dvcmxk"
      post_init = "ZWNobyBoZWxsb3dvcmxk"
    }
    runtime_root_dir   = "/var/lib/docker"
    enable_autoscaling = true
    replicas           = 2
    internet_accessible {
      max_bandwidth_out = 50
      charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    }
    data_disks {
      disk_type             = "CLOUD_PREMIUM"
      file_system           = "ext4"
      disk_size             = 60
      mount_target          = "/var/lib/containerd"
      auto_format_and_mount = true
    }
    key_ids = ["skey-9pcs2100"]
  }

  annotations {
    name  = "cluster-autoscaler.kubernetes.io/scale-down-disabled"
    value = "true"
  }
  annotations {
    name  = "node.tke.cloud.tencent.com/security-agent"
    value = "false"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) ID of the cluster.
* `name` - (Required, String) Node pool name.
* `native` - (Required, List) Native node pool creation parameters.
* `type` - (Required, String) Node pool type. Optional value is `Native`.
* `annotations` - (Optional, Set) Node Annotation List.
* `deletion_protection` - (Optional, Bool) Whether to enable deletion protection.
* `labels` - (Optional, List) Node Labels.
* `tags` - (Optional, List) Node tags.
* `taints` - (Optional, List) Node taint.
* `unschedulable` - (Optional, Bool) Whether the node is not schedulable by default. The native node is not aware of it and passes false by default.

The `annotations` object supports the following:

* `name` - (Required, String) Name in the map table.
* `value` - (Required, String) Value in the map table.

The `data_disks` object of `native` supports the following:

* `auto_format_and_mount` - (Required, Bool) Whether to automatically format the disk and mount it.
* `disk_size` - (Required, Int) Cloud disk size (G).
* `disk_type` - (Required, String) Cloud disk type. Valid values: `CLOUD_PREMIUM`: Premium Cloud Storage, `CLOUD_SSD`: cloud SSD disk, `CLOUD_BSSD`: Basic SSD, `CLOUD_HSSD`: Enhanced SSD, `CLOUD_TSSD`: Tremendous SSD, `LOCAL_NVME`: local NVME disk.
* `disk_partition` - (Optional, String) Mount device name or partition name.
* `encrypt` - (Optional, String) Pass in this parameter to create an encrypted cloud disk. The value is fixed to `ENCRYPT`.
* `file_system` - (Optional, String) File system (ext3/ext4/xfs).
* `kms_key_id` - (Optional, String) Customize the key when purchasing an encrypted disk. When this parameter is passed in, the Encrypt parameter is not empty.
* `mount_target` - (Optional, String) Mount directory.
* `snapshot_id` - (Optional, String) Snapshot ID. If passed in, the cloud disk will be created based on this snapshot. The snapshot type must be a data disk snapshot.
* `throughput_performance` - (Optional, Int) Cloud disk performance, unit: MB/s. Use this parameter to purchase additional performance for the cloud disk.

The `instance_charge_prepaid` object of `native` supports the following:

* `period` - (Required, Int) Postpaid billing cycle, unit (month): 1, 2, 3, 4, 5,, 6, 7, 8, 9, 10, 11, 12, 24, 36, 48, 60.
* `renew_flag` - (Optional, String) Prepaid renewal method:
  - `NOTIFY_AND_AUTO_RENEW`: Notify users of expiration and automatically renew (default).
  - `NOTIFY_AND_MANUAL_RENEW`: Notify users of expiration, but do not automatically renew.
  - `DISABLE_NOTIFY_AND_MANUAL_RENEW`: Do not notify users of expiration and do not automatically renew.

The `internet_accessible` object of `native` supports the following:

* `charge_type` - (Required, String) Network billing method. Optional value is `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.
* `max_bandwidth_out` - (Required, Int) Maximum bandwidth output. Note: When chargeType is `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_POSTPAID_BY_HOUR`, the valid range is 1~100. When chargeType is `BANDWIDTH_PACKAG`, the valid range is 1~2000.
* `bandwidth_package_id` - (Optional, String) Bandwidth package ID. Note: When ChargeType is BANDWIDTH_PACKAG, the value cannot be empty; otherwise, the value must be empty.

The `labels` object supports the following:

* `name` - (Required, String) Name in the map table.
* `value` - (Required, String) Value in map table.

The `lifecycle` object of `native` supports the following:

* `post_init` - (Optional, String) Custom script after node initialization.
* `pre_init` - (Optional, String) Custom script before node initialization.

The `management` object of `native` supports the following:

* `hosts` - (Optional, List) Hosts configuration.
* `kernel_args` - (Optional, List) Kernel parameter configuration.
* `nameservers` - (Optional, List) Dns configuration.

The `native` object supports the following:

* `instance_charge_type` - (Required, String, ForceNew) Node billing type. `PREPAID` is a yearly and monthly subscription, `POSTPAID_BY_HOUR` is a pay-as-you-go plan. The default is `POSTPAID_BY_HOUR`.
* `instance_types` - (Required, List) Model list.
* `security_group_ids` - (Required, List) Security group list.
* `subnet_ids` - (Required, List) Subnet list.
* `system_disk` - (Required, List, ForceNew) System disk configuration.
* `auto_repair` - (Optional, Bool) Whether to enable self-healing ability.
* `data_disks` - (Optional, List) Native node pool data disk list.
* `enable_autoscaling` - (Optional, Bool) Whether to enable elastic scaling.
* `health_check_policy_name` - (Optional, String) Fault self-healing rule name.
* `host_name_pattern` - (Optional, String) Native node pool hostName pattern string.
* `instance_charge_prepaid` - (Optional, List) Billing configuration for yearly and monthly models.
* `internet_accessible` - (Optional, List) Public network bandwidth settings.
* `key_ids` - (Optional, List) Node pool ssh public key id array.
* `kubelet_args` - (Optional, List) Kubelet custom parameters.
* `lifecycle` - (Optional, List) Predefined scripts.
* `management` - (Optional, List) Node pool management parameter settings.
* `replicas` - (Optional, Int) Desired number of nodes.
* `runtime_root_dir` - (Optional, String, ForceNew) Runtime root directory.
* `scaling` - (Optional, List) Node pool scaling configuration.

The `scaling` object of `native` supports the following:

* `create_policy` - (Optional, String) Node pool expansion strategy. `ZoneEquality`: multiple availability zones are broken up; `ZonePriority`: the preferred availability zone takes precedence.
* `max_replicas` - (Optional, Int) Maximum number of replicas in node pool.
* `min_replicas` - (Optional, Int) Minimum number of replicas in node pool.

The `system_disk` object of `native` supports the following:

* `disk_size` - (Required, Int, ForceNew) Cloud disk size (G).
* `disk_type` - (Required, String, ForceNew) Cloud disk type. Valid values: `CLOUD_PREMIUM`: Premium Cloud Storage, `CLOUD_SSD`: cloud SSD disk, `CLOUD_BSSD`: Basic SSD, `CLOUD_HSSD`: Enhanced SSD.

The `tags` object of `tags` supports the following:

* `key` - (Optional, String) Tag Key.
* `value` - (Optional, String) Tag Value.

The `tags` object supports the following:

* `resource_type` - (Optional, String) The resource type bound to the label.
* `tags` - (Optional, List) Tag pair list.

The `taints` object supports the following:

* `effect` - (Optional, String) Effect of the taint.
* `key` - (Optional, String) Key of the taint.
* `value` - (Optional, String) Value of the taint.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_at` - Creation time.
* `life_state` - Node pool status.


## Import

tke kubernetes_native_node_pool can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_native_node_pool.kubernetes_native_node_pool cls-xxx#np-xxx
```

