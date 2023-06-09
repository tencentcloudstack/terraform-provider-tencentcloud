---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_complete_expansion"
sidebar_current: "docs-tencentcloud-resource-sqlserver_complete_expansion"
description: |-
  Provides a resource to create a sqlserver complete_expansion
---

# tencentcloud_sqlserver_complete_expansion

Provides a resource to create a sqlserver complete_expansion

## Example Usage

```hcl
resource "tencentcloud_sqlserver_complete_expansion" "complete_expansion" {
  instance_id = "mssql-qelbzgwf"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of imported target instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



