---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_cloud_product_log_task_v2"
sidebar_current: "docs-tencentcloud-resource-cls_cloud_product_log_task_v2"
description: |-
  Provides a resource to create a cls cloud product log task
---

# tencentcloud_cls_cloud_product_log_task_v2

Provides a resource to create a cls cloud product log task

~> **NOTE:** In the destruction of resources, if cascading deletion of logset and topic is required, please set `force_delete` to `true`.

## Example Usage

### Create log delivery using the default newly created logset and topic

```hcl
resource "tencentcloud_cls_cloud_product_log_task_v2" "example" {
  instance_id          = "postgres-0an6hpv3"
  assumer_name         = "PostgreSQL"
  log_type             = "PostgreSQL-SLOW"
  cloud_product_region = "gz"
  cls_region           = "ap-guangzhou"
  logset_name          = "tf-example"
  topic_name           = "tf-example"
  force_delete         = true
}
```

### Create log delivery using existing logset and topic

```hcl
resource "tencentcloud_cls_cloud_product_log_task_v2" "example" {
  instance_id          = "postgres-0an6hpv3"
  assumer_name         = "PostgreSQL"
  log_type             = "PostgreSQL-SLOW"
  cloud_product_region = "gz"
  cls_region           = "ap-guangzhou"
  logset_id            = "ca5b4f56-1174-4eee-bc4c-69e48e0e8c45"
  topic_id             = "d8177ca9-466b-42f4-a110-5933daf0a83a"
  force_delete         = false
}
```

## Argument Reference

The following arguments are supported:

* `assumer_name` - (Required, String, ForceNew) Cloud product identification, Values: CDS, CWP, CDB, TDSQL-C, MongoDB, TDStore, DCDB, MariaDB, PostgreSQL, BH, APIS.
* `cloud_product_region` - (Required, String, ForceNew) Cloud product region. There are differences in the input format of different log types in different regions. Please refer to the following example:
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
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `log_type` - (Required, String, ForceNew) Log type, Values: CDS-AUDIT, CDS-RISK, CDB-AUDIT, TDSQL-C-AUDIT, MongoDB-AUDIT, MongoDB-SlowLog, MongoDB-ErrorLog, TDMYSQL-SLOW, DCDB-AUDIT, DCDB-SLOW, DCDB-ERROR, MariaDB-AUDIT, MariaDB-SLOW, MariaDB-ERROR, PostgreSQL-SLOW, PostgreSQL-ERROR, PostgreSQL-AUDIT, BH-FILELOG, BH-COMMANDLOG, APIS-ACCESS.
* `extend` - (Optional, String) Log configuration extension information, generally used to store additional log delivery configurations.
* `force_delete` - (Optional, Bool) Indicate whether to forcibly delete the corresponding logset and topic. If set to true, it will be forcibly deleted. Default is false.
* `logset_id` - (Optional, String, ForceNew) Log set ID.
* `logset_name` - (Optional, String, ForceNew) Log set name, required if `logset_id` is not filled in. If the log set does not exist, it will be automatically created.
* `topic_id` - (Optional, String, ForceNew) Log theme ID.
* `topic_name` - (Optional, String, ForceNew) The name of the log topic is required when `topic_id` is not filled in. If the log theme does not exist, it will be automatically created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls cloud product log task can be imported using the id, e.g.

```
terraform import tencentcloud_cls_cloud_product_log_task_v2.example postgres-1p7xvpc1#PostgreSQL#PostgreSQL-SLOW#gz
```

