---
subcategory: "CVM Dedicated Host(CDH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdh_instance"
sidebar_current: "docs-tencentcloud-resource-cdh_instance"
description: |-
  Provides a resource to manage CDH instance.
---

# tencentcloud_cdh_instance

Provides a resource to manage CDH instance.

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_cdh_instance" "foo" {
  availability_zone  = var.availability_zone
  host_type          = "HC20"
  charge_type        = "PREPAID"
  prepaid_period     = 1
  host_name          = "test"
  prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, String, ForceNew) The available zone for the CDH instance.
* `charge_type` - (Optional, String) The charge type of instance. Valid values are `PREPAID`. The default is `PREPAID`.
* `host_name` - (Optional, String) The name of the CDH instance. The max length of host_name is 60.
* `host_type` - (Optional, String, ForceNew) The type of the CDH instance.
* `prepaid_period` - (Optional, Int) The tenancy (time unit is month) of the prepaid instance, NOTE: it only works when charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.
* `prepaid_renew_flag` - (Optional, String) Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when charge_type is set to `PREPAID`.
* `project_id` - (Optional, Int) The project the instance belongs to, default to 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the instance.
* `cvm_instance_ids` - Id of CVM instances that have been created on the CDH instance.
* `expired_time` - Expired time of the instance.
* `host_resource` - An information list of host resource. Each element contains the following attributes:
  * `cpu_available_num` - The number of available CPU cores of the instance.
  * `cpu_total_num` - The number of total CPU cores of the instance.
  * `disk_available_size` - Instance disk available capacity, unit in GB.
  * `disk_total_size` - Instance disk total capacity, unit in GB.
  * `disk_type` - Type of the disk.
  * `memory_available_size` - Instance memory available capacity, unit in GB.
  * `memory_total_size` - Instance memory total capacity, unit in GB.
* `host_state` - State of the CDH instance.


## Import

CDH instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_cdh_instance.foo host-d6s7i5q4
```

