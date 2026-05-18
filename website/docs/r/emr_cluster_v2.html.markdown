---
subcategory: "MapReduce(EMR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_emr_cluster_v2"
sidebar_current: "docs-tencentcloud-resource-emr_cluster_v2"
description: |-
  Provides a resource to create a EMR cluster (v2) using the `CreateCluster` API.
---

# tencentcloud_emr_cluster_v2

Provides a resource to create a EMR cluster (v2) using the `CreateCluster` API.

~> **NOTE:** The current resource `tencentcloud_emr_cluster_v2` does not yet support in-place updates. Any field change after creation will be silently ignored; destroy and re-create the resource to apply modifications. Modify APIs (scale-out/scale-in, tag modification, rename, ...) will be added in a follow-up release.

~> **NOTE:** Node counts are driven implicitly by the number of `master_resource_spec`, `core_resource_spec`, `task_resource_spec`, and `common_resource_spec` blocks declared. Each block represents one node of that role, and the first block is used as the resource template passed to the `CreateCluster` API.

## Example Usage

```hcl
resource "tencentcloud_emr_cluster_v2" "example" {
  product_version        = "EMR-V3.5.0"
  enable_support_ha_flag = false
  instance_name          = "tf-example-emr-v2"
  instance_charge_type   = "POSTPAID_BY_HOUR"
  need_master_wan        = "NEED_MASTER_WAN"

  login_settings {
    password = "Tencent@2026"
  }

  scene_software_config {
    scene_name = "Hadoop-Default"
    software = [
      "hdfs-3.2.2",
      "yarn-3.2.2",
      "zookeeper-3.6.3",
      "openldap-2.4.44",
      "knox-1.6.1",
      "hive-3.1.3",
    ]
  }

  meta_db_info {
    meta_type = "EMR_DEFAULT_META"
  }

  tags {
    tag_key   = "createBy"
    tag_value = "terraform"
  }

  zone_resource_configuration {
    virtual_private_cloud {
      vpc_id    = "vpc-xxxxxxxx"
      subnet_id = "subnet-xxxxxxxx"
    }

    placement {
      zone       = "ap-guangzhou-7"
      project_id = 0
    }

    all_node_resource_spec {
      # One master_resource_spec block = 1 master node.
      master_resource_spec {
        instance_type = "S6.2XLARGE32"

        system_disk {
          count     = 1
          disk_size = 50
          disk_type = "CLOUD_HSSD"
        }

        data_disk {
          count     = 1
          disk_size = 100
          disk_type = "CLOUD_SSD"
        }
      }

      # Two core_resource_spec blocks = 2 core nodes (first block is the template).
      core_resource_spec {
        instance_type = "SA4.8XLARGE64"

        system_disk {
          count     = 1
          disk_size = 50
          disk_type = "CLOUD_SSD"
        }

        data_disk {
          count     = 1
          disk_size = 300
          disk_type = "CLOUD_TSSD"
        }
      }

      core_resource_spec {
        instance_type = "SA4.8XLARGE64"

        system_disk {
          count     = 1
          disk_size = 50
          disk_type = "CLOUD_SSD"
        }

        data_disk {
          count     = 1
          disk_size = 300
          disk_type = "CLOUD_TSSD"
        }
      }
    }
  }
}
```

### Prepaid cluster with auto-renew

```hcl
resource "tencentcloud_emr_cluster_v2" "prepaid" {
  product_version        = "EMR-V3.5.0"
  enable_support_ha_flag = true
  instance_name          = "tf-example-emr-v2-prepaid"
  instance_charge_type   = "PREPAID"

  instance_charge_prepaid {
    period     = 1
    renew_flag = true
  }

  login_settings {
    password = "Tencent@2026"
  }

  scene_software_config {
    scene_name = "Hadoop-Default"
    software   = ["hdfs-3.2.2", "yarn-3.2.2", "zookeeper-3.6.3"]
  }

  zone_resource_configuration {
    virtual_private_cloud {
      vpc_id    = "vpc-xxxxxxxx"
      subnet_id = "subnet-xxxxxxxx"
    }

    placement {
      zone       = "ap-guangzhou-7"
      project_id = 0
    }

    all_node_resource_spec {
      # Two master_resource_spec blocks = 2 master nodes (HA).
      master_resource_spec {
        instance_type = "S6.2XLARGE32"
        system_disk {
          count     = 1
          disk_size = 50
          disk_type = "CLOUD_HSSD"
        }
      }

      master_resource_spec {
        instance_type = "S6.2XLARGE32"
        system_disk {
          count     = 1
          disk_size = 50
          disk_type = "CLOUD_HSSD"
        }
      }

      # Three core_resource_spec blocks = 3 core nodes.
      core_resource_spec {
        instance_type = "SA4.8XLARGE64"
        system_disk {
          count     = 1
          disk_size = 50
          disk_type = "CLOUD_SSD"
        }
        data_disk {
          count     = 2
          disk_size = 300
          disk_type = "CLOUD_TSSD"
        }
      }

      core_resource_spec {
        instance_type = "SA4.8XLARGE64"
        system_disk {
          count     = 1
          disk_size = 50
          disk_type = "CLOUD_SSD"
        }
        data_disk {
          count     = 2
          disk_size = 300
          disk_type = "CLOUD_TSSD"
        }
      }

      core_resource_spec {
        instance_type = "SA4.8XLARGE64"
        system_disk {
          count     = 1
          disk_size = 50
          disk_type = "CLOUD_SSD"
        }
        data_disk {
          count     = 2
          disk_size = 300
          disk_type = "CLOUD_TSSD"
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `enable_support_ha_flag` - (Required, Bool) Whether to enable node high availability. `true` means enabled, `false` means disabled.
* `instance_charge_type` - (Required, String) Instance billing mode. Valid values: `PREPAID` (monthly/yearly subscription), `POSTPAID_BY_HOUR` (pay-as-you-go).
* `instance_name` - (Required, String) Instance name. Length 6-36 characters. Only Chinese, letters, numbers, `-`, `_` are allowed.
* `login_settings` - (Required, List) Login settings for purchased nodes.
* `product_version` - (Required, String) EMR product version name, e.g., `EMR-V3.5.0`.
* `scene_software_config` - (Required, List) Cluster scenario and components to deploy.
* `client_token` - (Optional, String) Unique random token for idempotency (5-minute validity).
* `cos_bucket` - (Optional, String) COS bucket path, used when creating StarRocks storage-compute separation clusters.
* `custom_conf` - (Optional, String) Custom software configuration in JSON format.
* `default_meta_version` - (Optional, String) Default metadata DB version. Valid values: `mysql8`, `tdsql8`, `mysql5`.
* `depend_service` - (Optional, List) Shared component dependency information.
* `disaster_recover_group_ids` - (Optional, List: [`String`]) Spread placement group IDs. Currently supports only one ID.
* `enable_cbs_encrypt_flag` - (Optional, Bool) Whether to enable cluster-level CBS encryption. Default is false.
* `enable_kerberos_flag` - (Optional, Bool) Whether to enable Kerberos authentication. Default is false.
* `enable_remote_login_flag` - (Optional, Bool) Whether to enable external remote login. Invalid when `security_group_ids` is set. Default is false.
* `instance_charge_prepaid` - (Optional, List) Prepaid (monthly/yearly) billing parameters. Required when `instance_charge_type` is `PREPAID`.
* `load_balancer_id` - (Optional, String) CLB instance ID, e.g., `lb-xxxxxxxx`.
* `meta_db_info` - (Optional, List) Metadata database information. When `meta_type` is `EMR_NEW_META`/`EMR_DEFAULT_META`, no extra fields are required; when `EMR_EXIT_META`, `unify_meta_instance_id` must be set; when `USER_CUSTOM_META`, `meta_data_jdbc_url`/`meta_data_user`/`meta_data_pass` must be set.
* `need_cdb_audit` - (Optional, Int) Whether to enable database auditing.
* `need_master_wan` - (Optional, String) Whether to enable master public network. Valid values: `NEED_MASTER_WAN` (default), `NOT_NEED_MASTER_WAN`.
* `node_marks` - (Optional, List) Node identification information (TF platform use only).
* `partition_number` - (Optional, Int) Partition placement group partition number.
* `script_bootstrap_action_config` - (Optional, List) Bootstrap script configurations.
* `security_group_ids` - (Optional, List: [`String`]) Security group IDs bound to the instance, e.g., `["sg-xxxxxxxx"]`.
* `sg_ip` - (Optional, String) Security source IP, e.g., `10.0.0.0/8`.
* `tags` - (Optional, List) Tags to bind to the cluster instance.
* `web_ui_version` - (Optional, Int) Service UI address version. `0`: single URL (default); `1`: all URLs.
* `zone_resource_configuration` - (Optional, List) Per-zone resource configuration. One entry per AZ (first: primary, second: backup, third: arbitration).

The `all_node_resource_spec` object of `zone_resource_configuration` supports the following:

* `common_resource_spec` - (Optional, List) Common node resource specifications. Each block represents one common node; block count drives `CommonCount`.
* `core_resource_spec` - (Optional, List) Core node resource specifications. Each block represents one core node; block count drives `CoreCount`.
* `master_resource_spec` - (Optional, List) Master node resource specifications. Each block represents one master node; block count drives `MasterCount`.
* `task_resource_spec` - (Optional, List) Task node resource specifications. Each block represents one task node; block count drives `TaskCount`.

The `common_resource_spec` object of `all_node_resource_spec` supports the following:

* `data_disk` - (Optional, List) Cloud data disk specifications.
* `instance_type` - (Optional, String) CVM instance type, e.g., `S6.2XLARGE32`, `SA4.8XLARGE64`.
* `local_data_disk` - (Optional, List) Local data disk specifications.
* `system_disk` - (Optional, List) System disk specifications.
* `tags` - (Optional, List) Tags to bind to the node.

The `core_resource_spec` object of `all_node_resource_spec` supports the following:

* `data_disk` - (Optional, List) Cloud data disk specifications.
* `instance_type` - (Optional, String) CVM instance type, e.g., `S6.2XLARGE32`, `SA4.8XLARGE64`.
* `local_data_disk` - (Optional, List) Local data disk specifications.
* `system_disk` - (Optional, List) System disk specifications.
* `tags` - (Optional, List) Tags to bind to the node.

The `data_disk` object of `common_resource_spec` supports the following:

* `count` - (Optional, Int) Number of disks.
* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `CLOUD_HSSD`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `data_disk` object of `core_resource_spec` supports the following:

* `count` - (Optional, Int) Number of disks.
* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `CLOUD_HSSD`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `data_disk` object of `master_resource_spec` supports the following:

* `count` - (Optional, Int) Number of disks.
* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `CLOUD_HSSD`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `data_disk` object of `task_resource_spec` supports the following:

* `count` - (Optional, Int) Number of disks.
* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `CLOUD_HSSD`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `depend_service` object supports the following:

* `instance_id` - (Optional, String) Shared component cluster instance ID.
* `service_name` - (Optional, String) Shared component name.

The `instance_charge_prepaid` object supports the following:

* `period` - (Optional, Int) Purchase duration in months. Valid values: 1-12, 24, 36, 48, 60.
* `renew_flag` - (Optional, Bool) Auto-renew flag. Default is false.

The `local_data_disk` object of `common_resource_spec` supports the following:

* `count` - (Optional, Int) Number of disks.
* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `CLOUD_HSSD`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `local_data_disk` object of `core_resource_spec` supports the following:

* `count` - (Optional, Int) Number of disks.
* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `CLOUD_HSSD`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `local_data_disk` object of `master_resource_spec` supports the following:

* `count` - (Optional, Int) Number of disks.
* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `CLOUD_HSSD`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `local_data_disk` object of `task_resource_spec` supports the following:

* `count` - (Optional, Int) Number of disks.
* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `CLOUD_HSSD`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `login_settings` object supports the following:

* `password` - (Optional, String) Login password, 8-16 characters, must contain uppercase, lowercase, digit and special character (supported special chars: `!@%^*`). First char cannot be special.
* `public_key_id` - (Optional, String) Public key ID for key-based login.

The `master_resource_spec` object of `all_node_resource_spec` supports the following:

* `data_disk` - (Optional, List) Cloud data disk specifications.
* `instance_type` - (Optional, String) CVM instance type, e.g., `S6.2XLARGE32`, `SA4.8XLARGE64`.
* `local_data_disk` - (Optional, List) Local data disk specifications.
* `system_disk` - (Optional, List) System disk specifications.
* `tags` - (Optional, List) Tags to bind to the node.

The `meta_db_info` object supports the following:

* `meta_data_jdbc_url` - (Optional, String) Custom MetaDB JDBC URL, e.g., `jdbc:mysql://10.10.10.10:3306/dbname`.
* `meta_data_pass` - (Optional, String) Custom MetaDB password.
* `meta_data_user` - (Optional, String) Custom MetaDB username.
* `meta_type` - (Optional, String) Hive shared metadata DB type. Valid values: `EMR_DEFAULT_META`, `EMR_EXIT_META`, `USER_CUSTOM_META`.
* `unify_meta_instance_id` - (Optional, String) EMR-MetaDB instance ID.

The `node_marks` object supports the following:

* `node_names` - (Optional, List) Node marker names.
* `node_type` - (Optional, String) Node type: `master`, `core`, `task`, `router`.
* `zone` - (Optional, String) Availability zone name.

The `placement` object of `zone_resource_configuration` supports the following:

* `project_id` - (Optional, Int) Project ID. Defaults to default project if omitted.
* `zone` - (Optional, String) Availability zone, e.g., `ap-guangzhou-7`.

The `scene_software_config` object supports the following:

* `software` - (Required, List) List of components with versions, e.g., `["hdfs-3.2.2", "yarn-3.2.2"]`.
* `scene_name` - (Optional, String) Scenario name, e.g., `Hadoop-Default`, `Hadoop-Kudu`, `Hadoop-Zookeeper`, `Hadoop-Presto`, `Hadoop-Hbase`.

The `script_bootstrap_action_config` object supports the following:

* `args` - (Optional, List) Script arguments, following standard Shell convention.
* `cos_file_name` - (Optional, String) Script file name.
* `cos_file_uri` - (Optional, String) COS URI of the script.
* `execution_moment` - (Optional, String) Execution timing. Valid values: `resourceAfter`, `clusterAfter`, `clusterBefore`.
* `remark` - (Optional, String) Remark.

The `system_disk` object of `common_resource_spec` supports the following:

* `count` - (Optional, Int) Number of disks.
* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `CLOUD_HSSD`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `system_disk` object of `core_resource_spec` supports the following:

* `count` - (Optional, Int) Number of disks.
* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `CLOUD_HSSD`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `system_disk` object of `master_resource_spec` supports the following:

* `count` - (Optional, Int) Number of disks.
* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `CLOUD_HSSD`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `system_disk` object of `task_resource_spec` supports the following:

* `count` - (Optional, Int) Number of disks.
* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `CLOUD_HSSD`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `tags` object of `common_resource_spec` supports the following:

* `tag_key` - (Optional, String) Tag key.
* `tag_value` - (Optional, String) Tag value.

The `tags` object of `core_resource_spec` supports the following:

* `tag_key` - (Optional, String) Tag key.
* `tag_value` - (Optional, String) Tag value.

The `tags` object of `master_resource_spec` supports the following:

* `tag_key` - (Optional, String) Tag key.
* `tag_value` - (Optional, String) Tag value.

The `tags` object of `task_resource_spec` supports the following:

* `tag_key` - (Optional, String) Tag key.
* `tag_value` - (Optional, String) Tag value.

The `tags` object supports the following:

* `tag_key` - (Optional, String) Tag key.
* `tag_value` - (Optional, String) Tag value.

The `task_resource_spec` object of `all_node_resource_spec` supports the following:

* `data_disk` - (Optional, List) Cloud data disk specifications.
* `instance_type` - (Optional, String) CVM instance type, e.g., `S6.2XLARGE32`, `SA4.8XLARGE64`.
* `local_data_disk` - (Optional, List) Local data disk specifications.
* `system_disk` - (Optional, List) System disk specifications.
* `tags` - (Optional, List) Tags to bind to the node.

The `virtual_private_cloud` object of `zone_resource_configuration` supports the following:

* `subnet_id` - (Optional, String) Subnet ID.
* `vpc_id` - (Optional, String) VPC ID.

The `zone_resource_configuration` object supports the following:

* `all_node_resource_spec` - (Optional, List) Resource specifications for all node roles.
* `placement` - (Optional, List) Zone and project placement.
* `virtual_private_cloud` - (Optional, List) VPC/Subnet information.
* `zone_tag` - (Optional, String) Zone tag. For single-AZ, leave empty. For multi-AZ: `master`, `standby`, `third-party`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cluster_id` - Cluster ID (same as the resource ID).
* `status` - Cluster status code. `2` indicates the cluster is running.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `1h0m`) Used when creating the resource.
* `read` - (Defaults to `20m`) Used when reading the resource.
* `delete` - (Defaults to `30m`) Used when deleting the resource.

## Import

EMR cluster (v2) can be imported using the instance ID, e.g.

```
terraform import tencentcloud_emr_cluster_v2.example emr-xxxxxxxx
```

