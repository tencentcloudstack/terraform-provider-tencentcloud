---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_readonly_instance"
sidebar_current: "docs-tencentcloud-resource-mysql_readonly_instance"
description: |-
  Provides a mysql instance resource to create read-only database instances.
---

# tencentcloud_mysql_readonly_instance

Provides a mysql instance resource to create read-only database instances.

~> **NOTE:** Read-only instances can be purchased only for two-node or three-node source instances on MySQL 5.6 or above with the InnoDB engine at a specification of 1 GB memory and 50 GB disk capacity or above.
~> **NOTE:** The terminate operation of read only mysql does NOT take effect immediately, maybe takes for several hours. so during that time, VPCs associated with that mysql instance can't be terminated also.

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cdb"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-mysql"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-mysql"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-mysql"
  description = "mysql test"
}

resource "tencentcloud_mysql_instance" "example" {
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "PassWord123"
  slave_deploy_mode = 0
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  slave_sync_mode   = 1
  instance_name     = "tf-example-mysql"
  mem_size          = 4000
  volume_size       = 200
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  intranet_port     = 3306
  security_groups   = [tencentcloud_security_group.security_group.id]

  tags = {
    name = "test"
  }

  parameters = {
    character_set_server = "UTF8"
    max_connections      = "1000"
  }
}

resource "tencentcloud_mysql_readonly_instance" "example" {
  master_instance_id = tencentcloud_mysql_instance.example.id
  instance_name      = "tf-example"
  mem_size           = 128000
  volume_size        = 255
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
  intranet_port      = 3306
  security_groups    = [tencentcloud_security_group.security_group.id]

  tags = {
    createBy = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, String) The name of a mysql instance.
* `master_instance_id` - (Required, String) Indicates the master instance ID of recovery instances.
* `mem_size` - (Required, Int) Memory size (in MB).
* `volume_size` - (Required, Int) Disk size (in GB).
* `auto_renew_flag` - (Optional, Int) Auto renew flag. NOTES: Only supported prepaid instance.
* `charge_type` - (Optional, String, ForceNew) Pay type of instance. Valid values:`PREPAID`, `POSTPAID`. Default is `POSTPAID`.
* `cluster_topology` - (Optional, List) Cluster Edition node topology configuration. Note: If you purchased a cluster edition instance, this parameter is required. You need to set the RW and RO node topology of the cluster edition instance. The RO node range is 1-5. Please set at least 1 RO node.
* `cpu` - (Optional, Int) CPU cores.
* `destroy_protect` - (Optional, String) Instance destroy protection status. Valid values: `on` (enable destroy protection), `off` (disable destroy protection).
* `device_type` - (Optional, String) Specify device type, available values:
	- `UNIVERSAL` (default): universal instance,
	- `EXCLUSIVE`: exclusive instance,
	- `BASIC_V2`: ONTKE single-node instance,
	- `CLOUD_NATIVE_CLUSTER`: cluster version standard type,
	- `CLOUD_NATIVE_CLUSTER_EXCLUSIVE`: cluster version enhanced type.
If it is not specified, it defaults to a universal instance.
* `disk_encryption` - (Optional, String, ForceNew) Whether to enable disk encryption for the CDB instance. Valid value: `on` (enable disk encryption); otherwise disk encryption is not enabled. Only cloud-disk edition instances support this feature. This parameter can only be set at creation time and cannot be modified.
* `disk_type` - (Optional, String, ForceNew) Disk Type: This parameter can be specified for Single-Node (Cloud Disk) or Cloud Disk Edition instances. `CLOUD_SSD` designates an SSD cloud disk; `CLOUD_HSSD` designates an Enhanced SSD cloud disk; and `CLOUD_PREMIUM` designates a High-Performance cloud disk. Note: The regions that support the disk types for Single-Node (Cloud Disk) and Cloud Disk Edition instances vary slightly; please refer to `Regions and Availability Zones` for specific support details.
* `fast_upgrade` - (Optional, Int) Specify whether to enable fast upgrade when upgrade instance spec, available value: `1` - enabled, `0` - disabled.
* `force_delete` - (Optional, Bool) Indicate whether to delete instance directly or not. Default is `false`. If set true, the instance will be deleted instead of staying recycle bin. Note: only works for `PREPAID` instance. When the main mysql instance set true, this para of the readonly mysql instance will not take effect.
* `intranet_port` - (Optional, Int) Public access port. Valid value ranges: [1024~65535]. The default value is `3306`.
* `master_region` - (Optional, String) The region information of the master instance. This field is required when purchasing a cross-region subscription.
* `param_template_id` - (Optional, Int) Specify parameter template id.
* `pay_type` - (Optional, Int, **Deprecated**) It has been deprecated from version 1.36.0. Please use `charge_type` instead. Pay type of instance. Valid values: `0`, `1`. `0`: prepaid, `1`: postpaid.
* `period` - (Optional, Int, **Deprecated**) It has been deprecated from version 1.36.0. Please use `prepaid_period` instead. Period of instance. NOTES: Only supported prepaid instance.
* `prepaid_period` - (Optional, Int) Period of instance. NOTES: Only supported prepaid instance.
* `ro_group_id` - (Optional, String) Read only group id. If rogroupId is empty, a new ro group is created by default. If it is not empty, the existing ro group is used. Cross-region query requires master instance permission.
* `security_groups` - (Optional, Set: [`String`]) Security groups to use.
* `slave_deploy_mode` - (Optional, Int) Availability zone deployment method. Available values: 0 - Single availability zone; 1 - Multiple availability zones.
* `subnet_id` - (Optional, String) Private network ID. If `vpc_id` is set, this value is required.
* `tags` - (Optional, Map) Instance tags.
* `vpc_id` - (Optional, String) ID of VPC, which can be modified once every 24 hours and can't be removed.
* `wait_switch` - (Optional, Int) Switch the method of accessing new instances, default is `0`. Supported values include: `0` - switch immediately, `1` - switch in time window.
* `zone` - (Optional, String) Zone information, this parameter defaults to, the system automatically selects an Availability Zone.

The `cluster_topology` object supports the following:

* `read_only_nodes` - (Optional, Set) RO Node Topology.
* `read_write_node` - (Optional, List) RW Node Topology.

The `read_only_nodes` object of `cluster_topology` supports the following:

* `is_random_zone` - (Optional, Bool) Whether to distribute in random availability zones. Enter `true` to specify a random availability zone. Otherwise, use the availability zone specified by Zone.
* `node_id` - (Optional, String) When upgrading a cluster instance, if you want to adjust the availability zone of a read-only node, you need to specify the node ID.
* `zone` - (Optional, String) Specifies the availability zone where the node is distributed.

The `read_write_node` object of `cluster_topology` supports the following:

* `zone` - (Required, String) The availability zone where the RW node is located.
* `node_id` - (Optional, String) When upgrading a cluster instance, if you want to adjust the availability zone of a read-only node, you need to specify the node ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `intranet_ip` - instance intranet IP.
* `locked` - Indicates whether the instance is locked. Valid values: `0`, `1`. `0` - No; `1` - Yes.
* `ro_vip` - VIP-only read access.
* `ro_vport` - VIP port number (read-only).
* `status` - Instance status. Valid values: `0`, `1`, `4`, `5`. `0` - Creating; `1` - Running; `4` - Isolating; `5` - Isolated.
* `task_status` - Indicates which kind of operations is being executed.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateDBInstance](https://cloud.tencent.com/document/api/236/15871) - This API is used to create a cloud database instance with an annual/monthly subscription, including primary instance, disaster recovery instance, and read-only instance.
* Create: [CreateDBInstanceHour](https://cloud.tencent.com/document/api/236/15865) - This API is used to create pay-as-you-go instances. You can create a cloud database instance by inputting the instance specification, MySQL version number, quantity, etc.
* Create: [DescribeDBInstances](https://cloud.tencent.com/document/api/236/15872) - This API is used to query the list of TencentDB for MySQL instances. It supports filtering instances by conditions such as project ID, instance ID, access address, and instance status.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Read: [DescribeDBInstances](https://cloud.tencent.com/document/api/236/15872) - This API is used to query the list of TencentDB for MySQL instances. It supports filtering instances by conditions such as project ID, instance ID, access address, and instance status.
* Read: [DescribeDBSecurityGroups](https://cloud.tencent.com/document/api/236/15854) - This API (DescribeDBSecurityGroups) is used to query the security group details of an instance.
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [DescribeAsyncRequestInfo](https://cloud.tencent.com/document/api/236/20410) - This API (DescribeAsyncRequestInfo) is used to query the async task execution result of a TencentDB instance.
* Update: [DescribeDBInstances](https://cloud.tencent.com/document/api/236/15872) - This API is used to query the list of TencentDB for MySQL instances. It supports filtering instances by conditions such as project ID, instance ID, access address, and instance status.
* Update: [DisassociateSecurityGroups](https://cloud.tencent.com/document/api/236/15851) - This API (DisassociateSecurityGroups) is used to unbind security groups from instances in batches.
* Update: [ModifyAutoRenewFlag](https://cloud.tencent.com/document/api/236/19652) - This API is used to modify the auto-renewal flag of a TencentDB instance.
* Update: [ModifyDBInstanceName](https://cloud.tencent.com/document/api/236/15877) - This API (ModifyDBInstanceName) is used to rename a TencentDB instance.
* Update: [ModifyDBInstanceSecurityGroups](https://cloud.tencent.com/document/api/236/15853) - This API (ModifyDBInstanceSecurityGroups) is used to modify the security groups bound to a TencentDB instance.
* Update: [ModifyDBInstanceVipVport](https://cloud.tencent.com/document/api/236/15867) - This API is used to modify the IP and port number of a cloud database instance. It can also perform basic network to VPC network and subnet change under VPC network.
* Update: [ModifyInstanceDestroyProtect](https://cloud.tencent.com/document/api/236/135817)
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Update: [UpgradeDBInstance](https://cloud.tencent.com/document/api/236/15876) - This API is used to upgrade or downgrade the configuration of a cloud database instance. Supported instance types include primary instance, disaster recovery instance and read-only instance.
* Update: [UpgradeDBInstanceEngineVersion](https://cloud.tencent.com/document/api/236/15870) - This API is used to upgrade the version of a cloud database instance. Supported instance types include primary instance, disaster recovery instance, and read-only instance.
* Delete: [DescribeDBInstances](https://cloud.tencent.com/document/api/236/15872) - This API is used to query the list of TencentDB for MySQL instances. It supports filtering instances by conditions such as project ID, instance ID, access address, and instance status.
* Delete: [IsolateDBInstance](https://cloud.tencent.com/document/api/236/15869) - This API is used to isolate a cloud database instance. After an instance is isolated, you cannot access the database via IP and port. The isolated instance can be started in the recycle bin.
* Delete: [OfflineIsolatedInstances](https://cloud.tencent.com/document/api/236/36197) - This api is used to deactivate cloud database instances in quarantined state now. The instance Status for row operations must be quarantined state, such as instances with Status value 5 queried through the query instance list api.

## Import

mysql read-only database instances can be imported using the id, e.g.
```
terraform import tencentcloud_mysql_readonly_instance.default cdb-dnqksd9f
```

