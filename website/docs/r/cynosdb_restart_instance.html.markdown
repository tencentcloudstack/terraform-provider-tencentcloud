---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_restart_instance"
sidebar_current: "docs-tencentcloud-resource-cynosdb_restart_instance"
description: |-
  Provides a resource to create a cynosdb restart_instance
---

# tencentcloud_cynosdb_restart_instance

Provides a resource to create a cynosdb restart_instance

## Example Usage

```hcl
resource "tencentcloud_cynosdb_restart_instance" "restart_instance" {
  instance_id = "cynosdbmysql-ins-afqx1hy0"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - instance state.


