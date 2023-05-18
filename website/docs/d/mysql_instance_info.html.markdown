---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_instance_info"
sidebar_current: "docs-tencentcloud-datasource-mysql_instance_info"
description: |-
  Use this data source to query detailed information of mysql instance_info
---

# tencentcloud_mysql_instance_info

Use this data source to query detailed information of mysql instance_info

## Example Usage

```hcl
data "tencentcloud_mysql_instance_info" "instance_info" {
  instance_id = "cdb-fitq5t9h"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `default_kms_region` - The default region of the KMS service used by the current CDB backend service.
* `encryption` - Whether to enable encryption, YES is enabled, NO is not enabled.
* `instance_name` - instance name.
* `key_id` - The key ID used for encryption.
* `key_region` - The region where the key is located.


