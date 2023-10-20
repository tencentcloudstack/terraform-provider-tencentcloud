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
  service_type              = "Cloud"
  db_type                   = "Database type."
  charset                   = "utf8"
  db_user                   = "root"
  db_pwd                    = "abc"
  db_name                   = "abc"
  source_name               = "abc"
  project_id                = "123"
  vip                       = "1.2.3.4"
  vport                     = "3306"
  vpc_id                    = ""
  uniq_vpc_id               = ""
  region_id                 = ""
  extra_param               = ""
  instance_id               = ""
  prod_db_name              = ""
  data_origin               = "abc"
  data_origin_project_id    = "abc"
  data_origin_datasource_id = "abc"
}
```

## Argument Reference

The following arguments are supported:

* `charset` - (Required, String) Charset.
* `db_name` - (Required, String) Database name.
* `db_pwd` - (Required, String) Password.
* `db_type` - (Required, String) MYSQL.
* `db_user` - (Required, String) User name.
* `project_id` - (Required, String) Project id.
* `service_type` - (Required, String) Own or Cloud.
* `source_name` - (Required, String) Datasource name in BI.
* `data_origin_datasource_id` - (Optional, String) Third-party datasource project id, this parameter can be ignored.
* `data_origin_project_id` - (Optional, String) Third-party datasource project id, this parameter can be ignored.
* `data_origin` - (Optional, String) Third-party datasource identification, this parameter can be ignored.
* `extra_param` - (Optional, String) Extended parameters.
* `instance_id` - (Optional, String) Instance id.
* `prod_db_name` - (Optional, String) Datasource product name.
* `region_id` - (Optional, String) Region identifier.
* `uniq_vpc_id` - (Optional, String) Unified vpc identification.
* `vip` - (Optional, String) Public cloud intranet ip.
* `vpc_id` - (Optional, String) Vpc identification.
* `vport` - (Optional, String) Public cloud intranet port.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

bi datasource_cloud can be imported using the id, e.g.

```
terraform import tencentcloud_bi_datasource_cloud.datasource_cloud datasource_cloud_id
```

