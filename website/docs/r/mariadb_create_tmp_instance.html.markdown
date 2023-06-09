---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_create_tmp_instance"
sidebar_current: "docs-tencentcloud-resource-mariadb_create_tmp_instance"
description: |-
  Provides a resource to create a mariadb create_tmp_instance
---

# tencentcloud_mariadb_create_tmp_instance

Provides a resource to create a mariadb create_tmp_instance

## Example Usage

```hcl
resource "tencentcloud_mariadb_create_tmp_instance" "create_tmp_instance" {
  instance_id   = "tdsql-9vqvls95"
  rollback_time = "2023-06-05 01:00:00"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `rollback_time` - (Required, String, ForceNew) Rollback time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



