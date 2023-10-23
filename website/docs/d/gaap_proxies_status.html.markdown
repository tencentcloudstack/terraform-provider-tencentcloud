---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_proxies_status"
sidebar_current: "docs-tencentcloud-datasource-gaap_proxies_status"
description: |-
  Use this data source to query detailed information of gaap proxies status
---

# tencentcloud_gaap_proxies_status

Use this data source to query detailed information of gaap proxies status

## Example Usage

```hcl
data "tencentcloud_gaap_proxies_status" "proxies_status" {
  proxy_ids = ["link-xxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `proxy_ids` - (Optional, Set: [`String`]) List of Proxy IDs.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_status_set` - Proxy status list.
  * `instance_id` - proxy instance ID.
  * `status` - proxy status.Among them:
- RUNNING indicates running;
- CREATING indicates being created;
- DESTROYING indicates being destroyed;
- OPENING indicates being opened;
- CLOSING indicates being closed;
- Closed indicates that it has been closed;
- ADJUSTING represents a configuration change in progress;
- ISOLATING indicates being isolated;
- ISOLATED indicates that it has been isolated;
- MOVING indicates that migration is in progress.


