---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_readonly_group_network_access_attachment"
sidebar_current: "docs-tencentcloud-resource-postgresql_readonly_group_network_access_attachment"
description: |-
  Provides a resource to create a postgresql readonly_group_network_access_attachment
---

# tencentcloud_postgresql_readonly_group_network_access_attachment

Provides a resource to create a postgresql readonly_group_network_access_attachment

## Example Usage

```hcl
resource "tencentcloud_postgresql_readonly_group_network_access_attachment" "readonly_group_network_access_attachment" {
  readonly_group_id = "pgro-xxxx"
  vpc_id            = "vpc-xxx"
  subnet_id         = "subnet-xxx"
  is_assign_vip     = false
  vip               = ""
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) Master database instance ID.
* `is_assign_vip` - (Required, Bool, ForceNew) Whether to manually assign the VIP. Valid values:true (manually assign), false (automatically assign).
* `readonly_group_id` - (Required, String, ForceNew) RO group identifier.
* `subnet_id` - (Required, String, ForceNew) Subnet ID.
* `vpc_id` - (Required, String, ForceNew) Unified VPC ID.
* `tags` - (Optional, Map, ForceNew) Tag description list.
* `vip` - (Optional, String, ForceNew) Target VIP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

postgresql readonly_group_network_access_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_readonly_group_network_access_attachment.readonly_group_network_access_attachment db_instance_id#readonly_group_id#vpc_id#vip#port
```

