---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_diag_event"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_diag_event"
description: |-
  Use this data source to query detailed information of dbbrain diag_event
---

# tencentcloud_dbbrain_diag_event

Use this data source to query detailed information of dbbrain diag_event

## Example Usage

```hcl
data "tencentcloud_dbbrain_diag_history" "diag_history" {
  instance_id = "%s"
  start_time  = "%s"
  end_time    = "%s"
  product     = "mysql"
}

data "tencentcloud_dbbrain_diag_event" "diag_event" {
  instance_id = "%s"
  event_id    = data.tencentcloud_dbbrain_diag_history.diag_history.events.0.event_id
  product     = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) isntance id.
* `event_id` - (Optional, Int) Event ID. Obtain it through `Get Instance Diagnosis History DescribeDBDiagHistory`.
* `product` - (Optional, String) Service product type, supported values include: `mysql` - cloud database MySQL, `cynosdb` - cloud database CynosDB for MySQL, the default is `mysql`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `diag_item` - diagnostic item.
* `diag_type` - Diagnostic type.
* `end_time` - End Time.
* `explanation` - Diagnostic event details, output is empty if there is no additional explanatory information.
* `metric` - reserved text. Note: This field may return null, indicating that no valid value can be obtained.
* `outline` - Diagnostic summary.
* `problem` - Diagnosed problem.
* `severity` - severity. The severity is divided into 5 levels, according to the degree of impact from high to low: 1: Fatal, 2: Serious, 3: Warning, 4: Prompt, 5: Healthy.
* `start_time` - Starting time.
* `suggestions` - A diagnostic suggestion, or empty if there is no suggestion.


