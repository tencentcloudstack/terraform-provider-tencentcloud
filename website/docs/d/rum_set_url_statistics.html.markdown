---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_set_url_statistics"
sidebar_current: "docs-tencentcloud-datasource-rum_set_url_statistics"
description: |-
  Use this data source to query detailed information of rum set_url_statistics
---

# tencentcloud_rum_set_url_statistics

Use this data source to query detailed information of rum set_url_statistics

## Example Usage

```hcl
data "tencentcloud_rum_set_url_statistics" "set_url_statistics" {
  start_time = 1625444040
  type       = "allcount"
  end_time   = 1625454840
  project_id = 1
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, Int) End time but is represented using a timestamp in seconds.
* `project_id` - (Required, Int) Project ID.
* `start_time` - (Required, Int) Start time but is represented using a timestamp in seconds.
* `type` - (Required, String) Query Data Type. `allcount`:CostType allcount, `data`: CostType group by data, `component`:miniProgram component, `day`:query data in day, `nettype`:query data group by nettype, `performance`:query data group by performance, `version`: CostType sort by version, `platform`: CostType sort by platform, `isp`: CostType sort by isp, `region`: CostType sort by region, `device`: CostType sort by device, `browser`: CostType sort by browser, `ext1`: CostType sort by ext1, `ext2`: CostType sort by ext2, `ext3`: CostType sort by ext3, `ret`: CostType sort by ret, `status`: CostType sort by status, `from`: CostType sort by from, `url`: CostType sort by url, `env`: CostType sort by env.
* `area` - (Optional, String) The region where the data reporting takes place.
* `brand` - (Optional, String) The mobile phone brand used for data reporting.
* `browser` - (Optional, String) The browser type used for data reporting.
* `cost_type` - (Optional, String) The method used for calculating the elapsed time `50`: 50th percentile, `75`: 75th percentile., `90`: 90th percentile., `95`: 95th percentile., `99`: 99th percentile., `99.5`: 99.5th percentile., `avg`: Mean.
* `device` - (Optional, String) The device used for data reporting.
* `engine` - (Optional, String) The browser engine used for data reporting.
* `env` - (Optional, String) The code environment where the data reporting takes place.(`production`: production env, `development`: development env, `gray`: gray env, `pre`: pre env, `daily`: daily env, `local`: local env, `others`: others env).
* `ext_first` - (Optional, String) First Expansion parameter.
* `ext_second` - (Optional, String) Second Expansion parameter.
* `ext_third` - (Optional, String) Third Expansion parameter.
* `from` - (Optional, String) The source page of the data reporting.
* `is_abroad` - (Optional, String) Whether it is non-China region.`1`: yes; `0`: no.
* `isp` - (Optional, String) The internet service provider used for data reporting.
* `level` - (Optional, String) Log level for data reporting(`1`: whitelist, `2`: normal, `4`: error, `8`: promise error, `16`: ajax request error, `32`: js resource load error, `64`: image resource load error, `128`: css resource load error, `256`: console.error, `512`: video resource load error, `1024`: request retcode error, `2048`: sdk self monitor error, `4096`: pv log, `8192`: event log).
* `net_type` - (Optional, String) The network type used for data reporting.(`1`: Wifi, `2`: 2G, `3`: 3G, `4`: 4G, `5`: 5G, `6`: 6G, `100`: Unknown).
* `os` - (Optional, String) The operating system used for data reporting.
* `package_type` - (Optional, String) Package Type.
* `platform` - (Optional, String) The platform where the data reporting takes place.(`1`: Android, `2`: IOS, `3`: Windows, `4`: Mac, `5`: Linux, `100`: Other).
* `result_output_file` - (Optional, String) Used to save results.
* `version_num` - (Optional, String) The SDK version used for data reporting.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Return value.


