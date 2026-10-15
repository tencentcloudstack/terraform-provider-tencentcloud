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

* Create: [CreateDBInstance](https://cloud.tencent.com/document/api/236/15871) - 创建云数据库实例（包年包月）
* Create: [CreateDBInstanceHour](https://cloud.tencent.com/document/api/236/15865) - 创建云数据库实例（按量计费）
* Create: [DescribeDBInstances](https://cloud.tencent.com/document/api/236/15872) - 查询实例列表
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Read: [DescribeDBInstances](https://cloud.tencent.com/document/api/236/15872) - 查询实例列表
* Read: [DescribeDBSecurityGroups](https://cloud.tencent.com/document/api/236/15854) - 查询实例安全组信息
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - 批量查看资源关联的标签
* Update: [DescribeAsyncRequestInfo](https://cloud.tencent.com/document/api/236/20410) - 查询异步任务的执行结果
* Update: [DescribeDBInstances](https://cloud.tencent.com/document/api/236/15872) - 查询实例列表
* Update: [DisassociateSecurityGroups](https://cloud.tencent.com/document/api/236/15851) - 安全组批量解绑云资源
* Update: [ModifyAutoRenewFlag](https://cloud.tencent.com/document/api/236/19652) - 修改云数据库实例的自动续费标记
* Update: [ModifyDBInstanceName](https://cloud.tencent.com/document/api/236/15877) - 修改云数据库实例名
* Update: [ModifyDBInstanceSecurityGroups](https://cloud.tencent.com/document/api/236/15853) - 修改云数据库安全组
* Update: [ModifyDBInstanceVipVport](https://cloud.tencent.com/document/api/236/15867) - 修改云数据库实例的IP和端口号
* Update: [ModifyInstanceDestroyProtect](https://cloud.tencent.com/document/api/236/135817) - 更新实例销毁保护状态
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Update: [RenewDBInstance](https://cloud.tencent.com/document/api/236/30160) - 续费云数据库实例
* Update: [UpgradeDBInstance](https://cloud.tencent.com/document/api/236/15876) - 调整云数据库实例的配置
* Update: [UpgradeDBInstanceEngineVersion](https://cloud.tencent.com/document/api/236/15870) - 升级实例版本
* Delete: [DescribeDBInstances](https://cloud.tencent.com/document/api/236/15872) - 查询实例列表
* Delete: [IsolateDBInstance](https://cloud.tencent.com/document/api/236/15869) - 隔离云数据库实例
* Delete: [OfflineIsolatedInstances](https://cloud.tencent.com/document/api/236/36197) - 下线隔离状态的云数据库实例

## Import

CDB dr(disaster recovery) instancecan be imported using the id, e.g.

```
terraform import tencentcloud_mysql_dr_instance.example cdb-bcet7sdb
```

