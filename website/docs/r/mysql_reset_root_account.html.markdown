---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_reset_root_account"
sidebar_current: "docs-tencentcloud-resource-mysql_reset_root_account"
description: |-
  Provides a resource to create a mysql reset_root_account
---

# tencentcloud_mysql_reset_root_account

Provides a resource to create a mysql reset_root_account

## Example Usage

```hcl
resource "tencentcloud_mysql_reset_root_account" "reset_root_account" {
  instance_id = "cdb-fitq5t9h"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



