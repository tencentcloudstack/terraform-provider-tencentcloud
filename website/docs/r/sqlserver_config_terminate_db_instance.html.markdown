---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_config_terminate_db_instance"
sidebar_current: "docs-tencentcloud-resource-sqlserver_config_terminate_db_instance"
description: |-
  Provides a resource to create a sqlserver config_terminate_db_instance
---

# tencentcloud_sqlserver_config_terminate_db_instance

Provides a resource to create a sqlserver config_terminate_db_instance

## Example Usage

```hcl
resource "tencentcloud_sqlserver_config_terminate_db_instance" "config_terminate_db_instance" {
  instance_id = "mssql-4tgeyeeh"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver config_terminate_db_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_terminate_db_instance.config_terminate_db_instance config_terminate_db_instance_id
```

