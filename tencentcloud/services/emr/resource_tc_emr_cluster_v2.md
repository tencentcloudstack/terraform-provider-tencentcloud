Provides a resource to create a EMR cluster (v2).

~> **NOTE:** At create time, every block of the same role within a zone (i.e. all `master_resource_spec` blocks, all `core_resource_spec` blocks, etc.) must declare an identical configuration — including `instance_type`, `system_disk`, `data_disk`, and `software`. The EMR `CreateCluster` API only accepts a single resource template per role and provisions the requested count of identical nodes from it. To run heterogeneous configurations within the same role, first create the cluster with uniform blocks, then customize individual nodes via subsequent `terraform apply` updates (resize disks, change `instance_type`, add data disks, etc.).

Example Usage

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

Import

EMR cluster (v2) can be imported using the instance ID, e.g.

```
terraform import tencentcloud_emr_cluster_v2.example emr-dvacz2w2
```
