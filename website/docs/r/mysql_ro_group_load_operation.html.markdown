---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_ro_group_load_operation"
sidebar_current: "docs-tencentcloud-resource-mysql_ro_group_load_operation"
description: |-
  Provides a resource to create a mysql ro_group_load_operation
---

# tencentcloud_mysql_ro_group_load_operation

Provides a resource to create a mysql ro_group_load_operation

## Example Usage

```hcl
resource "tencentcloud_mysql_ro_group_load_operation" "ro_group_load_operation" {
  ro_group_id = "cdbrg-bdlvcfpj"
}
```

## Argument Reference

The following arguments are supported:

* `ro_group_id` - (Required, String, ForceNew) The ID of the RO group, in the format: cdbrg-c1nl9rpv.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



