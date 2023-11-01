---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_diagnose"
sidebar_current: "docs-tencentcloud-datasource-elasticsearch_diagnose"
description: |-
  Use this data source to query detailed information of elasticsearch diagnose
---

# tencentcloud_elasticsearch_diagnose

Use this data source to query detailed information of elasticsearch diagnose

## Example Usage

```hcl
data "tencentcloud_elasticsearch_diagnose" "diagnose" {
  instance_id = "es-xxxxxx"
  date        = "20231030"
  limit       = 1
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `date` - (Optional, String) Report date, format 20210301.
* `limit` - (Optional, Int) Number of copies returned in the report. Default value 1.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `diagnose_results` - List of diagnostic reports.
  * `completed` - Whether the diagnosis is complete or not.
  * `create_time` - Create time.
  * `instance_id` - Instance id.
  * `job_param` - Diagnostic parameters such as diagnostic time, diagnostic index, etc.
    * `indices` - Diagnostic indices.
    * `interval` - Historical diagnosis time.
    * `jobs` - Diagnostic item list.
  * `job_results` - Diagnostic item result list.
    * `advise` - Diagnostic advice.
    * `detail` - Diagnosis details.
    * `job_name` - Diagnostic item name.
    * `log_details` - Diagnostic log details.
      * `advise` - Log exception handling recommendation.
      * `count` - Number of occurrences of log exception names.
      * `key` - Log exception name.
    * `metric_details` - Details of diagnostic metrics.
      * `key` - Metric detail name.
      * `metrics` - Metric detail value.
        * `dimensions` - Index dimension family.
          * `key` - Intelligent operation and maintenance index dimension Key.
          * `value` - Dimension value of intelligent operation and maintenance index.
        * `value` - Value.
    * `score` - Diagnostic item score.
    * `setting_details` - Diagnostic configuration detail.
      * `advise` - Configuration processing recommendations.
      * `key` - Key.
      * `value` - Value.
    * `status` - Diagnostic item status:-2 failed,-1 to be retried, 0 running, 1 successful.
    * `summary` - Diagnostic summary.
  * `job_type` - Diagnosis type, 2 timing diagnosis, 3 customer manual trigger diagnosis.
  * `request_id` - Request id.
  * `score` - Total diagnostic score.


