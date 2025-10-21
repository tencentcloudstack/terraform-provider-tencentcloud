---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_proxies"
sidebar_current: "docs-tencentcloud-datasource-gaap_proxies"
description: |-
  Use this data source to query gaap proxies.
---

# tencentcloud_gaap_proxies

Use this data source to query gaap proxies.

## Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

data "tencentcloud_gaap_proxies" "foo" {
  ids = [tencentcloud_gaap_proxy.foo.id]
}
```

## Argument Reference

The following arguments are supported:

* `access_region` - (Optional, String) Access region of the GAAP proxy to be queried. Conflict with `ids`.
* `ids` - (Optional, Set: [`String`]) ID of the GAAP proxy to be queried. Conflict with `project_id`, `access_region` and `realserver_region`.
* `project_id` - (Optional, Int) Project ID of the GAAP proxy to be queried. Conflict with `ids`.
* `realserver_region` - (Optional, String) Region of the GAAP realserver to be queried. Conflict with `ids`.
* `result_output_file` - (Optional, String) Used to save results.
* `tags` - (Optional, Map) Tags of the GAAP proxy to be queried. Support up to 5, display the information as long as it matches one.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `proxies` - An information list of GAAP proxy. Each element contains the following attributes:
  * `access_region` - Access region of the GAAP proxy.
  * `bandwidth` - Maximum bandwidth of the GAAP proxy, unit is Mbps.
  * `concurrent` - Maximum concurrency of the GAAP proxy, unit is 10k.
  * `create_time` - Creation time of the GAAP proxy.
  * `domain` - Access domain of the GAAP proxy.
  * `forward_ip` - Forwarding IP of the GAAP proxy.
  * `id` - ID of the GAAP proxy.
  * `ip` - Access domain of the GAAP proxy.
  * `is_auto_scale_proxy` - Indicates whether the auto scale channel is enabled, with 0 for no and 1 for yes.
  * `name` - Name of the GAAP proxy.
  * `policy_id` - Security policy ID of the GAAP proxy.
  * `project_id` - ID of the project within the GAAP proxy, '0' means is default project.
  * `realserver_region` - Region of the GAAP realserver.
  * `scalable` - Indicates whether GAAP proxy can scalable.
  * `status` - Status of the GAAP proxy.
  * `support_protocols` - Supported protocols of the GAAP proxy.
  * `tags` - Tags of the GAAP proxy.
  * `version` - Version of the GAAP proxy.


