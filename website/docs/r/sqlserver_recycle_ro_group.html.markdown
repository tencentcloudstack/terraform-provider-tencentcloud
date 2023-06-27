---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_recycle_ro_group"
sidebar_current: "docs-tencentcloud-resource-sqlserver_recycle_ro_group"
description: |-
  Provides a resource to create a sqlserver recycle_ro_group
---

# tencentcloud_sqlserver_recycle_ro_group

Provides a resource to create a sqlserver recycle_ro_group

## Example Usage

```hcl
resource "tencentcloud_sqlserver_recycle_ro_group" "recycle_ro_group" {
  instance_id        = "mssql-qelbzgwf"
  read_only_group_id = "mssqlrg-c9ld954d"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of the master instance.
* `read_only_group_id` - (Required, String, ForceNew) ID of the read-only group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



