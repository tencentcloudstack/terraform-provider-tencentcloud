---
subcategory: "Business Intelligence(BI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bi_datasource"
sidebar_current: "docs-tencentcloud-resource-bi_datasource"
description: |-
  Provides a resource to create a bi datasource
---

# tencentcloud_bi_datasource

Provides a resource to create a bi datasource

## Example Usage

```hcl
resource "tencentcloud_bi_datasource" "datasource" {
  charset     = "utf8"
  db_host     = "bj-cdb-1lxqg5r6.sql.tencentcdb.com"
  db_name     = "tf-test"
  db_port     = 63694
  db_type     = "MYSQL"
  db_pwd      = "ABc123,,,"
  db_user     = "root"
  project_id  = 11015030
  source_name = "tf-source-name"
}
```

## Argument Reference

The following arguments are supported:

* `charset` - (Required, String) Charset.
* `db_host` - (Required, String) Host.
* `db_name` - (Required, String) Database name.
* `db_port` - (Required, Int) Port.
* `db_pwd` - (Required, String) Password.
* `db_type` - (Required, String) `MYSQL`, `MSSQL`, `POSTGRE`, `ORACLE`, `CLICKHOUSE`, `TIDB`, `HIVE`, `PRESTO`.
* `db_user` - (Required, String) User name.
* `project_id` - (Required, Int) Project id.
* `source_name` - (Required, String) Datasource name in BI.
* `catalog` - (Optional, String) Catalog.
* `data_origin_datasource_id` - (Optional, String) Third-party datasource project id, this parameter can be ignored.
* `data_origin_project_id` - (Optional, String) Third-party datasource project id, this parameter can be ignored.
* `data_origin` - (Optional, String) Third-party datasource identification, this parameter can be ignored.
* `service_type` - (Optional, String) Own or Cloud, default: `Own`.
* `uniq_vpc_id` - (Optional, String) Tencent cloud private network unified identity.
* `vpc_id` - (Optional, String) Tencent cloud private network identity.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

bi datasource can be imported using the id, e.g.

```
terraform import tencentcloud_bi_datasource.datasource datasource_id
```

