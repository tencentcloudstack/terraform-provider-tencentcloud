---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_libra_db_instance_attachment"
sidebar_current: "docs-tencentcloud-resource-cynosdb_libra_db_instance_attachment"
description: |-
  Provides a resource to create a CynosDB (TDSQL-C) LibraDB read-only analytics engine instance attachment
---

# tencentcloud_cynosdb_libra_db_instance_attachment

Provides a resource to create a CynosDB (TDSQL-C) LibraDB read-only analytics engine instance attachment

## Example Usage

```hcl
resource "tencentcloud_cynosdb_libra_db_instance_attachment" "example" {
  cluster_id       = "cynosdbmysql-xxxxxxxx"
  zone             = "ap-guangzhou-3"
  cpu              = 4
  mem              = 8
  storage_size     = 100
  pay_mode         = 0
  instance_name    = "tf-example"
  instance_type    = "Common"
  storage_type     = "CLOUD_SSD"
  vpc_id           = "vpc-xxxxxxxx"
  subnet_id        = "subnet-xxxxxxxx"
  libra_db_version = "3.1.2"
  src_instance_id  = "cynosdbmysql-ins-xxxxxxxx"
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
* `goods_num` - (Optional, Int) Number of new read-only instances, value range (0, 15].
* `instance_name` - (Optional, String) Instance name.
* `instance_type` - (Optional, String) Instance type.
* `isolate_reason_types` - (Optional, List: [`Int`]) Isolation reason types for delete operation.
* `isolate_reason` - (Optional, String) Isolation reason for delete operation.
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
* `big_deal_ids` - Big deal IDs.
* `deal_names` - Post-paid order names.
* `instance_id` - Instance ID.
* `resource_ids` - Resource ID list.
* `tran_id` - Frozen transaction ID.


## Import

CynosDB LibraDB instance attachment can be imported using the cluster_id#instance_id, e.g.

```
terraform import tencentcloud_cynosdb_libra_db_instance_attachment.example cynosdbmysql-xxxxxxxx#cynosdbmysql-ins-yyyyyyyy
```

