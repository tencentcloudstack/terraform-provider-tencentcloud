---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_dr_instance"
sidebar_current: "docs-tencentcloud-resource-mysql_dr_instance"
description: |-
  Provides a mysql instance resource to create CDB dr(disaster recovery) instance.
---

# tencentcloud_mysql_dr_instance

Provides a mysql instance resource to create CDB dr(disaster recovery) instance.

~> **NOTE:** Field `charge_type` only supports modification from `POSTPAID` to `PREPAID`. And the default renewal period is 1 month. and you can also use the `prepaid_period` field to customize the renewal period.

## Example Usage

### Create POSTPAID dr instance

```hcl
resource "tencentcloud_mysql_dr_instance" "example" {
  master_instance_id = "cdb-3kwa3gfj"
  master_region      = "ap-guangzhou"
  auto_renew_flag    = 0
  availability_zone  = "ap-guangzhou-6"
  charge_type        = "POSTPAID"
  cpu                = 4
  device_type        = "UNIVERSAL"
  first_slave_zone   = "ap-guangzhou-7"
  instance_name      = "tf-example"
  mem_size           = 8000
  project_id         = 0
  security_groups = [
    "sg-e6a8xxib",
  ]
  slave_deploy_mode = 1
  slave_sync_mode   = 0
  subnet_id         = "subnet-hhi88a58"
  volume_size       = 100
  vpc_id            = "vpc-i5yyodl9"
  intranet_port     = 3360
  tags = {
    createBy = "Terraform"
  }
}
```

### Create PREPAID dr instance

```hcl
resource "tencentcloud_mysql_dr_instance" "example" {
  master_instance_id = "cdb-3kwa3gfj"
  master_region      = "ap-guangzhou"
  availability_zone  = "ap-guangzhou-6"
  charge_type        = "PREPAID"
  prepaid_period     = 1
  auto_renew_flag    = 1
  cpu                = 4
  device_type        = "UNIVERSAL"
  first_slave_zone   = "ap-guangzhou-7"
  instance_name      = "tf-example"
  mem_size           = 8000
  project_id         = 0
  security_groups = [
    "sg-e6a8xxib",
  ]
  slave_deploy_mode = 1
  slave_sync_mode   = 0
  subnet_id         = "subnet-hhi88a58"
  volume_size       = 100
  vpc_id            = "vpc-i5yyodl9"
  intranet_port     = 3360
  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, String) The name of a mysql instance.
* `master_instance_id` - (Required, String) Indicates the master instance ID of recovery instances.
* `master_region` - (Required, String) The zone information of the primary instance is required when you purchase a disaster recovery instance.
* `mem_size` - (Required, Int) Memory size (in MB).
* `volume_size` - (Required, Int) Disk size (in GB).
* `auto_renew_flag` - (Optional, Int) Auto renew flag. NOTES: Only supported prepaid instance.
* `availability_zone` - (Optional, String) Indicates which availability zone will be used.
* `charge_type` - (Optional, String, ForceNew) Pay type of instance. Valid values:`PREPAID`, `POSTPAID`. Default is `POSTPAID`.
* `cpu` - (Optional, Int) CPU cores.
* `device_type` - (Optional, String) Specify device type, available values: `UNIVERSAL` (default), `EXCLUSIVE`, `BASIC`.
* `disk_type` - (Optional, String, ForceNew) Disk Type: This parameter can be specified for Single-Node (Cloud Disk) or Cloud Disk Edition instances. `CLOUD_SSD` designates an SSD cloud disk; `CLOUD_HSSD` designates an Enhanced SSD cloud disk; and `CLOUD_PREMIUM` designates a High-Performance cloud disk. Note: The regions that support the disk types for Single-Node (Cloud Disk) and Cloud Disk Edition instances vary slightly; please refer to `Regions and Availability Zones` for specific support details.
* `first_slave_zone` - (Optional, String) Zone information about first slave instance.
* `force_delete` - (Optional, Bool) Indicate whether to delete instance directly or not. Default is `false`. If set true, the instance will be deleted instead of staying recycle bin. Note: only works for `PREPAID` instance.
* `intranet_port` - (Optional, Int) Public access port. Valid value ranges: [1024~65535]. The default value is `3306`.
* `pay_type` - (Optional, Int, **Deprecated**) It has been deprecated from version 1.36.0. Please use `charge_type` instead. Pay type of instance. Valid values: `0`, `1`. `0`: prepaid, `1`: postpaid.
* `period` - (Optional, Int, **Deprecated**) It has been deprecated from version 1.36.0. Please use `prepaid_period` instead. Period of instance. NOTES: Only supported prepaid instance.
* `prepaid_period` - (Optional, Int) Period of instance. NOTES: Only supported prepaid instance.
* `project_id` - (Optional, Int) Project ID, default value is 0.
* `second_slave_zone` - (Optional, String) Zone information about second slave instance.
* `security_groups` - (Optional, Set: [`String`]) Security groups to use.
* `slave_deploy_mode` - (Optional, Int) Availability zone deployment method. Available values: 0 - Single availability zone; 1 - Multiple availability zones.
* `slave_sync_mode` - (Optional, Int) Data replication mode. 0 - Async replication; 1 - Semisync replication; 2 - Strongsync replication.
* `subnet_id` - (Optional, String) Private network ID. If `vpc_id` is set, this value is required.
* `tags` - (Optional, Map) Instance tags.
* `vpc_id` - (Optional, String) ID of VPC, which can be modified once every 24 hours and can't be removed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `intranet_ip` - instance intranet IP.

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
* Update: [RenewDBInstance](https://cloud.tencent.com/document/api/236/30160) - This API is used to renew cloud database instances. It supports yearly/monthly subscription instances. Pay-as-you-go instances can be renewed as yearly/monthly subscription instances through this API.
* Update: [UpgradeDBInstance](https://cloud.tencent.com/document/api/236/15876) - This API is used to upgrade or downgrade the configuration of a cloud database instance. Supported instance types include primary instance, disaster recovery instance and read-only instance.
* Update: [UpgradeDBInstanceEngineVersion](https://cloud.tencent.com/document/api/236/15870) - This API is used to upgrade the version of a cloud database instance. Supported instance types include primary instance, disaster recovery instance, and read-only instance.
* Delete: [DescribeDBInstances](https://cloud.tencent.com/document/api/236/15872) - This API is used to query the list of TencentDB for MySQL instances. It supports filtering instances by conditions such as project ID, instance ID, access address, and instance status.
* Delete: [IsolateDBInstance](https://cloud.tencent.com/document/api/236/15869) - This API is used to isolate a cloud database instance. After an instance is isolated, you cannot access the database via IP and port. The isolated instance can be started in the recycle bin.
* Delete: [OfflineIsolatedInstances](https://cloud.tencent.com/document/api/236/36197) - This api is used to deactivate cloud database instances in quarantined state now. The instance Status for row operations must be quarantined state, such as instances with Status value 5 queried through the query instance list api.

## Import

CDB dr(disaster recovery) instancecan be imported using the id, e.g.

```
terraform import tencentcloud_mysql_dr_instance.example cdb-bcet7sdb
```

