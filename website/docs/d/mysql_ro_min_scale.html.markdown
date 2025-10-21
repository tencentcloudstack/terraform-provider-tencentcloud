---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_ro_min_scale"
sidebar_current: "docs-tencentcloud-datasource-mysql_ro_min_scale"
description: |-
  Use this data source to query detailed information of mysql ro_min_scale
---

# tencentcloud_mysql_ro_min_scale

Use this data source to query detailed information of mysql ro_min_scale

## Example Usage

```hcl
data "tencentcloud_mysql_ro_min_scale" "ro_min_scale" {
  # ro_instance_id = ""
  master_instance_id = "cdb-fitq5t9h"
}
```

## Argument Reference

The following arguments are supported:

* `master_instance_id` - (Optional, String) The primary instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page. This parameter and the RoInstanceId parameter cannot be empty at the same time. Note that when the input parameter contains RoInstanceId, the return value is the minimum specification when the read-only instance is upgraded; when the input parameter only contains MasterInstanceId, the return value is the minimum specification when the read-only instance is purchased.
* `result_output_file` - (Optional, String) Used to save results.
* `ro_instance_id` - (Optional, String) The read-only instance ID, in the format: cdbro-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page. This parameter and the MasterInstanceId parameter cannot be empty at the same time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `memory` - Memory specification size, unit: MB.
* `volume` - Disk specification size, unit: GB.


