---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_instance_network_access"
sidebar_current: "docs-tencentcloud-resource-postgresql_instance_network_access"
description: |-
  Provides a resource to create a postgres instance network access
---

# tencentcloud_postgresql_instance_network_access

Provides a resource to create a postgres instance network access

## Example Usage

### Create by custom vip

```hcl
resource "tencentcloud_postgresql_instance_network_access" "example" {
  db_instance_id = "postgres-ai46555b"
  vpc_id         = "vpc-i5yyodl9"
  subnet_id      = "subnet-d4umunpy"
  vip            = "10.0.10.11"
}
```

### Create by automatic allocation vip

```hcl
resource "tencentcloud_postgresql_instance_network_access" "example" {
  db_instance_id = "postgres-ai46555b"
  vpc_id         = "vpc-i5yyodl9"
  subnet_id      = "subnet-d4umunpy"
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) Instance ID in the format of postgres-6bwgamo3.
* `subnet_id` - (Required, String, ForceNew) Subnet ID.
* `vpc_id` - (Required, String, ForceNew) Unified VPC ID.
* `vip` - (Optional, String, ForceNew) Target VIP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

postgres instance network access can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_instance_network_access.example postgres-ai46555b#vpc-i5yyodl9#subnet-d4umunpy#10.0.10.11
```

