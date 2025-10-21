---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_diag_history"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_diag_history"
description: |-
  Use this data source to query detailed information of dbbrain diag_history
---

# tencentcloud_dbbrain_diag_history

Use this data source to query detailed information of dbbrain diag_history

## Example Usage

```hcl
data "tencentcloud_dbbrain_diag_history" "diag_history" {
  instance_id = "%s"
  start_time  = "%s"
  end_time    = "%s"
  product     = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End time, such as `2019-09-11 12:13:14`, the interval between the end time and the start time can be up to 2 days.
* `instance_id` - (Required, String) instance id.
* `start_time` - (Required, String) Start time, such as `2019-09-10 12:13:14`.
* `product` - (Optional, String) Service product type, supported values include: `mysql` - cloud database MySQL, `cynosdb` - cloud database CynosDB for MySQL, the default is `mysql`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `events` - Event description.
  * `diag_item` - Description of the diagnostic item.
  * `diag_type` - Diagnostic type.
  * `end_time` - End Time.
  * `event_id` - Event unique ID.
  * `instance_id` - instance id.
  * `metric` - reserved text. Note: This field may return null, indicating that no valid value can be obtained.
  * `outline` - Diagnostic summary.
  * `region` - region.
  * `severity` - severity. The severity is divided into 5 levels, according to the degree of impact from high to low: 1: Fatal, 2: Serious, 3: Warning, 4: Prompt, 5: Healthy.
  * `start_time` - start Time.


