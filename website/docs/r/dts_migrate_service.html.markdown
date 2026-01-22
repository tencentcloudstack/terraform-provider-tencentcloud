---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_migrate_service"
sidebar_current: "docs-tencentcloud-resource-dts_migrate_service"
description: |-
  Provides a resource to create a DTS migrate service
---

# tencentcloud_dts_migrate_service

Provides a resource to create a DTS migrate service

## Example Usage

```hcl
resource "tencentcloud_dts_migrate_service" "example" {
  src_database_type = "mysql"
  dst_database_type = "cynosdbmysql"
  src_region        = "ap-guangzhou"
  dst_region        = "ap-guangzhou"
  instance_class    = "small"
  job_name          = "tf-example"
  tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `dst_database_type` - (Required, String) destination database type, optional value is mysql/redis/percona/mongodb/postgresql/sqlserver/mariadb.
* `dst_region` - (Required, String) destination region.
* `instance_class` - (Required, String) instance class, optional value is small/medium/large/xlarge/2xlarge.
* `src_database_type` - (Required, String) source database type, optional value is mysql/redis/percona/mongodb/postgresql/sqlserver/mariadb.
* `src_region` - (Required, String) source region.
* `job_name` - (Optional, String) job name.
* `tags` - (Optional, List) tags.

The `tags` object supports the following:

* `tag_key` - (Optional, String) tag key.
* `tag_value` - (Optional, String) tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

DTS migrate service can be imported using the id, e.g.
```
$ terraform import tencentcloud_dts_migrate_service.example dts-iy98oxba
```

