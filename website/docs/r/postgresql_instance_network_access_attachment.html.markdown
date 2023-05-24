---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_instance_network_access_attachment"
sidebar_current: "docs-tencentcloud-resource-postgresql_instance_network_access_attachment"
description: |-
  Provides a resource to create a postgresql instance_network_access_attachment
---

# tencentcloud_postgresql_instance_network_access_attachment

Provides a resource to create a postgresql instance_network_access_attachment

## Example Usage

Vip assigned by system.

```hcl
resource "tencentcloud_postgresql_instance_network_access_attachment" "instance_network_access_attachment" {
  db_instance_id = tencentcloud_postgresql_instance.test.id
  vpc_id         = local.vpc_id
  subnet_id      = local.subnet_id
  is_assign_vip  = false
  tags = {
    "createdBy" = "terraform"
  }
}
```

Vip specified by user.

```hcl
resource "tencentcloud_postgresql_instance_network_access_attachment" "instance_network_access_attachment" {
  db_instance_id = tencentcloud_postgresql_instance.test.id
  vpc_id         = local.my_vpc_id
  subnet_id      = local.my_subnet_id
  is_assign_vip  = true
  vip            = "172.18.111.111"
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) Instance ID in the format of postgres-6bwgamo3.
* `is_assign_vip` - (Required, Bool, ForceNew) Whether to manually assign the VIP. Valid values:true (manually assign), false (automatically assign).
* `subnet_id` - (Required, String, ForceNew) Subnet ID.
* `vpc_id` - (Required, String, ForceNew) Unified VPC ID.
* `tags` - (Optional, Map, ForceNew) Tag description list.
* `vip` - (Optional, String, ForceNew) Target VIP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

postgresql instance_network_access_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_instance_network_access_attachment.instance_network_access_attachment instance_network_access_attachment_id
```

