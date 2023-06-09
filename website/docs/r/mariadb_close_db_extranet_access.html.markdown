---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_close_db_extranet_access"
sidebar_current: "docs-tencentcloud-resource-mariadb_close_db_extranet_access"
description: |-
  Provides a resource to create a mariadb close_db_extranet_access
---

# tencentcloud_mariadb_close_db_extranet_access

Provides a resource to create a mariadb close_db_extranet_access

## Example Usage

```hcl
resource "tencentcloud_mariadb_close_db_extranet_access" "close_db_extranet_access" {
  instance_id = "tdsql-9vqvls95"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of instance for which to enable public network access. The ID is in the format of `tdsql-ow728lmc` and can be obtained through the `DescribeDBInstances` API.
* `ipv6_flag` - (Optional, Int, ForceNew) Whether IPv6 is used. Default value: 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



