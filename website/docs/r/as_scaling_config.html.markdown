---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_scaling_config"
sidebar_current: "docs-tencentcloud-resource-as_scaling_config"
description: |-
  Provides a resource to create a configuration for an AS (Auto scaling) instance.
---

# tencentcloud_as_scaling_config

Provides a resource to create a configuration for an AS (Auto scaling) instance.

~> **NOTE:**  In order to ensure the integrity of customer data, if the cvm instance was destroyed due to shrinking, it will keep the cbs associate with cvm by default. If you want to destroy together, please set `delete_with_instance` to `true`.

## Example Usage

```hcl
resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "launch-configuration"
  image_id           = "img-9qabwvbn"
  instance_types     = ["SA1.SMALL1"]
  project_id         = 0
  system_disk_type   = "CLOUD_PREMIUM"
  system_disk_size   = "50"

  data_disk {
    disk_type = "CLOUD_PREMIUM"
    disk_size = 50
  }

  internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
  internet_max_bandwidth_out = 10
  public_ip_assigned         = true
  password                   = "test123#"
  enhanced_security_service  = false
  enhanced_monitor_service   = false
  user_data                  = "dGVzdA=="

  instance_tags = {
    tag = "as"
  }
}
```

Using SPOT charge type

```hcl
resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name   = "launch-configuration"
  image_id             = "img-9qabwvbn"
  instance_types       = ["SA1.SMALL1"]
  instance_charge_type = "SPOTPAID"
  spot_instance_type   = "one-time"
  spot_max_price       = "1000"
}
```

## Argument Reference

The following arguments are supported:

* `configuration_name` - (Required, String) Name of a launch configuration.
* `image_id` - (Required, String) An available image ID for a cvm instance.
* `instance_types` - (Required, List: [`String`]) Specified types of CVM instances.
* `cam_role_name` - (Optional, String) CAM role name authorized to access.
* `data_disk` - (Optional, List) Configurations of data disk.
* `disk_type_policy` - (Optional, String) Policy of cloud disk type. Valid values: `ORIGINAL` and `AUTOMATIC`. Default is `ORIGINAL`.
* `enhanced_monitor_service` - (Optional, Bool) To specify whether to enable cloud monitor service. Default is `TRUE`.
* `enhanced_security_service` - (Optional, Bool) To specify whether to enable cloud security service. Default is `TRUE`.
* `instance_charge_type_prepaid_period` - (Optional, Int) The tenancy (in month) of the prepaid instance, NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.
* `instance_charge_type_prepaid_renew_flag` - (Optional, String) Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.
* `instance_charge_type` - (Optional, String) Charge type of instance. Valid values are `PREPAID`, `POSTPAID_BY_HOUR`, `SPOTPAID`. The default is `POSTPAID_BY_HOUR`. NOTE: `SPOTPAID` instance must set `spot_instance_type` and `spot_max_price` at the same time.
* `instance_name_settings` - (Optional, List) Settings of CVM instance names.
* `instance_tags` - (Optional, Map) A list of tags used to associate different resources.
* `internet_charge_type` - (Optional, String) Charge types for network traffic. Valid values: `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR`, `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.
* `internet_max_bandwidth_out` - (Optional, Int) Max bandwidth of Internet access in Mbps. Default is `0`.
* `keep_image_login` - (Optional, Bool) Specify whether to keep original settings of a CVM image. And it can't be used with password or key_ids together.
* `key_ids` - (Optional, List: [`String`]) ID list of keys.
* `password` - (Optional, String) Password to access.
* `project_id` - (Optional, Int) Specifys to which project the configuration belongs.
* `public_ip_assigned` - (Optional, Bool) Specify whether to assign an Internet IP address.
* `security_group_ids` - (Optional, List: [`String`]) Security groups to which a CVM instance belongs.
* `spot_instance_type` - (Optional, String) Type of spot instance, only support `one-time` now. Note: it only works when instance_charge_type is set to `SPOTPAID`.
* `spot_max_price` - (Optional, String) Max price of a spot instance, is the format of decimal string, for example "0.50". Note: it only works when instance_charge_type is set to `SPOTPAID`.
* `system_disk_size` - (Optional, Int) Volume of system disk in GB. Default is `50`.
* `system_disk_type` - (Optional, String) Type of a CVM disk. Valid values: `CLOUD_PREMIUM` and `CLOUD_SSD`. Default is `CLOUD_PREMIUM`. valid when disk_type_policy is ORIGINAL.
* `user_data` - (Optional, String) ase64-encoded User Data text, the length limit is 16KB.

The `data_disk` object supports the following:

* `delete_with_instance` - (Optional, Bool) Indicates whether the disk remove after instance terminated. Default is `false`.
* `disk_size` - (Optional, Int) Volume of disk in GB. Default is `0`.
* `disk_type` - (Optional, String) Types of disk. Valid values: `CLOUD_PREMIUM` and `CLOUD_SSD`. valid when disk_type_policy is ORIGINAL.
* `snapshot_id` - (Optional, String) Data disk snapshot ID.

The `instance_name_settings` object supports the following:

* `instance_name` - (Required, String) CVM instance name.
* `instance_name_style` - (Optional, String) Type of CVM instance name. Valid values: `ORIGINAL` and `UNIQUE`. Default is `ORIGINAL`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The time when the launch configuration was created.
* `status` - Current statues of a launch configuration.


## Import

AutoScaling Configuration can be imported using the id, e.g.

```
$ terraform import tencentcloud_as_scaling_config.scaling_config asc-n32ymck2
```

