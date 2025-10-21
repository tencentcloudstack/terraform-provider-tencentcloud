---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_instance_reboot_time"
sidebar_current: "docs-tencentcloud-datasource-mysql_instance_reboot_time"
description: |-
  Use this data source to query detailed information of mysql instance_reboot_time
---

# tencentcloud_mysql_instance_reboot_time

Use this data source to query detailed information of mysql instance_reboot_time

## Example Usage

```hcl
data "tencentcloud_mysql_instance_reboot_time" "instance_reboot_time" {
  instance_ids = ["cdb-fitq5t9h"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_ids` - (Required, Set: [`String`]) The instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Returned parameter information.
  * `instance_id` - Instance ID, the format is: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page.
  * `time_in_seconds` - expected restart time.


