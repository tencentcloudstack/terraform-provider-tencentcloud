---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_instance"
sidebar_current: "docs-tencentcloud-resource-mysql_instance"
description: |-
  Provides a mysql instance resource to create master database instances.
---

# tencentcloud_mysql_instance

Provides a mysql instance resource to create master database instances.

~> **NOTE:** If this mysql has readonly instance, the terminate operation of the mysql does NOT take effect immediately, maybe takes for several hours. so during that time, VPCs associated with that mysql instance can't be terminated also.

## Example Usage

```hcl
resource "tencentcloud_mysql_instance" "default" {
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "********"
  slave_deploy_mode = 0
  first_slave_zone  = "ap-guangzhou-4"
  second_slave_zone = "ap-guangzhou-4"
  slave_sync_mode   = 1
  availability_zone = "ap-guangzhou-4"
  project_id        = 201901010001
  instance_name     = "myTestMysql"
  mem_size          = 128000
  volume_size       = 250
  vpc_id            = "vpc-12mt3l31"
  subnet_id         = "subnet-9uivyb1g"
  intranet_port     = 3306
  security_groups   = ["sg-ot8eclwz"]

  tags = {
    name = "test"
  }

  parameters = {
    character_set_server = "UTF8"
    max_connections      = "1000"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, String) The name of a mysql instance.
* `mem_size` - (Required, Int) Memory size (in MB).
* `volume_size` - (Required, Int) Disk size (in GB).
* `auto_renew_flag` - (Optional, Int) Auto renew flag. NOTES: Only supported prepaid instance.
* `availability_zone` - (Optional, String, ForceNew) Indicates which availability zone will be used.
* `charge_type` - (Optional, String, ForceNew) Pay type of instance. Valid values:`PREPAID`, `POSTPAID`. Default is `POSTPAID`.
* `cpu` - (Optional, Int) CPU cores.
* `device_type` - (Optional, String) Specify device type, available values: `UNIVERSAL` (default), `EXCLUSIVE`, `BASIC`.
* `engine_version` - (Optional, String) The version number of the database engine to use. Supported versions include 5.5/5.6/5.7/8.0, and default is 5.7. Upgrade the instance engine version to support 5.6/5.7 and switch immediately.
* `fast_upgrade` - (Optional, Int) Specify whether to enable fast upgrade when upgrade instance spec, available value: `1` - enabled, `0` - disabled.
* `first_slave_zone` - (Optional, String) Zone information about first slave instance.
* `force_delete` - (Optional, Bool) Indicate whether to delete instance directly or not. Default is `false`. If set true, the instance will be deleted instead of staying recycle bin. Note: only works for `PREPAID` instance. When the main mysql instance set true, this para of the readonly mysql instance will not take effect.
* `internet_service` - (Optional, Int) Indicates whether to enable the access to an instance from public network: 0 - No, 1 - Yes.
* `intranet_port` - (Optional, Int) Public access port. Valid value ranges: [1024~65535]. The default value is `3306`.
* `max_deay_time` - (Optional, Int) Latency threshold. Value range 1~10. Only need to fill in when upgrading kernel subversion and engine version.
* `param_template_id` - (Optional, Int) Specify parameter template id.
* `parameters` - (Optional, Map) List of parameters to use.
* `pay_type` - (Optional, Int, **Deprecated**) It has been deprecated from version 1.36.0. Please use `charge_type` instead. Pay type of instance. Valid values: `0`, `1`. `0`: prepaid, `1`: postpaid.
* `period` - (Optional, Int, **Deprecated**) It has been deprecated from version 1.36.0. Please use `prepaid_period` instead. Period of instance. NOTES: Only supported prepaid instance.
* `prepaid_period` - (Optional, Int) Period of instance. NOTES: Only supported prepaid instance.
* `project_id` - (Optional, Int) Project ID, default value is 0.
* `root_password` - (Optional, String) Password of root account. This parameter can be specified when you purchase master instances, but it should be ignored when you purchase read-only instances or disaster recovery instances.
* `second_slave_zone` - (Optional, String) Zone information about second slave instance.
* `security_groups` - (Optional, Set: [`String`]) Security groups to use.
* `slave_deploy_mode` - (Optional, Int) Availability zone deployment method. Available values: 0 - Single availability zone; 1 - Multiple availability zones.
* `slave_sync_mode` - (Optional, Int) Data replication mode. 0 - Async replication; 1 - Semisync replication; 2 - Strongsync replication.
* `subnet_id` - (Optional, String) Private network ID. If `vpc_id` is set, this value is required.
* `tags` - (Optional, Map) Instance tags.
* `upgrade_subversion` - (Optional, Int) Whether it is a kernel subversion upgrade, supported values: 1 - upgrade the kernel subversion; 0 - upgrade the database engine version. Only need to fill in when upgrading kernel subversion and engine version.
* `vpc_id` - (Optional, String) ID of VPC, which can be modified once every 24 hours and can't be removed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `gtid` - Indicates whether GTID is enable. `0` - Not enabled; `1` - Enabled.
* `internet_host` - host for public access.
* `internet_port` - Access port for public access.
* `intranet_ip` - instance intranet IP.
* `locked` - Indicates whether the instance is locked. Valid values: `0`, `1`. `0` - No; `1` - Yes.
* `status` - Instance status. Valid values: `0`, `1`, `4`, `5`. `0` - Creating; `1` - Running; `4` - Isolating; `5` - Isolated.
* `task_status` - Indicates which kind of operations is being executed.


## Import

MySQL instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_mysql_instance.foo cdb-12345678"
```

