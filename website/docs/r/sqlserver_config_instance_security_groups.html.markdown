---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_config_instance_security_groups"
sidebar_current: "docs-tencentcloud-resource-sqlserver_config_instance_security_groups"
description: |-
  Provides a resource to create a sqlserver config_instance_security_groups
---

# tencentcloud_sqlserver_config_instance_security_groups

Provides a resource to create a sqlserver config_instance_security_groups

## Example Usage

```hcl
resource "tencentcloud_sqlserver_config_instance_security_groups" "config_instance_security_groups" {
  instance_id           = "mssql-qelbzgwf"
  security_group_id_set = ["sg-mayqdlt1", "sg-5aubsf8n"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `security_group_id_set` - (Required, Set: [`String`]) A list of security group IDs to modify, an array of one or more security group IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver config_instance_security_groups can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_instance_security_groups.config_instance_security_groups config_instance_security_groups_id
```

