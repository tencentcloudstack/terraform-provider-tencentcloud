---
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
  internet_service = 1
  engine_version   = "5.7"

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
    max_connections = "1000"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required) The name of a mysql instance.
* `mem_size` - (Required) Memory size (in MB).
* `root_password` - (Required) Password of root account. This parameter can be specified when you purchase master instances, but it should be ignored when you purchase read-only instances or disaster recovery instances.
* `volume_size` - (Required) Disk size (in GB).
* `auto_renew_flag` - (Optional) Auto renew flag, works for prepay instance.
* `availability_zone` - (Optional, ForceNew) Indicates which availability zone will be used.
* `engine_version` - (Optional, ForceNew) The version number of the database engine to use. Supported versions include 5.5/5.6/5.7, and default is 5.7.
* `first_slave_zone` - (Optional, ForceNew) Zone information about first slave instance.
* `internet_service` - (Optional) Indicates whether to enable the access to an instance from public network: 0 - No, 1 - Yes.
* `intranet_port` - (Optional) Public access port, rang form 1024 to 65535 and default value is 3306.
* `parameters` - (Optional) List of parameters to use.
* `pay_type` - (Optional, ForceNew) Pay type of instance, 0: prepay, 1: postpay. Now only supported postpay.
* `period` - (Optional) Period of instance, works for prepay instance.
* `project_id` - (Optional) Project ID, default value is 0.
* `second_slave_zone` - (Optional, ForceNew) Zone information about second slave instance.
* `security_groups` - (Optional) Security groups to use.
* `slave_deploy_mode` - (Optional, ForceNew) Availability zone deployment method. Available values: 0 - Single availability zone; 1 - Multiple availability zones.
* `slave_sync_mode` - (Optional, ForceNew) Data replication mode. 0 - Async replication; 1 - Semisync replication; 2 - Strongsync replication.
* `subnet_id` - (Optional) Private network ID. If vpc_id is set, this value is required.
* `tags` - (Optional) Instance tags.
* `vpc_id` - (Optional) ID of VPC, which can be modified once every 24 hours and can't be removed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `gtid` - Indicates whether GTID is enable. 0 - Not enabled; 1 - Enabled.
* `internet_host` - host for public access.
* `internet_port` - Access port for public access.
* `intranet_ip` - instance intranet IP.
* `locked` - Indicates whether the instance is locked. 0 - No; 1 - Yes.
* `status` - Instance status. Available values: 0 - Creating; 1 - Running; 4 - Isolating; 5 - Isolated.
* `task_status` - Indicates which kind of operations is being executed.


