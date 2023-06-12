---
subcategory: "VPN Connections(VPN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpn_default_health_check_ip"
sidebar_current: "docs-tencentcloud-datasource-vpn_default_health_check_ip"
description: |-
  Use this data source to query detailed information of vpn default_health_check_ip
---

# tencentcloud_vpn_default_health_check_ip

Use this data source to query detailed information of vpn default_health_check_ip

## Example Usage

```hcl
data "tencentcloud_vpn_default_health_check_ip" "default_health_check_ip" {
  vpn_gateway_id = "vpngw-gt8bianl"
}
```

## Argument Reference

The following arguments are supported:

* `vpn_gateway_id` - (Required, String) vpn gateway id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `health_check_local_ip` - local ip of health check.
* `health_check_remote_ip` - remote ip for health check.


