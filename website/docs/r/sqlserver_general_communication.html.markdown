---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_general_communication"
sidebar_current: "docs-tencentcloud-resource-sqlserver_general_communication"
description: |-
  Provides a resource to create a sqlserver general_communication
---

# tencentcloud_sqlserver_general_communication

Provides a resource to create a sqlserver general_communication

## Example Usage

```hcl
resource "tencentcloud_sqlserver_general_communication" "general_communication" {
  instance_id = "mssql-qelbzgwf"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of instances.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver general_communication can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_general_communication.general_communication general_communication_id
```

