---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_binlog_save_days"
sidebar_current: "docs-tencentcloud-resource-cynosdb_binlog_save_days"
description: |-
  Provides a resource to create a cynosdb binlog_save_days
---

# tencentcloud_cynosdb_binlog_save_days

Provides a resource to create a cynosdb binlog_save_days

## Example Usage

```hcl
resource "tencentcloud_cynosdb_binlog_save_days" "binlog_save_days" {
  cluster_id       = "cynosdbmysql-123"
  binlog_save_days = 7
}
```

## Argument Reference

The following arguments are supported:

* `binlog_save_days` - (Required, Int) Binlog retention days.
* `cluster_id` - (Required, String, ForceNew) Cluster ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cynosdb binlog_save_days can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_binlog_save_days.binlog_save_days binlog_save_days_id
```

