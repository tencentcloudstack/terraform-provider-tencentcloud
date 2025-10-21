---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_db_features"
sidebar_current: "docs-tencentcloud-datasource-mysql_db_features"
description: |-
  Use this data source to query detailed information of mysql db_features
---

# tencentcloud_mysql_db_features

Use this data source to query detailed information of mysql db_features

## Example Usage

```hcl
data "tencentcloud_mysql_db_features" "db_features" {
  instance_id = "cdb-fitq5t9h"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID, the format is: cdb-c1nl9rpv or cdbro-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `audit_need_upgrade` - Whether to enable auditing needs to upgrade the kernel version.
* `current_sub_version` - Current kernel version.
* `encryption_need_upgrade` - Whether to enable encryption needs to upgrade the kernel version.
* `is_remote_ro` - Whether it is a remote read-only instance.
* `is_support_audit` - Whether to support the database audit function.
* `is_support_encryption` - Whether to support the database encryption function.
* `is_support_update_sub_version` - Whether to support minor version upgrades.
* `master_region` - The region where the master instance is located.
* `target_sub_version` - Available kernel versions for upgrade.


