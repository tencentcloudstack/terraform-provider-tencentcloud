---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_renew_instance"
sidebar_current: "docs-tencentcloud-resource-mariadb_renew_instance"
description: |-
  Provides a resource to create a mariadb renew_instance
---

# tencentcloud_mariadb_renew_instance

Provides a resource to create a mariadb renew_instance

## Example Usage

```hcl
resource "tencentcloud_mariadb_renew_instance" "renew_instance" {
  instance_id = "tdsql-9vqvls95"
  period      = 1
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `period` - (Required, Int, ForceNew) Renewal duration, unit: month.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



