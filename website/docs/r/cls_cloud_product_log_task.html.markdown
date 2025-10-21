---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_cloud_product_log_task"
sidebar_current: "docs-tencentcloud-resource-cls_cloud_product_log_task"
description: |-
  Provides a resource to create a cls cloud product log task
---

# tencentcloud_cls_cloud_product_log_task

Provides a resource to create a cls cloud product log task

~> **NOTE:** This resource has been deprecated in Terraform TencentCloud provider version 1.81.188. Please use `tencentcloud_cls_cloud_product_log_task_v2` instead.

~> **NOTE:** Using this resource will create new `logset` and `topic`

## Example Usage

```hcl
resource "tencentcloud_cls_cloud_product_log_task" "example" {
  instance_id          = "postgres-1p7xvpc1"
  assumer_name         = "PostgreSQL"
  log_type             = "PostgreSQL-SLOW"
  cloud_product_region = "gz"
  cls_region           = "ap-guangzhou"
  logset_name          = "tf-example"
  topic_name           = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `assumer_name` - (Required, String) Cloud product identification, Values: CDS, CWP, CDB, TDSQL-C, MongoDB, TDStore, DCDB, MariaDB, PostgreSQL, BH, APIS.
* `cloud_product_region` - (Required, String) Cloud product region. There are differences in the input format of different log types in different regions. Please refer to the following example:
- CDS(all log type): ap-guangzhou
- CDB-AUDIT: gz
- TDSQL-C-AUDIT: gz
- MongoDB-AUDIT: gz
- MongoDB-SlowLog: ap-guangzhou
- MongoDB-ErrorLog: ap-guangzhou
- TDMYSQL-SLOW: gz
- DCDB(all log type): gz
- MariaDB(all log type): gz
- PostgreSQL(all log type): gz
- BH(all log type): overseas-polaris(Domestic sites overseas)/fsi-polaris(Domestic sites finance)/general-polaris(Domestic sites)/intl-sg-prod(International sites)
- APIS(all log type): gz.
* `cls_region` - (Required, String) CLS target region.
* `instance_id` - (Required, String) Instance ID.
* `log_type` - (Required, String) Log type, Values: CDS-AUDIT, CDS-RISK, CDB-AUDIT, TDSQL-C-AUDIT, MongoDB-AUDIT, MongoDB-SlowLog, MongoDB-ErrorLog, TDMYSQL-SLOW, DCDB-AUDIT, DCDB-SLOW, DCDB-ERROR, MariaDB-AUDIT, MariaDB-SLOW, MariaDB-ERROR, PostgreSQL-SLOW, PostgreSQL-ERROR, PostgreSQL-AUDIT, BH-FILELOG, BH-COMMANDLOG, APIS-ACCESS.
* `extend` - (Optional, String) Log configuration extension information, generally used to store additional log delivery configurations.
* `logset_name` - (Optional, String) Log set name, it will be automatically created.
* `topic_name` - (Optional, String) The name of the log topic, it will be automatically created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `logset_id` - Log set ID.
* `topic_id` - Log theme ID.


## Import

cls cloud product log task can be imported using the id, e.g.

```
terraform import tencentcloud_cls_cloud_product_log_task.example postgres-1p7xvpc1#PostgreSQL#PostgreSQL-SLOW#gz
```

