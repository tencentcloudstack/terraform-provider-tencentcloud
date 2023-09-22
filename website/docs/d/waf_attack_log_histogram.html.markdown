---
subcategory: "Waf"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_attack_log_histogram"
sidebar_current: "docs-tencentcloud-datasource-waf_attack_log_histogram"
description: |-
  Use this data source to query detailed information of waf attack_log_histogram
---

# tencentcloud_waf_attack_log_histogram

Use this data source to query detailed information of waf attack_log_histogram

## Example Usage

### Obtain the specified domain name log information

```hcl
data "tencentcloud_waf_attack_log_histogram" "example" {
  domain       = "domain.com"
  start_time   = "2023-09-01 00:00:00"
  end_time     = "2023-09-29 00:00:00"
  query_string = "method:GET"
}
```

### Obtain all domain name log information

```hcl
data "tencentcloud_waf_attack_log_histogram" "example" {
  domain       = "all"
  start_time   = "2023-09-01 00:00:00"
  end_time     = "2023-09-29 00:00:00"
  query_string = "method:GET"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain for query, all domain use all.
* `end_time` - (Required, String) End time.
* `query_string` - (Required, String) Lucene grammar.
* `start_time` - (Required, String) Begin time.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - The statistics detail.
  * `count` - The count of logs.
  * `time_stamp` - Timestamp.
* `period` - Period.
* `total_count` - total count.


