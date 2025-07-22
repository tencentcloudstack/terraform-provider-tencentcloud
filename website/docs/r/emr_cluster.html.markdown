---
subcategory: "MapReduce(EMR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_emr_cluster"
sidebar_current: "docs-tencentcloud-resource-emr_cluster"
description: |-
  Provide a resource to create an emr cluster.
---

# tencentcloud_emr_cluster

Provide a resource to create an emr cluster.

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_instance_types" "cvm4c8m" {
  exclude_sold_out = true
  cpu_core_count   = 4
  memory_size      = 8
  filter {
    name   = "instance-charge-type"
    values = ["POSTPAID_BY_HOUR"]
  }
  filter {
    name   = "zone"
    values = [var.availability_zone]
  }
}

resource "tencentcloud_vpc" "emr_vpc" {
  name       = "emr-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "emr_subnet" {
  availability_zone = var.availability_zone
  name              = "emr-subnets"
  vpc_id            = tencentcloud_vpc.emr_vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_security_group" "emr_sg" {
  name        = "emr-sg"
  description = "emr sg"
  project_id  = 0
}

resource "tencentcloud_emr_cluster" "emr_cluster" {
  product_id = 38
  vpc_settings = {
    vpc_id    = tencentcloud_vpc.emr_vpc.id
    subnet_id = tencentcloud_subnet.emr_subnet.id
  }
  softwares = [
    "hdfs-2.8.5",
    "knox-1.6.1",
    "openldap-2.4.44",
    "yarn-2.8.5",
    "zookeeper-3.6.3",
  ]
  support_ha    = 0
  instance_name = "emr-cluster-test"
  resource_spec {
    master_resource_spec {
      mem_size     = 8192
      cpu          = 4
      disk_size    = 100
      disk_type    = "CLOUD_PREMIUM"
      spec         = "CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
      storage_type = 5
      root_size    = 50
    }
    core_resource_spec {
      mem_size     = 8192
      cpu          = 4
      disk_size    = 100
      disk_type    = "CLOUD_PREMIUM"
      spec         = "CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
      storage_type = 5
      root_size    = 50
    }
    master_count = 1
    core_count   = 2
  }
  login_settings = {
    password = "Tencent@cloud123"
  }
  time_span = 3600
  time_unit = "s"
  pay_mode  = 0
  placement_info {
    zone       = var.availability_zone
    project_id = 0
  }
  sg_id = tencentcloud_security_group.emr_sg.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, String, ForceNew) Name of the instance, which can contain 6 to 36 English letters, Chinese characters, digits, dashes(-), or underscores(_).
* `pay_mode` - (Required, Int) The pay mode of instance. 0 represent POSTPAID_BY_HOUR, 1 represent PREPAID.
* `product_id` - (Required, Int, ForceNew) Product ID. Different products ID represents different EMR product versions. Value range:
	- 16: represents EMR-V2.3.0
	- 20: represents EMR-V2.5.0
	- 25: represents EMR-V3.1.0
	- 27: represents KAFKA-V1.0.0
	- 30: represents EMR-V2.6.0
	- 33: represents EMR-V3.2.1
	- 34: represents EMR-V3.3.0
	- 37: represents EMR-V3.4.0
	- 38: represents EMR-V2.7.0
	- 44: represents EMR-V3.5.0
	- 50: represents KAFKA-V2.0.0
	- 51: represents STARROCKS-V1.4.0
	- 53: represents EMR-V3.6.0
	- 54: represents STARROCKS-V2.0.0.
* `softwares` - (Required, Set: [`String`], ForceNew) The softwares of a EMR instance.
* `support_ha` - (Required, Int, ForceNew) The flag whether the instance support high availability.(0=>not support, 1=>support).
* `vpc_settings` - (Required, Map, ForceNew) The private net config of EMR instance.
* `auto_renew` - (Optional, Int) 0 means turn off automatic renewal, 1 means turn on automatic renewal. Default is 0.
* `display_strategy` - (Optional, String, **Deprecated**) It will be deprecated in later versions. Display strategy of EMR instance.
* `extend_fs_field` - (Optional, String) Access the external file system.
* `login_settings` - (Optional, Map) Instance login settings. There are two optional fields:- password: Instance login password: 8-16 characters, including uppercase letters, lowercase letters, numbers and special characters. Special symbols only support! @% ^ *. The first bit of the password cannot be a special character;- public_key_id: Public key id. After the key is associated, the instance can be accessed through the corresponding private key.
* `multi_zone_setting` - (Optional, List) The specification of node resources is as follows: fill in a few available areas. In order, the first one is the main available area, the second one is the backup available area, and the third one is the arbitration available area.
* `multi_zone` - (Optional, Bool, ForceNew) true means that cross-AZ deployment is enabled; it is only a user parameter when creating a new cluster, and no subsequent adjustment is supported.
* `need_master_wan` - (Optional, String, ForceNew) Whether to enable the cluster Master node public network. Value range:
				- NEED_MASTER_WAN: Indicates that the cluster Master node public network is enabled.
				- NOT_NEED_MASTER_WAN: Indicates that it is not turned on.
				By default, the cluster Master node internet is enabled.
* `placement_info` - (Optional, List) The location of the instance.
* `placement` - (Optional, Map, **Deprecated**) It will be deprecated in later versions. Use `placement_info` instead. The location of the instance.
* `pre_executed_file_settings` - (Optional, List, ForceNew) Pre executed file settings. It can only be set at the time of creation, and cannot be modified.
* `resource_spec` - (Optional, List) Resource specification of EMR instance.
* `scene_name` - (Optional, String) Scene-based value:
	- Hadoop-Kudu
	- Hadoop-Zookeeper
	- Hadoop-Presto
	- Hadoop-Hbase.
* `sg_id` - (Optional, String, ForceNew) The ID of the security group to which the instance belongs, in the form of sg-xxxxxxxx.
* `tags` - (Optional, Map) Tag description list.
* `terminate_node_info` - (Optional, List) Terminate nodes. Note: it only works when the number of nodes decreases.
* `time_span` - (Optional, Int) The length of time the instance was purchased. Use with TimeUnit.When TimeUnit is s, the parameter can only be filled in at 3600, representing a metered instance.
When TimeUnit is m, the number filled in by this parameter indicates the length of purchase of the monthly instance of the package year, such as 1 for one month of purchase.
* `time_unit` - (Optional, String) The unit of time in which the instance was purchased. When PayMode is 0, TimeUnit can only take values of s(second). When PayMode is 1, TimeUnit can only take the value m(month).

The `common_resource_spec` object of `resource_spec` supports the following:

* `cpu` - (Optional, Int, ForceNew) Number of CPU cores.
* `disk_size` - (Optional, Int, ForceNew) Data disk capacity.
* `disk_type` - (Optional, String, ForceNew) disk types. Value range:
	- CLOUD_SSD: Represents cloud SSD;
	- CLOUD_PREMIUM: Represents efficient cloud disk;
	- CLOUD_BASIC: Represents Cloud Block Storage.
* `mem_size` - (Optional, Int, ForceNew) Memory size in M.
* `multi_disks` - (Optional, Set, ForceNew) Cloud disk list. When the data disk is a cloud disk, use disk_type and disk_size parameters directly, and use multi_disks for excess parts.
* `root_size` - (Optional, Int, ForceNew) Root disk capacity.
* `spec` - (Optional, String, ForceNew) Node specification description, such as CVM.SA2.
* `storage_type` - (Optional, Int, ForceNew) Storage type. Value range:
	- 4: Represents cloud SSD;
	- 5: Represents efficient cloud disk;
	- 6: Represents enhanced SSD Cloud Block Storage;
	- 11: Represents throughput Cloud Block Storage;
	- 12: Represents extremely fast SSD Cloud Block Storage.

The `core_resource_spec` object of `resource_spec` supports the following:

* `cpu` - (Optional, Int, ForceNew) Number of CPU cores.
* `disk_size` - (Optional, Int, ForceNew) Data disk capacity.
* `disk_type` - (Optional, String, ForceNew) disk types. Value range:
	- CLOUD_SSD: Represents cloud SSD;
	- CLOUD_PREMIUM: Represents efficient cloud disk;
	- CLOUD_BASIC: Represents Cloud Block Storage.
* `mem_size` - (Optional, Int, ForceNew) Memory size in M.
* `multi_disks` - (Optional, Set, ForceNew) Cloud disk list. When the data disk is a cloud disk, use disk_type and disk_size parameters directly, and use multi_disks for excess parts.
* `root_size` - (Optional, Int, ForceNew) Root disk capacity.
* `spec` - (Optional, String, ForceNew) Node specification description, such as CVM.SA2.
* `storage_type` - (Optional, Int, ForceNew) Storage type. Value range:
	- 4: Represents cloud SSD;
	- 5: Represents efficient cloud disk;
	- 6: Represents enhanced SSD Cloud Block Storage;
	- 11: Represents throughput Cloud Block Storage;
	- 12: Represents extremely fast SSD Cloud Block Storage.

The `master_resource_spec` object of `resource_spec` supports the following:

* `cpu` - (Optional, Int, ForceNew) Number of CPU cores.
* `disk_size` - (Optional, Int, ForceNew) Data disk capacity.
* `disk_type` - (Optional, String, ForceNew) disk types. Value range:
	- CLOUD_SSD: Represents cloud SSD;
	- CLOUD_PREMIUM: Represents efficient cloud disk;
	- CLOUD_BASIC: Represents Cloud Block Storage.
* `mem_size` - (Optional, Int, ForceNew) Memory size in M.
* `multi_disks` - (Optional, Set, ForceNew) Cloud disk list. When the data disk is a cloud disk, use disk_type and disk_size parameters directly, and use multi_disks for excess parts.
* `root_size` - (Optional, Int, ForceNew) Root disk capacity.
* `spec` - (Optional, String, ForceNew) Node specification description, such as CVM.SA2.
* `storage_type` - (Optional, Int, ForceNew) Storage type. Value range:
	- 4: Represents cloud SSD;
	- 5: Represents efficient cloud disk;
	- 6: Represents enhanced SSD Cloud Block Storage;
	- 11: Represents throughput Cloud Block Storage;
	- 12: Represents extremely fast SSD Cloud Block Storage.

The `multi_disks` object of `common_resource_spec` supports the following:

* `count` - (Optional, Int, ForceNew) Number of cloud disks of this type.
* `disk_type` - (Optional, String, ForceNew) Cloud disk type
	- CLOUD_SSD: Represents cloud SSD;
	- CLOUD_PREMIUM: Represents efficient cloud disk;
	- CLOUD_HSSD: Represents enhanced SSD Cloud Block Storage.
* `volume` - (Optional, Int, ForceNew) Cloud disk size.

The `multi_disks` object of `core_resource_spec` supports the following:

* `count` - (Optional, Int, ForceNew) Number of cloud disks of this type.
* `disk_type` - (Optional, String, ForceNew) Cloud disk type
	- CLOUD_SSD: Represents cloud SSD;
	- CLOUD_PREMIUM: Represents efficient cloud disk;
	- CLOUD_HSSD: Represents enhanced SSD Cloud Block Storage.
* `volume` - (Optional, Int, ForceNew) Cloud disk size.

The `multi_disks` object of `master_resource_spec` supports the following:

* `count` - (Optional, Int, ForceNew) Number of cloud disks of this type.
* `disk_type` - (Optional, String, ForceNew) Cloud disk type
	- CLOUD_SSD: Represents cloud SSD;
	- CLOUD_PREMIUM: Represents efficient cloud disk;
	- CLOUD_HSSD: Represents enhanced SSD Cloud Block Storage.
* `volume` - (Optional, Int, ForceNew) Cloud disk size.

The `multi_disks` object of `task_resource_spec` supports the following:

* `count` - (Optional, Int, ForceNew) Number of cloud disks of this type.
* `disk_type` - (Optional, String, ForceNew) Cloud disk type
	- CLOUD_SSD: Represents cloud SSD;
	- CLOUD_PREMIUM: Represents efficient cloud disk;
	- CLOUD_HSSD: Represents enhanced SSD Cloud Block Storage.
* `volume` - (Optional, Int, ForceNew) Cloud disk size.

The `multi_zone_setting` object supports the following:

* `vpc_settings` - (Required, Map, ForceNew) The private net config of EMR instance.
* `placement` - (Optional, List) The location of the instance.
* `resource_spec` - (Optional, List) Resource specification of EMR instance.

The `placement_info` object supports the following:

* `zone` - (Required, String) Zone.
* `project_id` - (Optional, Int) Project id.

The `placement` object of `multi_zone_setting` supports the following:

* `zone` - (Required, String, ForceNew) Zone.

The `pre_executed_file_settings` object supports the following:

* `args` - (Optional, List, ForceNew) Execution script parameters.
* `cos_file_name` - (Optional, String, ForceNew) Script file name.
* `cos_file_uri` - (Optional, String, ForceNew) The cos address of the script.
* `cos_secret_id` - (Optional, String, ForceNew) Cos secretId.
* `cos_secret_key` - (Optional, String, ForceNew) Cos secretKey.
* `remark` - (Optional, String, ForceNew) Remark.
* `run_order` - (Optional, Int, ForceNew) Run order.
* `when_run` - (Optional, String, ForceNew) `resourceAfter` or `clusterAfter`.

The `resource_spec` object of `multi_zone_setting` supports the following:

* `common_count` - (Optional, Int, ForceNew) The number of common node.
* `common_resource_spec` - (Optional, List, ForceNew) Resource details.
* `core_count` - (Optional, Int) The number of core node.
* `core_resource_spec` - (Optional, List, ForceNew) Resource details.
* `master_count` - (Optional, Int) The number of master node.
* `master_resource_spec` - (Optional, List, ForceNew) Resource details.
* `task_count` - (Optional, Int) The number of core node.
* `task_resource_spec` - (Optional, List, ForceNew) Resource details.

The `resource_spec` object supports the following:

* `common_count` - (Optional, Int, ForceNew) The number of common node.
* `common_resource_spec` - (Optional, List, ForceNew) Resource details.
* `core_count` - (Optional, Int) The number of core node.
* `core_resource_spec` - (Optional, List, ForceNew) Resource details.
* `master_count` - (Optional, Int) The number of master node.
* `master_resource_spec` - (Optional, List, ForceNew) Resource details.
* `task_count` - (Optional, Int) The number of core node.
* `task_resource_spec` - (Optional, List, ForceNew) Resource details.

The `task_resource_spec` object of `resource_spec` supports the following:

* `cpu` - (Optional, Int, ForceNew) Number of CPU cores.
* `disk_size` - (Optional, Int, ForceNew) Data disk capacity.
* `disk_type` - (Optional, String, ForceNew) disk types. Value range:
	- CLOUD_SSD: Represents cloud SSD;
	- CLOUD_PREMIUM: Represents efficient cloud disk;
	- CLOUD_BASIC: Represents Cloud Block Storage.
* `mem_size` - (Optional, Int, ForceNew) Memory size in M.
* `multi_disks` - (Optional, Set, ForceNew) Cloud disk list. When the data disk is a cloud disk, use disk_type and disk_size parameters directly, and use multi_disks for excess parts.
* `root_size` - (Optional, Int, ForceNew) Root disk capacity.
* `spec` - (Optional, String, ForceNew) Node specification description, such as CVM.SA2.
* `storage_type` - (Optional, Int, ForceNew) Storage type. Value range:
	- 4: Represents cloud SSD;
	- 5: Represents efficient cloud disk;
	- 6: Represents enhanced SSD Cloud Block Storage;
	- 11: Represents throughput Cloud Block Storage;
	- 12: Represents extremely fast SSD Cloud Block Storage.

The `terminate_node_info` object supports the following:

* `cvm_instance_ids` - (Optional, List) Destroy resource list.
* `node_flag` - (Optional, String) Value range of destruction node type: `MASTER`, `TASK`, `CORE`, `ROUTER`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_id` - Created EMR instance id.


