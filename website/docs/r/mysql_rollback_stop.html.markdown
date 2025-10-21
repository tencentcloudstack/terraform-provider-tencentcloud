---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_rollback_stop"
sidebar_current: "docs-tencentcloud-resource-mysql_rollback_stop"
description: |-
  Provides a resource to create a mysql rollback_stop
---

# tencentcloud_mysql_rollback_stop

Provides a resource to create a mysql rollback_stop

## Example Usage

### Revoke the ongoing rollback task of the instance

```hcl
resource "tencentcloud_mysql_rollback_stop" "example" {
  instance_id = "cdb-fitq5t9h"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Cloud database instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



