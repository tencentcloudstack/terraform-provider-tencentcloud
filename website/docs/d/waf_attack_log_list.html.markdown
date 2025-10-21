---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_attack_log_list"
sidebar_current: "docs-tencentcloud-datasource-waf_attack_log_list"
description: |-
  Use this data source to query detailed information of waf attack_log_list
---

# tencentcloud_waf_attack_log_list

Use this data source to query detailed information of waf attack_log_list

## Example Usage

### Obtain the specified domain name attack log list

```hcl
data "tencentcloud_waf_attack_log_list" "example" {
  domain       = "domain.com"
  start_time   = "2023-09-01 00:00:00"
  end_time     = "2023-09-07 00:00:00"
  query_string = "method:GET"
  sort         = "desc"
  query_count  = 10
  page         = 0
}
```

### Obtain all domain name attack log list

```hcl
data "tencentcloud_waf_attack_log_list" "example" {
  domain       = "all"
  start_time   = "2023-09-01 00:00:00"
  end_time     = "2023-09-07 00:00:00"
  query_string = "method:GET"
  sort         = "asc"
  query_count  = 20
  page         = 1
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain for query, all domain use all.
* `end_time` - (Required, String) End time.
* `query_string` - (Required, String) Lucene grammar.
* `start_time` - (Required, String) Begin time.
* `page` - (Optional, Int) Number of pages, starting from 0 by default.
* `query_count` - (Optional, Int) Number of queries, default to 10, maximum of 100.
* `result_output_file` - (Optional, String) Used to save results.
* `sort` - (Optional, String) Default desc, support desc, asc.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Attack log array.
  * `content` - The detail of attack log.
  * `file_name` - Useless.
  * `source` - Useless.
  * `time_stamp` - Time string.


