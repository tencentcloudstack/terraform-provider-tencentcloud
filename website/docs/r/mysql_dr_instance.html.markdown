---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_dr_instance"
sidebar_current: "docs-tencentcloud-resource-mysql_dr_instance"
description: |-
  Provides a mysql instance resource to create read-only database instances.
---

# tencentcloud_mysql_dr_instance

Provides a mysql instance resource to create read-only database instances.

## Example Usage

```hcl
resource "tencentcloud_mysql_dr_instance" "mysql_dr" {
  master_instance_id = "cdb-adjdu3t5"
  master_region      = "ap-guangzhou"
  auto_renew_flag    = 0
  availability_zone  = "ap-shanghai-3"
  charge_type        = "POSTPAID"
  cpu                = 4
  device_type        = "UNIVERSAL"
  first_slave_zone   = "ap-shanghai-4"
  instance_name      = "mysql-dr-test-up"
  mem_size           = 8000
  prepaid_period     = 1
  project_id         = 0
  security_groups = [
    "sg-q4d821qk",
  ]
  slave_deploy_mode = 1
  slave_sync_mode   = 0
  subnet_id         = "subnet-5vfntba5"
  volume_size       = 100
  vpc_id            = "vpc-h6s1s3aa"
  intranet_port     = 3360
  tags = {
    test = "test-tf"
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


## Import

mysql dr database instances can be imported using the id, e.g.
```
terraform import tencentcloud_mysql_dr_instance.mysql_dr cdb-bcet7sdb
```

