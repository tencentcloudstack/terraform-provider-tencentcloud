---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_scaling_configs"
sidebar_current: "docs-tencentcloud-datasource-as_scaling_configs"
description: |-
  Use this data source to query scaling configuration information.
---

# tencentcloud_as_scaling_configs

Use this data source to query scaling configuration information.

## Example Usage

```hcl
data "tencentcloud_as_scaling_configs" "as_configs" {
    configuration_id   = "asc-oqio4yyj"
    result_output_file = "my_test_path"
}
```

## Argument Reference

The following arguments are supported:

* `configuration_id` - (Optional) Launch configuration ID.
* `configuration_name` - (Optional) Launch configuration name.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `configuration_list` - A list of configuration. Each element contains the following attributes:
  * `configuration_id` - Launch configuration ID.
  * `configuration_name` - Launch configuration name.
  * `create_time` - The time when the launch configuration was created.
  * `data_disk` - Configurations of data disk.
    * `disk_size` - Volume of disk in GB. Default is 0.
    * `disk_type` - Type of disk.
    * `snapshot_id` - Data disk snapshot ID.
  * `enhanced_monitor_service` - Whether to activate cloud monitor service.
  * `enhanced_security_service` - Whether to activate cloud security service.
  * `image_id` - ID of available image, for example img-8toqc6s3.
  * `instance_tags` - A tag list associates with an instance.
  * `instance_types` - Instance type list of the scaling configuration.
  * `internet_charge_type` - Charge types for network traffic.
  * `internet_max_bandwidth_out` - Max bandwidth of Internet access in Mbps.
  * `key_ids` - ID list of login keys
  * `project_id` - ID of the project to which the configuration belongs. Default value is 0.
  * `public_ip_assigned` - Specify whether to assign an Internet IP address.
  * `security_group_ids` - Security groups to which the instance belongs.
  * `status` - Current statues of a launch configuration.
  * `system_disk_size` - System disk size of the scaling configuration in GB.
  * `system_disk_type` - System disk category of the scaling configuration.
  * `user_data` - Base64-encoded User Data text.


