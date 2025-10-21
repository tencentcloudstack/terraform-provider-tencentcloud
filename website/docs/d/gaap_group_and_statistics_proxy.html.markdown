---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_group_and_statistics_proxy"
sidebar_current: "docs-tencentcloud-datasource-gaap_group_and_statistics_proxy"
description: |-
  Use this data source to query detailed information of gaap and statistics proxy
---

# tencentcloud_gaap_group_and_statistics_proxy

Use this data source to query detailed information of gaap and statistics proxy

## Example Usage

```hcl
data "tencentcloud_gaap_group_and_statistics_proxy" "group_and_statistics_proxy" {
  project_id = 0
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, Int) Project Id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `group_set` - Channel group information that can be counted.
  * `group_id` - Channel Group ID.
  * `group_name` - Channel Group name.
  * `proxy_set` - Channel list in the proxy group.
    * `listener_list` - listener list.
      * `listener_id` - listener Id.
      * `listener_name` - listener name.
      * `port` - listened port.
      * `protocol` - Listener protocol type.
    * `proxy_id` - Channel Id.
    * `proxy_name` - Channel name.


