---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cluster_databases"
sidebar_current: "docs-tencentcloud-resource-cynosdb_cluster_databases"
description: |-
  Provides a resource to create a cynosdb cluster_databases
---

# tencentcloud_cynosdb_cluster_databases

Provides a resource to create a cynosdb cluster_databases

## Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_databases" "cluster_databases" {
  cluster_id    = "cynosdbmysql-bws8h88b"
  db_name       = "terraform-test"
  character_set = "utf8"
  collate_rule  = "utf8_general_ci"
  user_host_privileges {
    db_user_name = "root"
    db_host      = "%"
    db_privilege = "READ_ONLY"
  }
  description = "terraform test"
}
```

## Argument Reference

The following arguments are supported:

* `character_set` - (Required, String) Character Set Type.
* `cluster_id` - (Required, String) Cluster ID.
* `collate_rule` - (Required, String) Sort Rules.
* `db_name` - (Required, String) Database name.
* `description` - (Optional, String) Remarks.
* `user_host_privileges` - (Optional, List) Authorize user host permissions.

The `user_host_privileges` object supports the following:

* `db_host` - (Required, String) .
* `db_privilege` - (Required, String) .
* `db_user_name` - (Required, String) Authorized Users.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cynosdb cluster_databases can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_databases.cluster_databases cluster_databases_id
```

