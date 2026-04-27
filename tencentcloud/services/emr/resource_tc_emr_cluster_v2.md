Provides a resource to create a EMR cluster (v2) using the `CreateCluster` API.

~> **NOTE:** The current resource `tencentcloud_emr_cluster_v2` does not yet support in-place updates. Any field change after creation will be silently ignored; destroy and re-create the resource to apply modifications. Modify APIs (scale-out/scale-in, tag modification, rename, ...) will be added in a follow-up release.

~> **NOTE:** Node counts are driven implicitly by the number of `master_resource_spec`, `core_resource_spec`, `task_resource_spec`, and `common_resource_spec` blocks declared. Each block represents one node of that role, and the first block is used as the resource template passed to the `CreateCluster` API.

Example Usage

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
    software   = [
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

Prepaid cluster with auto-renew

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

Import

EMR cluster (v2) can be imported using the instance ID, e.g.

```
terraform import tencentcloud_emr_cluster_v2.example emr-xxxxxxxx
```
