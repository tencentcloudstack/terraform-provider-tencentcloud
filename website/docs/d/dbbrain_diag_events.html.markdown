---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_diag_events"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_diag_events"
description: |-
  Use this data source to query detailed information of DBbrain diag events
---

# tencentcloud_dbbrain_diag_events

Use this data source to query detailed information of DBbrain diag events

## Example Usage

### Query events only by time

```hcl
data "tencentcloud_dbbrain_diag_events" "example" {
  start_time = "2025-01-01T00:00:00+08:00"
  end_time   = "2026-12-31T00:00:00+08:00"
}
```

### Or add another filters

```hcl
data "tencentcloud_dbbrain_diag_events" "example" {
  start_time = "2026-01-01T00:00:00+08:00"
  end_time   = "2026-12-31T00:00:00+08:00"
  instance_ids = [
    "crs-kpyy0txj"
  ]

  product    = "redis"
  severities = [1, 2, 3, 4, 5]
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End time.
* `start_time` - (Required, String) Start time.
* `instance_ids` - (Optional, Set: [`String`]) Instance ID list.
* `product` - (Optional, String) Service product type; supported values include: `mysql` - Cloud Database MySQL, `redis` - Cloud Database Redis, `mariadb` - MariaDB database. The default is `mysql`.
* `result_output_file` - (Optional, String) Used to save results.
* `severities` - (Optional, Set: [`Int`]) Severity list, optional value is 1-fatal, 2-severity, 3-warning, 4-tips, 5-health.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - Diag event list.
  * `diag_item` - Diag item.
  * `diag_type` - Diag type.
  * `end_time` - End time.
  * `event_id` - Event ID.
  * `instance_id` - Instance ID.
  * `metric` - Metric.
  * `outline` - Outline.
  * `region` - Region.
  * `severity` - Severity.
  * `start_time` - Start time.


