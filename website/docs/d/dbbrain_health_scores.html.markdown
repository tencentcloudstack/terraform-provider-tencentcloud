---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_health_scores"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_health_scores"
description: |-
  Use this data source to query detailed information of dbbrain health_scores
---

# tencentcloud_dbbrain_health_scores

Use this data source to query detailed information of dbbrain health_scores

## Example Usage

```hcl
data "tencentcloud_dbbrain_health_scores" "health_scores" {
  instance_id = ""
  time        = ""
  product     = ""
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The ID of the instance whose health score needs to be obtained.
* `product` - (Required, String) Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database TDSQL-C for MySQL, the default is mysql.
* `time` - (Required, String) The time to obtain the health score, the time format is as follows: 2019-09-10 12:13:14.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Health score and abnormal deduction items.
  * `events_total_count` - The total number of abnormal events.
  * `health_level` - Health level, such as: HEALTH, SUB_HEALTH, RISK, HIGH_RISK.
  * `health_score` - Health score.
  * `issue_types` - Exception details.
    * `events` - unusual event.
      * `count` - Number of alerts.
      * `diag_type` - Diagnostic type.
      * `end_time` - End Time.
      * `event_id` - Event ID.
      * `metric` - reserved text.
      * `outline` - overview.
      * `score_lost` - Points deducted.
      * `severity` - severity. The severity is divided into 5 levels, according to the degree of impact from high to low: 1: Fatal, 2: Serious, 3: Warning, 4: Prompt, 5: Healthy.
      * `start_time` - Starting time.
    * `issue_type` - Index classification: AVAILABILITY: availability, MAINTAINABILITY: maintainability, PERFORMANCE, performance, RELIABILITY reliability.
    * `total_count` - The total number of abnormal events.


