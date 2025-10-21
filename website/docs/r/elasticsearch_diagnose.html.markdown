---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_diagnose"
sidebar_current: "docs-tencentcloud-resource-elasticsearch_diagnose"
description: |-
  Provides a resource to create a elasticsearch diagnose
---

# tencentcloud_elasticsearch_diagnose

Provides a resource to create a elasticsearch diagnose

## Example Usage

```hcl
resource "tencentcloud_elasticsearch_diagnose" "diagnose" {
  instance_id = "es-xxxxxx"
  cron_time   = "15:00:00"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `cron_time` - (Optional, String) Intelligent operation and maintenance staff regularly patrol the inspection time every day, the time format is HH:00:00, such as 15:00:00.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `diagnose_job_metas` - Diagnostic items and meta-information of intelligent operation and maintenance.
  * `job_description` - Intelligent operation and maintenance diagnostic item description.
  * `job_name` - English name of diagnosis item for intelligent operation and maintenance.
  * `job_zh_name` - Chinese name of intelligent operation and maintenance diagnosis item.
* `max_count` - The maximum number of manual triggers per day for intelligent operation and maintenance staff.


## Import

es diagnose can be imported using the id, e.g.

```
terraform import tencentcloud_elasticsearch_diagnose.diagnose diagnose_id
```

