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

Vip assigned by system.

```hcl
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "172.18.111.0/24"
  name       = "test-pg-network-vpc"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.default_az
  cidr_block        = "172.18.111.0/24"
  name              = "test-pg-network-sub1"
  vpc_id            = tencentcloud_vpc.vpc.id
}

locals {
  my_vpc_id    = tencentcloud_subnet.subnet.vpc_id
  my_subnet_id = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_postgresql_readonly_group" "group" {
  master_db_instance_id       = local.pgsql_id
  name                        = "tf_test_postgresql_readonly_group"
  project_id                  = 0
  vpc_id                      = local.my_vpc_id
  subnet_id                   = local.my_subnet_id
  replay_lag_eliminate        = 1
  replay_latency_eliminate    = 1
  max_replay_lag              = 100
  max_replay_latency          = 512
  min_delay_eliminate_reserve = 1
}

resource "tencentcloud_postgresql_readonly_group_network_access_attachment" "readonly_group_network_access_attachment" {
  db_instance_id    = local.pgsql_id
  readonly_group_id = tencentcloud_postgresql_readonly_group.group.id
  vpc_id            = local.my_vpc_id
  subnet_id         = local.my_subnet_id
  is_assign_vip     = false
  tags = {
    "createdBy" = "terraform"
  }
}
```

Vip specified by user.

```hcl
resource "tencentcloud_postgresql_readonly_group_network_access_attachment" "readonly_group_network_access_attachment" {
  db_instance_id    = local.pgsql_id
  readonly_group_id = tencentcloud_postgresql_readonly_group.group.id
  vpc_id            = local.my_vpc_id
  subnet_id         = local.my_subnet_id
  is_assign_vip     = true
  vip               = "172.18.111.111"
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

