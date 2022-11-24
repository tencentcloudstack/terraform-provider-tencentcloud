---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_taw_instance"
sidebar_current: "docs-tencentcloud-resource-rum_taw_instance"
description: |-
  Provides a resource to create a rum taw_instance
---

# tencentcloud_rum_taw_instance

Provides a resource to create a rum taw_instance

## Example Usage

```hcl
resource "tencentcloud_rum_taw_instance" "taw_instance" {
  area_id             = "1"
  charge_type         = "1"
  data_retention_days = "30"
  instance_name       = "instanceName-1"
  tags = {
    createdBy = "terraform"
  }
  instance_desc = "instanceDesc-1"
}
```

## Argument Reference

The following arguments are supported:

* `area_id` - (Required, Int) Region ID (at least greater than 0).
* `charge_type` - (Required, Int) Billing type (1: Pay-as-you-go).
* `data_retention_days` - (Required, Int) Data retention period (at least greater than 0).
* `instance_name` - (Required, String) Instance name (up to 255 bytes).
* `instance_desc` - (Optional, String) Instance description (up to 1,024 bytes).
* `tags` - (Optional, Map) Tag description list. Up to 10 tag key-value pairs are supported and must be unique.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `charge_status` - Billing status (`1` = in use, `2` = expired, `3` = destroyed, `4` = assigning, `5` = failed).
* `cluster_id` - Cluster ID.
* `created_at` - Create time.
* `instance_status` - Instance status (`1` = creating, `2` = running, `3` = exception, `4` = restarting, `5` = stopping, `6` = stopped, `7` = deleted).
* `updated_at` - Update time.


## Import

rum taw_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_rum_taw_instance.taw_instance tawInstance_id
```

