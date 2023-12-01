---
subcategory: "Business Intelligence(BI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bi_datasource_cloud"
sidebar_current: "docs-tencentcloud-resource-bi_datasource_cloud"
description: |-
  Provides a resource to create a bi datasource_cloud
---

# tencentcloud_bi_datasource_cloud

Provides a resource to create a bi datasource_cloud

## Example Usage

```hcl
resource "tencentcloud_bi_datasource_cloud" "datasource_cloud" {
  charset    = "utf8"
  db_name    = "bi_dev"
  db_type    = "MYSQL"
  db_user    = "root"
  project_id = "11015056"
  db_pwd     = "xxxxxx"
  service_type {
    instance_id = "cdb-12viotu5"
    region      = "ap-guangzhou"
    type        = "Cloud"
  }
  source_name = "tf-test1"
  vip         = "10.0.0.4"
  vport       = "3306"
  region_id   = "gz"
  vpc_id      = 5292713
}
```

## Argument Reference

The following arguments are supported:

* `charset` - (Required, String) Charset.
* `db_name` - (Required, String) Database name.
* `db_pwd` - (Required, String) Password.
* `db_type` - (Required, String) `MYSQL`, `TDSQL-C_MYSQL`, `TDSQL_MYSQL`, `MSSQL`, `POSTGRESQL`, `MARIADB`.
* `db_user` - (Required, String) User name.
* `project_id` - (Required, String) Project id.
* `service_type` - (Required, List) Service type, Own or Cloud.
* `source_name` - (Required, String) Datasource name in BI.
* `vpc_id` - (Required, String) Vpc identification.
* `data_origin_datasource_id` - (Optional, String) Third-party datasource project id, this parameter can be ignored.
* `data_origin_project_id` - (Optional, String) Third-party datasource project id, this parameter can be ignored.
* `data_origin` - (Optional, String) Third-party datasource identification, this parameter can be ignored.
* `extra_param` - (Optional, String) Extended parameters.
* `region_id` - (Optional, String) Region identifier.
* `uniq_vpc_id` - (Optional, String) Unified vpc identification.
* `vip` - (Optional, String) Public cloud intranet ip.
* `vport` - (Optional, String) Public cloud intranet port.

The `service_type` object supports the following:

* `instance_id` - (Required, String) Instance Id.
* `region` - (Required, String) Region.
* `type` - (Required, String) Service type, Cloud.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



