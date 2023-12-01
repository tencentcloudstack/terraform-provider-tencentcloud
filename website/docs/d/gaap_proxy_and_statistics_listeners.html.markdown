---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_proxy_and_statistics_listeners"
sidebar_current: "docs-tencentcloud-datasource-gaap_proxy_and_statistics_listeners"
description: |-
  Use this data source to query detailed information of gaap proxy and statistics listeners
---

# tencentcloud_gaap_proxy_and_statistics_listeners

Use this data source to query detailed information of gaap proxy and statistics listeners

## Example Usage

```hcl
data "tencentcloud_gaap_proxy_and_statistics_listeners" "proxy_and_statistics_listeners" {
  project_id = 0
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, Int) Project Id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `proxy_set` - proxy information that can be counted.
  * `listener_list` - Listener List.
    * `listener_id` - Listener Id.
    * `listener_name` - Listener Name.
    * `port` - listerned port.
    * `protocol` - Listener protocol type.
  * `proxy_id` - Proxy Id.
  * `proxy_name` - Proxy Name.


