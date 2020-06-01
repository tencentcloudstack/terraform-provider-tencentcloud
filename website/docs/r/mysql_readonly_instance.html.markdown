---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_readonly_instance"
sidebar_current: "docs-tencentcloud-resource-mysql_readonly_instance"
description: |-
  Provides a mysql instance resource to create read-only database instances.
---

# tencentcloud_mysql_readonly_instance

Provides a mysql instance resource to create read-only database instances.

~> **NOTE:** The terminate operation of read only mysql does NOT take effect immediately, maybe takes for several hours. so during that time, VPCs associated with that mysql instance can't be terminated also.

## Example Usage

```hcl
resource "tencentcloud_mysql_readonly_instance" "default" {
  master_instance_id = "cdb-dnqksd9f"
  instance_name      = "myTestMysql"
  mem_size           = 128000
  volume_size        = 255
  vpc_id             = "vpc-12mt3l31"
  subnet_id          = "subnet-9uivyb1g"
  intranet_port      = 3306
  security_groups    = ["sg-ot8eclwz"]

  tags = {
    name = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required) The name of a mysql instance.
* `master_instance_id` - (Required, ForceNew) Indicates the master instance ID of recovery instances.
* `mem_size` - (Required) Memory size (in MB).
* `volume_size` - (Required) Disk size (in GB).
* `auto_renew_flag` - (Optional) Auto renew flag. NOTES: Only supported prepaid instance.
* `force_delete` - (Optional) Indicate whether to delete instance directly or not. Default is false. If set true, the instance will be deleted instead of staying recycle bin. Note: only works for `PREPAID` instance. When the main mysql instance set true, this para of the readonly mysql instance will not take effect.
* `intranet_port` - (Optional) Public access port, rang form 1024 to 65535 and default value is 3306.
* `pay_type` - (Optional, ForceNew) Pay type of instance, 0: prepay, 1: postpay. NOTES: Only supported prepay instance.
* `period` - (Optional) Period of instance. NOTES: Only supported prepaid instance.
* `security_groups` - (Optional) Security groups to use.
* `subnet_id` - (Optional) Private network ID. If vpc_id is set, this value is required.
* `tags` - (Optional) Instance tags.
* `vpc_id` - (Optional) ID of VPC, which can be modified once every 24 hours and can't be removed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `intranet_ip` - instance intranet IP.
* `locked` - Indicates whether the instance is locked. 0 - No; 1 - Yes.
* `status` - Instance status. Available values: 0 - Creating; 1 - Running; 4 - Isolating; 5 - Isolated.
* `task_status` - Indicates which kind of operations is being executed.


