---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_ro_start_replication"
sidebar_current: "docs-tencentcloud-resource-mysql_ro_start_replication"
description: |-
  Provides a resource to create a mysql ro_start_replication
---

# tencentcloud_mysql_ro_start_replication

Provides a resource to create a mysql ro_start_replication

## Example Usage

```hcl
resource "tencentcloud_mysql_ro_start_replication" "ro_start_replication" {
  instance_id = "cdb-fitq5t9h"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Read-Only instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



