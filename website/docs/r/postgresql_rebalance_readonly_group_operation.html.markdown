---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_rebalance_readonly_group_operation"
sidebar_current: "docs-tencentcloud-resource-postgresql_rebalance_readonly_group_operation"
description: |-
  Provides a resource to create a postgresql rebalance_readonly_group_operation
---

# tencentcloud_postgresql_rebalance_readonly_group_operation

Provides a resource to create a postgresql rebalance_readonly_group_operation

## Example Usage

```hcl
resource "tencentcloud_postgresql_rebalance_readonly_group_operation" "rebalance_readonly_group_operation" {
  read_only_group_id = "pgrogrp-test"
}
```

## Argument Reference

The following arguments are supported:

* `read_only_group_id` - (Required, String, ForceNew) readonly Group ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

postgresql rebalance_readonly_group_operation can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_rebalance_readonly_group_operation.rebalance_readonly_group_operation rebalance_readonly_group_operation_id
```

