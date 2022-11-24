---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_taw_instance"
sidebar_current: "docs-tencentcloud-datasource-rum_taw_instance"
description: |-
  Use this data source to query detailed information of rum tawInstance
---

# tencentcloud_rum_taw_instance

Use this data source to query detailed information of rum tawInstance

## Example Usage

```hcl
data "tencentcloud_rum_taw_instance" "tawInstance" {
  charge_statuses   = ""
  charge_types      = ""
  area_ids          = ""
  instance_statuses = ""
  instance_ids      = ""
}
```

## Argument Reference

The following arguments are supported:

* `area_ids` - (Optional, Set: [`Int`]) Region ID.
* `charge_statuses` - (Optional, Set: [`Int`]) Billing status.
* `charge_types` - (Optional, Set: [`Int`]) Billing type.
* `instance_ids` - (Optional, Set: [`String`]) Instance ID.
* `instance_statuses` - (Optional, Set: [`Int`]) Instance status (`1`: creating; `2`: running; `3`: exceptional; `4`: restarting; `5`: stopping; `6`: stopped; `7`: terminating; `8`: terminated).
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_set` - Instance list.
  * `area_id` - Area ID.
  * `charge_status` - Billing status (`1` = in use, `2` = expired, `3` = destroyed, `4` = assigning, `5` = failed).
  * `charge_type` - Billing type (`1` = free version, `2` = prepaid, `3` = postpaid).
  * `cluster_id` - Cluster ID.
  * `created_at` - Create time.
  * `data_retention_days` - Data retention time (days).
  * `instance_desc` - Instance Desc.
  * `instance_id` - Instance ID.
  * `instance_name` - Instance name.
  * `instance_status` - Instance status (`1` = creating, `2` = running, `3` = exception, `4` = restarting, `5` = stopping, `6` = stopped, `7` = deleted).
  * `tags` - Tag List.
    * `key` - Tag Key.
    * `value` - Tag Value.
  * `updated_at` - Update time.


