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

~> **NOTE:** In the destruction of resources, if cascading deletion of logset and topic is required, please set `is_delete_topic` and `is_delete_logset` to `true`.

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
  is_delete_topic      = true
  is_delete_logset     = true
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
}
```

### Create log delivery with tags bound to the associated logset and topic

```hcl
resource "tencentcloud_cls_cloud_product_log_task_v2" "example" {
  instance_id          = "postgres-0an6hpv3"
  assumer_name         = "PostgreSQL"
  log_type             = "PostgreSQL-SLOW"
  cloud_product_region = "gz"
  cls_region           = "ap-guangzhou"
  logset_name          = "tf-example"
  topic_name           = "tf-example"

  tags = {
    Environment = "production"
    Team        = "backend"
  }
}
```

## Argument Reference

The following arguments are supported:

* `assumer_name` - (Required, String, ForceNew) Cloud product identification. Supported values: APIS, BH, CDB, CDS, CFS, CLB, CSIP, CWP, DCDB, DNSPod, EMR, HTTPDNS, KHL, llmsgw, MariaDB, MDP, MongoDB, PostgreSQL, TCSS, TDSQL-C, TDStore, TencentDB-Redis, TEO, TokenHub, TSE.
* `cloud_product_region` - (Required, String, ForceNew) Cloud product region. The input format varies by log type:
- Short region code (e.g., `gz`, `sh`, `bj`): applies to APIS (all), CDB-AUDIT, TDSQL-C-AUDIT, TDMYSQL-SLOW, DCDB (all), MariaDB (all), PostgreSQL (all), MongoDB-AUDIT, TencentDB-Redis (all), EMR-OPERATION.
- Long region code (e.g., `ap-guangzhou`, `ap-shanghai`): applies to CDS (all), MongoDB-SlowLog, MongoDB-ErrorLog, MongoDB-OperationLog, DNSPod-RESOLVELOG, HTTPDNS-RESOLVELOG, MDP-SSAI, CFS-AUDIT, TEO-INEFERENCE, CSIP, TCSS, TSE, CWP, KHL.
- BH Polaris name: applies to BH (all), values: `overseas-polaris` (Hong Kong and overseas), `fsi-polaris` (finance zone), `general-polaris` (general zone), `intl-sg-prod` (international site).
* `cls_region` - (Required, String) CLS target region. Refer to the region list documentation for supported regions.
* `instance_id` - (Required, String, ForceNew) Instance ID. Obtain it from the official documentation of the corresponding cloud product.
* `log_type` - (Required, String, ForceNew) Log type, must correspond to the `assumer_name` value. Mapping:
- APIS: APIS-ACCESS
- BH: BH-COMMANDLOG, BH-FILELOG
- CDB: CDB-AUDIT
- CDS: CDS-AUDIT, CDS-RISK
- CFS: CFS-AUDIT
- CLB: CMR-SPEND
- CSIP: CSIP
- CWP: CWP
- DCDB: DCDB-AUDIT, DCDB-ERROR, DCDB-SLOW
- DNSPod: DNSPod-RESOLVELOG
- EMR: EMR-OPERATION
- HTTPDNS: HTTPDNS-RESOLVELOG
- MariaDB: MariaDB-AUDIT, MariaDB-ERROR, MariaDB-SLOW
- MDP: MDP-SSAI
- MongoDB: MongoDB-AUDIT, MongoDB-ErrorLog, MongoDB-OperationLog, MongoDB-SlowLog
- PostgreSQL: PostgreSQL-AUDIT, PostgreSQL-ERROR, PostgreSQL-SLOW
- TCSS: TCSS
- TDSQL-C: TDSQL-C-AUDIT
- TDStore: TDMYSQL-SLOW
- TencentDB-Redis: Redis-AUDIT, Redis-ERROR, Redis-SLOW
- TEO: TEO-INEFERENCE
- llmsgw: llmsgw-mcp-security-alarm.
* `extend` - (Optional, String) Log configuration extension information, generally used to store additional log delivery configurations. Example: `{"ServiceName":["HDFS","KNOX","YARN","ZOOKEEPER"],"Policy":0}`.
* `force_delete` - (Optional, Bool, **Deprecated**) It has been deprecated from version 1.82.102. Please use `is_delete_topic` or `is_delete_logset` instead. Indicate whether to forcibly delete the corresponding logset and topic. If set to true, it will be forcibly deleted. Default is false.
* `is_delete_logset` - (Optional, Bool) Whether to delete the associated Logset when deleting the log collection task. This field only takes effect when `force_delete` is false. If the Logset has other Topics, it will not be deleted. Default is false.
* `is_delete_topic` - (Optional, Bool) Whether to delete the associated Topic when deleting the log collection task. This field only takes effect when `force_delete` is false. Default is false.
* `logset_id` - (Optional, String, ForceNew) Log set ID. Obtain it via the DescribeLogsets API.
* `logset_name` - (Optional, String, ForceNew) Log set name, required when `logset_id` is not specified. If the log set does not exist, it will be created automatically.
* `tags` - (Optional, Map) Tag description list. Up to 10 tag key-value pairs are supported, and each tag key can only be bound to the same resource once. Tags are bound to the associated log topic.
* `topic_id` - (Optional, String, ForceNew) Log topic ID. Obtain it via the DescribeTopics API.
* `topic_name` - (Optional, String, ForceNew) Log topic name, required when `topic_id` is not specified. If the log topic does not exist, it will be created automatically.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls cloud product log task can be imported using the id, e.g.

```
terraform import tencentcloud_cls_cloud_product_log_task_v2.example postgres-1p7xvpc1#PostgreSQL#PostgreSQL-SLOW#gz
```

