---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_peak_points"
sidebar_current: "docs-tencentcloud-datasource-waf_peak_points"
description: |-
  Use this data source to query detailed information of waf peak_points
---

# tencentcloud_waf_peak_points

Use this data source to query detailed information of waf peak_points

## Example Usage

### Basic Query

```hcl
data "tencentcloud_waf_peak_points" "example" {
  from_time = "2023-09-01 00:00:00"
  to_time   = "2023-09-07 00:00:00"
}
```

### Query by filter

```hcl
data "tencentcloud_waf_peak_points" "example" {
  from_time   = "2023-09-01 00:00:00"
  to_time     = "2023-09-07 00:00:00"
  domain      = "domain.com"
  edition     = "clb-waf"
  instance_id = "waf_2kxtlbky00b2v1fn"
  metric_name = "access"
}
```

## Argument Reference

The following arguments are supported:

* `from_time` - (Required, String) Begin time.
* `to_time` - (Required, String) End time.
* `domain` - (Optional, String) The domain name to be queried. If all domain name data is queried, this parameter is not filled in.
* `edition` - (Optional, String) Only support sparta-waf and clb-waf. If not passed, there will be no filtering.
* `instance_id` - (Optional, String) WAF instance ID, if not passed, there will be no filtering.
* `metric_name` - (Optional, String) Twelve values are available: `access`-Peak qps trend chart; `botAccess`- bot peak qps trend chart; `down`-Downstream peak bandwidth trend chart; `up`-Upstream peak bandwidth trend chart; `attack`-Trend chart of total number of web attacks; `cc`-Trend chart of total number of CC attacks; `bw`- Black IP Attack Total Trend Chart; `tamper`- Anti Tamper Attack Total Trend Chart; `leak`- Trend chart of total number of anti leakage attacks; `acl`- Trend chart of total number of access control attacks; `http_status`- Trend chart of status code frequency; `wx_access`- WeChat Mini Program Peak QPS Trend Chart.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `points` - point list.
  * `access` - qps.
  * `attack` - Number of web attacks.
  * `bot_access` - Bot qps.
  * `cc` - Number of cc attacks.
  * `down` - Peak downlink bandwidth, unit B.
  * `status_client_error` - Trend chart of the number of status codes returned by WAF to the client.
  * `status_ok` - Trend chart of the number of status codes returned by WAF to the client.
  * `status_redirect` - Trend chart of the number of status codes returned by WAF to the client.
  * `status_server_error` - Trend chart of the number of status codes returned by WAF to the server.
  * `time` - Second level timestamp.
  * `up` - Peak uplink bandwidth, unit B.
  * `upstream_client_error` - Trend chart of the number of status codes returned to WAF by the origin site.
  * `upstream_redirect` - Trend chart of the number of status codes returned to WAF by the origin site.
  * `upstream_server_error` - Trend chart of the number of status codes returned to WAF by the origin site.


