---
subcategory: "MapReduce(EMR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_emr_cluster_v2"
sidebar_current: "docs-tencentcloud-resource-emr_cluster_v2"
description: |-
  Provides a resource to create a EMR cluster (v2).
---

# tencentcloud_emr_cluster_v2

Provides a resource to create a EMR cluster (v2).

~> **NOTE:** At create time, every block of the same role within a zone (i.e. all `master_resource_spec` blocks, all `core_resource_spec` blocks, etc.) must declare an identical configuration — including `instance_type`, `system_disk`, `data_disk`, and `software`. The EMR `CreateCluster` API only accepts a single resource template per role and provisions the requested count of identical nodes from it. To run heterogeneous configurations within the same role, first create the cluster with uniform blocks, then customize individual nodes via subsequent `terraform apply` updates (resize disks, change `instance_type`, add data disks, etc.).

## Example Usage

```hcl
resource "tencentcloud_emr_cluster_v2" "example" {
  product_version        = "EMR-V3.6.0"
  enable_support_ha_flag = true
  instance_name          = "tf-example"
  instance_charge_type   = "POSTPAID_BY_HOUR"
  need_master_wan        = "NEED_MASTER_WAN"

  login_settings {
    password = "Password@123"
  }

  scene_software_config {
    scene_name = "Hadoop-Default"
  }

  meta_db_info {
    meta_type = "EMR_DEFAULT_META"
  }

  tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }

  security_group_ids = ["sg-37tigqat"]
  sg_ip              = "10.0.0.0/8"

  zone_resource_configuration {
    virtual_private_cloud {
      vpc_id    = "vpc-i5yyodl9"
      subnet_id = "subnet-hhi88a58"
    }

    placement {
      zone = "ap-guangzhou-6"
    }

    all_node_resource_spec {
      master_resource_spec {
        _node_index   = "master_node_1"
        instance_type = "SA5.2XLARGE16"
        system_disk {
          disk_size = 100
          disk_type = "CLOUD_HSSD"
        }

        data_disk {
          _disk_index = "master_disk_1"
          disk_size   = 100
          disk_type   = "CLOUD_SSD"
        }

        data_disk {
          _disk_index = "master_disk_2"
          disk_size   = 200
          disk_type   = "CLOUD_SSD"
        }

        software {
          services = "RUNTIME-1.0.0"
          roles    = ["Sysctl"]
        }

        software {
          services = "FILEBEAT-7.2.0"
          roles    = ["Filebeat"]
        }

        software {
          services = "HDFS-3.2.2"
          roles    = ["NameNode", "ZKFailoverController"]
        }

        software {
          services = "YARN-3.2.2"
          roles    = ["ResourceManager", "JobHistoryServer", "TimeLineServer"]
        }

        software {
          services = "OPENLDAP-2.4.44"
          roles    = ["slapd"]
        }

        software {
          services = "KNOX-1.6.1"
          roles    = ["gateway", "ldap"]
        }
      }

      master_resource_spec {
        _node_index   = "master_node_2"
        instance_type = "SA5.2XLARGE16"
        system_disk {
          disk_size = 100
          disk_type = "CLOUD_HSSD"
        }

        data_disk {
          _disk_index = "master_disk_1"
          disk_size   = 100
          disk_type   = "CLOUD_SSD"
        }

        data_disk {
          _disk_index = "master_disk_2"
          disk_size   = 200
          disk_type   = "CLOUD_SSD"
        }

        software {
          services = "RUNTIME-1.0.0"
          roles    = ["Sysctl"]
        }

        software {
          services = "FILEBEAT-7.2.0"
          roles    = ["Filebeat"]
        }

        software {
          services = "HDFS-3.2.2"
          roles    = ["NameNode", "ZKFailoverController"]
        }

        software {
          services = "YARN-3.2.2"
          roles    = ["ResourceManager", "JobHistoryServer", "TimeLineServer"]
        }

        software {
          services = "OPENLDAP-2.4.44"
          roles    = ["slapd"]
        }

        software {
          services = "KNOX-1.6.1"
          roles    = ["gateway", "ldap"]
        }
      }

      core_resource_spec {
        _node_index   = "core_node_1"
        instance_type = "SA5.LARGE8"
        system_disk {
          disk_size = 100
          disk_type = "CLOUD_HSSD"
        }

        data_disk {
          _disk_index = "core_disk_1"
          disk_size   = 100
          disk_type   = "CLOUD_SSD"
        }

        data_disk {
          _disk_index = "core_disk_2"
          disk_size   = 100
          disk_type   = "CLOUD_SSD"
        }

        software {
          services = "RUNTIME-1.0.0"
          roles    = ["Sysctl"]
        }

        software {
          services = "FILEBEAT-7.2.0"
          roles    = ["Filebeat"]
        }

        software {
          services = "HDFS-3.2.2"
          roles    = ["DataNode"]
        }

        software {
          services = "YARN-3.2.2"
          roles    = ["NodeManager"]
        }
      }

      core_resource_spec {
        _node_index   = "core_node_2"
        instance_type = "SA5.LARGE8"
        system_disk {
          disk_size = 100
          disk_type = "CLOUD_HSSD"
        }

        data_disk {
          _disk_index = "core_disk_1"
          disk_size   = 100
          disk_type   = "CLOUD_SSD"
        }

        data_disk {
          _disk_index = "core_disk_2"
          disk_size   = 100
          disk_type   = "CLOUD_SSD"
        }

        software {
          services = "RUNTIME-1.0.0"
          roles    = ["Sysctl"]
        }

        software {
          services = "FILEBEAT-7.2.0"
          roles    = ["Filebeat"]
        }

        software {
          services = "HDFS-3.2.2"
          roles    = ["DataNode"]
        }

        software {
          services = "YARN-3.2.2"
          roles    = ["NodeManager"]
        }
      }

      core_resource_spec {
        _node_index   = "core_node_3"
        instance_type = "SA5.LARGE8"
        system_disk {
          disk_size = 100
          disk_type = "CLOUD_HSSD"
        }

        data_disk {
          _disk_index = "core_disk_1"
          disk_size   = 100
          disk_type   = "CLOUD_SSD"
        }

        data_disk {
          _disk_index = "core_disk_2"
          disk_size   = 100
          disk_type   = "CLOUD_SSD"
        }

        software {
          services = "RUNTIME-1.0.0"
          roles    = ["Sysctl"]
        }

        software {
          services = "FILEBEAT-7.2.0"
          roles    = ["Filebeat"]
        }

        software {
          services = "HDFS-3.2.2"
          roles    = ["DataNode"]
        }

        software {
          services = "YARN-3.2.2"
          roles    = ["NodeManager"]
        }
      }

      task_resource_spec {
        _node_index   = "task_node_1"
        instance_type = "SA5.LARGE8"
        system_disk {
          disk_size = 100
          disk_type = "CLOUD_HSSD"
        }

        data_disk {
          _disk_index = "task_disk_1"
          disk_size   = 100
          disk_type   = "CLOUD_SSD"
        }

        data_disk {
          _disk_index = "task_disk_2"
          disk_size   = 100
          disk_type   = "CLOUD_SSD"
        }

        software {
          services = "RUNTIME-1.0.0"
          roles    = ["Sysctl"]
        }

        software {
          services = "FILEBEAT-7.2.0"
          roles    = ["Filebeat"]
        }

        software {
          services = "YARN-3.2.2"
          roles    = ["NodeManager"]
        }
      }

      task_resource_spec {
        _node_index   = "task_node_2"
        instance_type = "SA5.LARGE8"
        system_disk {
          disk_size = 100
          disk_type = "CLOUD_HSSD"
        }

        data_disk {
          _disk_index = "task_disk_1"
          disk_size   = 100
          disk_type   = "CLOUD_SSD"
        }

        data_disk {
          _disk_index = "task_disk_2"
          disk_size   = 100
          disk_type   = "CLOUD_SSD"
        }

        software {
          services = "RUNTIME-1.0.0"
          roles    = ["Sysctl"]
        }

        software {
          services = "FILEBEAT-7.2.0"
          roles    = ["Filebeat"]
        }

        software {
          services = "YARN-3.2.2"
          roles    = ["NodeManager"]
        }
      }

      common_resource_spec {
        _node_index   = "common_node_1"
        instance_type = "SA5.LARGE8"
        system_disk {
          disk_size = 100
          disk_type = "CLOUD_HSSD"
        }

        data_disk {
          _disk_index = "common_disk_1"
          disk_size   = 100
          disk_type   = "CLOUD_SSD"
        }

        software {
          services = "RUNTIME-1.0.0"
          roles    = ["Sysctl"]
        }

        software {
          services = "FILEBEAT-7.2.0"
          roles    = ["Filebeat"]
        }

        software {
          services = "HDFS-3.2.2"
          roles    = ["JournalNode"]
        }

        software {
          services = "ZOOKEEPER-3.6.3"
          roles    = ["Zookeeper"]
        }
      }

      common_resource_spec {
        _node_index   = "common_node_2"
        instance_type = "SA5.LARGE8"
        system_disk {
          disk_size = 100
          disk_type = "CLOUD_HSSD"
        }

        data_disk {
          _disk_index = "common_disk_1"
          disk_size   = 100
          disk_type   = "CLOUD_SSD"
        }

        software {
          services = "RUNTIME-1.0.0"
          roles    = ["Sysctl"]
        }

        software {
          services = "FILEBEAT-7.2.0"
          roles    = ["Filebeat"]
        }

        software {
          services = "HDFS-3.2.2"
          roles    = ["JournalNode"]
        }

        software {
          services = "ZOOKEEPER-3.6.3"
          roles    = ["Zookeeper"]
        }
      }

      common_resource_spec {
        _node_index   = "common_node_3"
        instance_type = "SA5.LARGE8"
        system_disk {
          disk_size = 100
          disk_type = "CLOUD_HSSD"
        }

        data_disk {
          _disk_index = "common_disk_1"
          disk_size   = 100
          disk_type   = "CLOUD_SSD"
        }

        software {
          services = "RUNTIME-1.0.0"
          roles    = ["Sysctl"]
        }

        software {
          services = "FILEBEAT-7.2.0"
          roles    = ["Filebeat"]
        }

        software {
          services = "HDFS-3.2.2"
          roles    = ["JournalNode"]
        }

        software {
          services = "ZOOKEEPER-3.6.3"
          roles    = ["Zookeeper"]
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `enable_support_ha_flag` - (Required, Bool, ForceNew) Whether to enable node high availability. `true` means enabled, `false` means disabled.
* `instance_charge_type` - (Required, String, ForceNew) Instance billing mode. Valid values: `PREPAID` (monthly/yearly subscription), `POSTPAID_BY_HOUR` (pay-as-you-go).
* `instance_name` - (Required, String) Instance name. Length 6-36 characters. Only Chinese, letters, numbers, `-`, `_` are allowed.
* `login_settings` - (Required, List, ForceNew) Login settings for purchased nodes.
* `product_version` - (Required, String, ForceNew) EMR product version name, e.g., `EMR-V3.5.0`.
* `scene_software_config` - (Required, List) Cluster scenario configuration. Components to deploy are now declared per node role via the `soft_ware` field inside each `*_resource_spec` block; the cluster-level `software` list below is computed (read-back from API) for observability.
* `cos_bucket` - (Optional, String, ForceNew) COS bucket path, used when creating StarRocks storage-compute separation clusters.
* `custom_conf` - (Optional, String, ForceNew) Custom software configuration in JSON format.
* `default_meta_version` - (Optional, String, ForceNew) Default metadata DB version. Valid values: `mysql8`, `tdsql8`, `mysql5`.
* `depend_service` - (Optional, List, ForceNew) Shared component dependency information.
* `disaster_recover_group_ids` - (Optional, List: [`String`], ForceNew) Spread placement group IDs. Currently supports only one ID.
* `enable_cbs_encrypt_flag` - (Optional, Bool, ForceNew) Whether to enable cluster-level CBS encryption. Default is false.
* `enable_cbs_sys_encrypt_flag` - (Optional, Bool, ForceNew) Whether to enable CBS system encryption.
* `enable_kerberos_flag` - (Optional, Bool, ForceNew) Whether to enable Kerberos authentication. Default is false.
* `enable_remote_login_flag` - (Optional, Bool, ForceNew) Whether to enable external remote login. Invalid when `security_group_ids` is set. Default is false.
* `instance_charge_prepaid` - (Optional, List, ForceNew) Prepaid (monthly/yearly) billing parameters. Required when `instance_charge_type` is `PREPAID`.
* `load_balancer_id` - (Optional, String, ForceNew) CLB instance ID, e.g., `lb-xxxxxxxx`.
* `meta_db_info` - (Optional, List, ForceNew) Metadata database information. When `meta_type` is `EMR_NEW_META`/`EMR_DEFAULT_META`, no extra fields are required; when `EMR_EXIT_META`, `unify_meta_instance_id` must be set; when `USER_CUSTOM_META`, `meta_data_jdbc_url`/`meta_data_user`/`meta_data_pass` must be set.
* `need_cdb_audit` - (Optional, Int, ForceNew) Whether to enable database auditing.
* `need_master_wan` - (Optional, String, ForceNew) Whether to enable master public network. Valid values: `NEED_MASTER_WAN` (default), `NOT_NEED_MASTER_WAN`.
* `partition_number` - (Optional, Int, ForceNew) Partition placement group partition number.
* `script_bootstrap_action_config` - (Optional, List, ForceNew) Bootstrap script configurations.
* `security_group_ids` - (Optional, List: [`String`], ForceNew) Security group IDs bound to the instance, e.g., `["sg-xxxxxxxx"]`.
* `sg_ip` - (Optional, String, ForceNew) Security source IP, e.g., `10.0.0.0/8`.
* `tags` - (Optional, List) Tags to bind to the cluster instance.
* `web_ui_version` - (Optional, Int, ForceNew) Service UI address version. `0`: single URL (default); `1`: all URLs.
* `zone_resource_configuration` - (Optional, List) Per-zone resource configuration. Supports 1 (single-AZ) or up to 3 entries (multi-AZ: primary, backup, arbitration). ZoneTag is derived automatically from the list index.

The `all_node_resource_spec` object of `zone_resource_configuration` supports the following:

* `common_resource_spec` - (Optional, Set) Common node resource specifications. Number of blocks = `CommonCount`. All blocks must have identical configuration; the first block is the single resource template sent to the API. This field is a `TypeSet` keyed by `_node_index` only  block order in HCL is irrelevant.
* `core_resource_spec` - (Optional, Set) Core node resource specifications. Number of blocks = `CoreCount`. All blocks must have identical configuration; the first block is the single resource template sent to the API. This field is a `TypeSet` keyed by `_node_index` only  block order in HCL is irrelevant.
* `master_resource_spec` - (Optional, Set) Master node resource specifications. Number of blocks = `MasterCount`. All blocks must have identical configuration; the first block is the single resource template sent to the API. This field is a `TypeSet` keyed by `_node_index` only  block order in HCL is irrelevant.
* `task_resource_spec` - (Optional, Set) Task node resource specifications. Number of blocks = `TaskCount`. All blocks must have identical configuration; the first block is the single resource template sent to the API. This field is a `TypeSet` keyed by `_node_index` only  block order in HCL is irrelevant.

The `common_resource_spec` object of `all_node_resource_spec` supports the following:

* `_node_index` - (Required, String) **Required** stable identity key for this node spec block. Must be unique within the same role's set, and must remain stable across plan/apply. Used by the Update handler to pair old/new blocks for in-place modification. Renaming an existing `_node_index` is rejected by CustomizeDiff for master/common (length immutable); for core/task it is interpreted as scale-in old + scale-out new.
* `software` - (Required, Set) Per-role software components (with their role/process lists) deployed on this node role. Must be identical across every block of the same role at create time. Aggregated across all four roles (deduped by `services`) and passed to `CreateCluster` as `SceneSoftwareConfig.Software`. Immutable after create modification is rejected at plan time by CustomizeDiff.
* `data_disk` - (Optional, Set) Cloud data disk specifications. `TypeSet` keyed by full content (including `_disk_index`); block order in HCL is irrelevant.
* `instance_type` - (Optional, String) CVM instance type, e.g., `S6.2XLARGE32`, `SA4.8XLARGE64`.
* `system_disk` - (Optional, List) System disk specifications.

The `core_resource_spec` object of `all_node_resource_spec` supports the following:

* `_node_index` - (Required, String) **Required** stable identity key for this node spec block. Must be unique within the same role's set, and must remain stable across plan/apply. Used by the Update handler to pair old/new blocks for in-place modification. Renaming an existing `_node_index` is rejected by CustomizeDiff for master/common (length immutable); for core/task it is interpreted as scale-in old + scale-out new.
* `software` - (Required, Set) Per-role software components (with their role/process lists) deployed on this node role. Must be identical across every block of the same role at create time. Aggregated across all four roles (deduped by `services`) and passed to `CreateCluster` as `SceneSoftwareConfig.Software`. Immutable after create modification is rejected at plan time by CustomizeDiff.
* `data_disk` - (Optional, Set) Cloud data disk specifications. `TypeSet` keyed by full content (including `_disk_index`); block order in HCL is irrelevant.
* `instance_type` - (Optional, String) CVM instance type, e.g., `S6.2XLARGE32`, `SA4.8XLARGE64`.
* `system_disk` - (Optional, List) System disk specifications.

The `data_disk` object of `common_resource_spec` supports the following:

* `_disk_index` - (Required, String) **Required** stable identity key for this `data_disk` block. Must be unique within the same node's `data_disk` set, and must remain stable across plan/apply once written to state. Renaming an existing `_disk_index` is rejected (treated as remove+add, which violates the no-shrink rule).
* `disk_size` - (Optional, Int) Disk size in GB. Can only be increased after creation; shrinking is rejected at plan time.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_HSSD`, `CLOUD_THROUGHPUT`, `CLOUD_TSSD`, `CLOUD_BIGDATA`, `CLOUD_HIGHIO`, `CLOUD_BSSD`, `REMOTE_SSD`. Immutable after creation.

The `data_disk` object of `core_resource_spec` supports the following:

* `_disk_index` - (Required, String) **Required** stable identity key for this `data_disk` block. Must be unique within the same node's `data_disk` set, and must remain stable across plan/apply once written to state. Renaming an existing `_disk_index` is rejected (treated as remove+add, which violates the no-shrink rule).
* `disk_size` - (Optional, Int) Disk size in GB. Can only be increased after creation; shrinking is rejected at plan time.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_HSSD`, `CLOUD_THROUGHPUT`, `CLOUD_TSSD`, `CLOUD_BIGDATA`, `CLOUD_HIGHIO`, `CLOUD_BSSD`, `REMOTE_SSD`. Immutable after creation.

The `data_disk` object of `master_resource_spec` supports the following:

* `_disk_index` - (Required, String) **Required** stable identity key for this `data_disk` block. Must be unique within the same node's `data_disk` set, and must remain stable across plan/apply once written to state. Renaming an existing `_disk_index` is rejected (treated as remove+add, which violates the no-shrink rule).
* `disk_size` - (Optional, Int) Disk size in GB. Can only be increased after creation; shrinking is rejected at plan time.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_HSSD`, `CLOUD_THROUGHPUT`, `CLOUD_TSSD`, `CLOUD_BIGDATA`, `CLOUD_HIGHIO`, `CLOUD_BSSD`, `REMOTE_SSD`. Immutable after creation.

The `data_disk` object of `task_resource_spec` supports the following:

* `_disk_index` - (Required, String) **Required** stable identity key for this `data_disk` block. Must be unique within the same node's `data_disk` set, and must remain stable across plan/apply once written to state. Renaming an existing `_disk_index` is rejected (treated as remove+add, which violates the no-shrink rule).
* `disk_size` - (Optional, Int) Disk size in GB. Can only be increased after creation; shrinking is rejected at plan time.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_HSSD`, `CLOUD_THROUGHPUT`, `CLOUD_TSSD`, `CLOUD_BIGDATA`, `CLOUD_HIGHIO`, `CLOUD_BSSD`, `REMOTE_SSD`. Immutable after creation.

The `depend_service` object supports the following:

* `instance_id` - (Optional, String, ForceNew) Shared component cluster instance ID.
* `service_name` - (Optional, String, ForceNew) Shared component name.

The `instance_charge_prepaid` object supports the following:

* `period` - (Optional, Int, ForceNew) Purchase duration in months. Valid values: 1-12, 24, 36, 48, 60.
* `renew_flag` - (Optional, Bool, ForceNew) Auto-renew flag. Default is false.

The `login_settings` object supports the following:

* `password` - (Optional, String, ForceNew) Login password, 8-16 characters, must contain uppercase, lowercase, digit and special character (supported special chars: `!@%^*`). First char cannot be special.
* `public_key_id` - (Optional, String, ForceNew) Public key ID for key-based login.

The `master_resource_spec` object of `all_node_resource_spec` supports the following:

* `_node_index` - (Required, String) **Required** stable identity key for this node spec block. Must be unique within the same role's set, and must remain stable across plan/apply. Used by the Update handler to pair old/new blocks for in-place modification. Renaming an existing `_node_index` is rejected by CustomizeDiff for master/common (length immutable); for core/task it is interpreted as scale-in old + scale-out new.
* `software` - (Required, Set) Per-role software components (with their role/process lists) deployed on this node role. Must be identical across every block of the same role at create time. Aggregated across all four roles (deduped by `services`) and passed to `CreateCluster` as `SceneSoftwareConfig.Software`. Immutable after create modification is rejected at plan time by CustomizeDiff.
* `data_disk` - (Optional, Set) Cloud data disk specifications. `TypeSet` keyed by full content (including `_disk_index`); block order in HCL is irrelevant.
* `instance_type` - (Optional, String) CVM instance type, e.g., `S6.2XLARGE32`, `SA4.8XLARGE64`.
* `system_disk` - (Optional, List) System disk specifications.

The `meta_db_info` object supports the following:

* `meta_data_jdbc_url` - (Optional, String, ForceNew) Custom MetaDB JDBC URL, e.g., `jdbc:mysql://10.10.10.10:3306/dbname`.
* `meta_data_pass` - (Optional, String, ForceNew) Custom MetaDB password.
* `meta_data_user` - (Optional, String, ForceNew) Custom MetaDB username.
* `meta_type` - (Optional, String, ForceNew) Hive shared metadata DB type. Valid values: `EMR_DEFAULT_META`, `EMR_EXIT_META`, `USER_CUSTOM_META`.
* `unify_meta_instance_id` - (Optional, String, ForceNew) EMR-MetaDB instance ID.

The `placement` object of `zone_resource_configuration` supports the following:

* `project_id` - (Optional, Int, ForceNew) Project ID. Defaults to default project if omitted.
* `zone` - (Optional, String, ForceNew) Availability zone, e.g., `ap-guangzhou-7`.

The `scene_software_config` object supports the following:

* `scene_name` - (Optional, String, ForceNew) Scenario name, e.g., `Hadoop-Default`, `Hadoop-Kudu`, `Hadoop-Zookeeper`, `Hadoop-Presto`, `Hadoop-Hbase`.

The `script_bootstrap_action_config` object supports the following:

* `args` - (Optional, List, ForceNew) Script arguments, following standard Shell convention.
* `cos_file_name` - (Optional, String, ForceNew) Script file name.
* `cos_file_uri` - (Optional, String, ForceNew) COS URI of the script.
* `execution_moment` - (Optional, String, ForceNew) Execution timing. Valid values: `resourceAfter`, `clusterAfter`, `clusterBefore`.
* `remark` - (Optional, String, ForceNew) Remark.

The `software` object of `common_resource_spec` supports the following:

* `roles` - (Required, Set) Process list for this component on this role, e.g., `["NameNode", "ZKFailoverController"]` for hdfs.
* `services` - (Required, String) Component name with version, e.g., `hdfs-3.2.2`.

The `software` object of `core_resource_spec` supports the following:

* `roles` - (Required, Set) Process list for this component on this role, e.g., `["NameNode", "ZKFailoverController"]` for hdfs.
* `services` - (Required, String) Component name with version, e.g., `hdfs-3.2.2`.

The `software` object of `master_resource_spec` supports the following:

* `roles` - (Required, Set) Process list for this component on this role, e.g., `["NameNode", "ZKFailoverController"]` for hdfs.
* `services` - (Required, String) Component name with version, e.g., `hdfs-3.2.2`.

The `software` object of `task_resource_spec` supports the following:

* `roles` - (Required, Set) Process list for this component on this role, e.g., `["NameNode", "ZKFailoverController"]` for hdfs.
* `services` - (Required, String) Component name with version, e.g., `hdfs-3.2.2`.

The `system_disk` object of `common_resource_spec` supports the following:

* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `system_disk` object of `core_resource_spec` supports the following:

* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `system_disk` object of `master_resource_spec` supports the following:

* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `system_disk` object of `task_resource_spec` supports the following:

* `disk_size` - (Optional, Int) Disk size in GB.
* `disk_type` - (Optional, String) Disk type. Valid values: `CLOUD_SSD`, `CLOUD_PREMIUM`, `CLOUD_BASIC`, `LOCAL_BASIC`, `LOCAL_SSD`.

The `tags` object supports the following:

* `tag_key` - (Optional, String) Tag key.
* `tag_value` - (Optional, String) Tag value.

The `task_resource_spec` object of `all_node_resource_spec` supports the following:

* `_node_index` - (Required, String) **Required** stable identity key for this node spec block. Must be unique within the same role's set, and must remain stable across plan/apply. Used by the Update handler to pair old/new blocks for in-place modification. Renaming an existing `_node_index` is rejected by CustomizeDiff for master/common (length immutable); for core/task it is interpreted as scale-in old + scale-out new.
* `software` - (Required, Set) Per-role software components (with their role/process lists) deployed on this node role. Must be identical across every block of the same role at create time. Aggregated across all four roles (deduped by `services`) and passed to `CreateCluster` as `SceneSoftwareConfig.Software`. Immutable after create modification is rejected at plan time by CustomizeDiff.
* `data_disk` - (Optional, Set) Cloud data disk specifications. `TypeSet` keyed by full content (including `_disk_index`); block order in HCL is irrelevant.
* `instance_type` - (Optional, String) CVM instance type, e.g., `S6.2XLARGE32`, `SA4.8XLARGE64`.
* `system_disk` - (Optional, List) System disk specifications.

The `virtual_private_cloud` object of `zone_resource_configuration` supports the following:

* `subnet_id` - (Optional, String, ForceNew) Subnet ID.
* `vpc_id` - (Optional, String, ForceNew) VPC ID.

The `zone_resource_configuration` object supports the following:

* `all_node_resource_spec` - (Optional, List) Resource specifications for all node roles.
* `placement` - (Optional, List, ForceNew) Zone and project placement.
* `virtual_private_cloud` - (Optional, List, ForceNew) VPC/Subnet information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cluster_id` - Cluster ID (same as the resource ID).
* `status` - Cluster status code. `2` indicates the cluster is running.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `30m`) Used when creating the resource.
* `update` - (Defaults to `30m`) Used when updating the resource.
* `delete` - (Defaults to `30m`) Used when deleting the resource.

## Import

EMR cluster (v2) can be imported using the instance ID, e.g.

```
terraform import tencentcloud_emr_cluster_v2.example emr-dvacz2w2
```

