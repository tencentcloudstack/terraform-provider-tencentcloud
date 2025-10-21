---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_switch_record"
sidebar_current: "docs-tencentcloud-datasource-mysql_switch_record"
description: |-
  Use this data source to query detailed information of mysql switch_record
---

# tencentcloud_mysql_switch_record

Use this data source to query detailed information of mysql switch_record

## Example Usage

```hcl
data "tencentcloud_mysql_switch_record" "switch_record" {
  instance_id = "cdb-fitq5t9h"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID, the format is: cdb-c1nl9rpv or cdbro-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Instance switching record details.
  * `switch_time` - Switching time, the format is: 2017-09-03 01:34:31.
  * `switch_type` - Switch type, possible return values: TRANSFER - data migration; MASTER2SLAVE - master-standby switch; RECOVERY - master-slave recovery.


