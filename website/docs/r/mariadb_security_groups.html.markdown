---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_security_groups"
sidebar_current: "docs-tencentcloud-resource-mariadb_security_groups"
description: |-
  Provides a resource to create a mariadb security_groups
---

# tencentcloud_mariadb_security_groups

Provides a resource to create a mariadb security_groups

## Example Usage

```hcl
resource "tencentcloud_mariadb_security_groups" "security_groups" {
  instance_id       = "tdsql-4pzs5b67"
  security_group_id = "sg-7kpsbxdb"
  product           = "mariadb"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `product` - (Required, String) product name, fixed to mariadb.
* `security_group_id` - (Required, String) security group id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mariadb security_groups can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_security_groups.security_groups tdsql-4pzs5b67#sg-7kpsbxdb#mariadb
```

