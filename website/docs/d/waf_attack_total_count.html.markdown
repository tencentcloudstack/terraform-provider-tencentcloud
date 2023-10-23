---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_attack_total_count"
sidebar_current: "docs-tencentcloud-datasource-waf_attack_total_count"
description: |-
  Use this data source to query detailed information of waf attack_total_count
---

# tencentcloud_waf_attack_total_count

Use this data source to query detailed information of waf attack_total_count

## Example Usage

### Obtain the specified domain name attack log

```hcl
data "tencentcloud_waf_attack_total_count" "example" {
  start_time   = "2023-09-01 00:00:00"
  end_time     = "2023-09-07 00:00:00"
  domain       = "domain.com"
  query_string = "method:GET"
}
```

### Obtain all domain name attack log

```hcl
data "tencentcloud_waf_attack_total_count" "example" {
  start_time   = "2023-09-01 00:00:00"
  end_time     = "2023-09-07 00:00:00"
  domain       = "all"
  query_string = "method:GET"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Query domain name, all domain use all.
* `end_time` - (Required, String) End time.
* `start_time` - (Required, String) Begin time.
* `query_string` - (Optional, String) Query conditions.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `total_count` - Total number of attacks.


