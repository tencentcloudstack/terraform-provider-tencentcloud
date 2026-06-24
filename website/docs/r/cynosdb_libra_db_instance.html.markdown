---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_libra_db_instance"
sidebar_current: "docs-tencentcloud-resource-cynosdb_libra_db_instance"
description: |-
  Provides a resource to create a CynosDB (TDSQL-C) LibraDB read-only analytics engine instance
---

# tencentcloud_cynosdb_libra_db_instance

Provides a resource to create a CynosDB (TDSQL-C) LibraDB read-only analytics engine instance

## Example Usage

```hcl
resource "tencentcloud_cynosdb_libra_db_instance" "example" {
  cluster_id         = "cynosdbmysql-5oo78wv9"
  zone               = "ap-guangzhou-7"
  cpu                = 8
  mem                = 32
  storage_size       = 100
  pay_mode           = 0
  port               = 2000
  instance_name      = "tf-example"
  instance_type      = "Common"
  storage_type       = "CLOUD_TCS"
  vpc_id             = "vpc-i5yyodl9"
  subnet_id          = "subnet-5rrirqyc"
  libra_db_version   = "2.2410.18.0"
  src_instance_id    = "cynosdbmysql-ins-84ja0ye0"
  security_group_ids = ["sg-4rd5741x"]
  force_delete       = true
  objects {
    database_tables {
      migrate_db_mode = "all"
    }
  }
}
```

### or

```hcl
resource "tencentcloud_cynosdb_libra_db_instance" "example" {
  cluster_id         = "cynosdbmysql-5oo78wv9"
  zone               = "ap-guangzhou-7"
  cpu                = 8
  mem                = 32
  storage_size       = 100
  pay_mode           = 0
  port               = 2000
  instance_name      = "tf-example"
  instance_type      = "Common"
  storage_type       = "CLOUD_TCS"
  vpc_id             = "vpc-i5yyodl9"
  subnet_id          = "subnet-5rrirqyc"
  libra_db_version   = "2.2410.18.0"
  src_instance_id    = "cynosdbmysql-ins-84ja0ye0"
  security_group_ids = ["sg-4rd5741x"]
  force_delete       = true
  objects {
    database_tables {
      migrate_db_mode = "partial"
      databases {
        db_name            = "test"
        migrate_table_mode = "all"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `cpu` - (Required, Int) Number of CPU cores.
* `mem` - (Required, Int) Memory size in GB.
* `storage_size` - (Required, Int) Disk size.
* `zone` - (Required, String) Availability zone.
* `auto_voucher` - (Optional, Int) Whether to automatically select vouchers: 1 yes, 0 no, default 0.
* `deal_mode` - (Optional, Int) Transaction mode: 0 - place order and pay, 1 - place order only.
* `force_delete` - (Optional, Bool) Whether to force delete the instance. Default is false.
* `instance_name` - (Optional, String) Instance name.
* `instance_type` - (Optional, String) Instance type.
* `libra_db_version` - (Optional, String) Analytics engine version.
* `objects` - (Optional, List) Sync object list.
* `order_source` - (Optional, String) Order source.
* `pay_mode` - (Optional, Int) Payment mode.
* `port` - (Optional, Int) Port for the new RO group, value range [0, 65535).
* `replicas_num` - (Optional, Int) Number of replicas.
* `security_group_ids` - (Optional, List: [`String`]) Security group IDs for the new read-only instance.
* `src_instance_id` - (Optional, String) Source instance ID.
* `storage_type` - (Optional, String) Disk type.
* `subnet_id` - (Optional, String) Subnet ID. Required if VpcId is set.
* `time_span` - (Optional, Int) Purchase duration, effective with TimeUnit.
* `time_unit` - (Optional, String) Purchase duration unit. Options: d (day), m (month).
* `vpc_id` - (Optional, String) VPC network ID.

The `database_tables` object of `objects` supports the following:

* `databases` - (Optional, List) Database information list.
* `migrate_db_mode` - (Optional, String) Database migration mode.

The `databases` object of `database_tables` supports the following:

* `db_name` - (Optional, String) Database name.
* `migrate_table_mode` - (Optional, String) Table migration mode.
* `tables` - (Optional, List) Table information list.

The `objects` object supports the following:

* `database_tables` - (Optional, List) Database table information.

The `tables` object of `databases` supports the following:

* `table_name` - (Optional, String) Table name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_id` - Instance ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `10m`) Used when creating the resource.
* `delete` - (Defaults to `10m`) Used when deleting the resource.

## Import

CynosDB LibraDB instance can be imported using the cluster_id#instance_id, e.g.

```
terraform import tencentcloud_cynosdb_libra_db_instance.example cynosdbmysql-5oo78wv9#libradb-ins-irehx3rm
```

