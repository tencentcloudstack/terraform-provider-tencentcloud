---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_application_proxy_rule"
sidebar_current: "docs-tencentcloud-resource-teo_application_proxy_rule"
description: |-
  Provides a resource to create a teo application_proxy_rule
---

# tencentcloud_teo_application_proxy_rule

Provides a resource to create a teo application_proxy_rule

## Example Usage

```hcl
resource "tencentcloud_teo_application_proxy_rule" "application_proxy_rule" {
  forward_client_ip = "TOA"
  origin_type       = "custom"
  origin_port       = "8083"
  origin_value = [
    "127.0.0.1",
  ]
  port = [
    "8083",
  ]
  proto           = "TCP"
  proxy_id        = "proxy-6972528a-373a-11ed-afca-52540044a456"
  session_persist = false
  status          = "online"
  zone_id         = "zone-2983wizgxqvm"
}
```

## Argument Reference

The following arguments are supported:

* `origin_port` - (Required, String) Origin port, supported formats: single port: 80; Port segment: 81-90, 81 to 90 ports.
* `origin_type` - (Required, String) Origin server type. Valid values: `custom`: Specified origins; `origins`: An origin group.
* `origin_value` - (Required, Set: [`String`]) Origin site information: When `OriginType` is `custom`, it indicates one or more origin sites, such as `['8.8.8.8', '9.9.9.9']` or `OriginValue=['test.com']`; When `OriginType` is `origins`, there is required to be one and only one element, representing the origin site group ID, such as `['origin-537f5b41-162a-11ed-abaa-525400c5da15']`.
* `port` - (Required, Set: [`String`]) Valid values: `80` means port 80; `81-90` means port range 81-90.
* `proto` - (Required, String) Protocol. Valid values: `TCP`, `UDP`.
* `proxy_id` - (Required, String) Proxy ID.
* `zone_id` - (Required, String) Site ID.
* `forward_client_ip` - (Optional, String) Passes the client IP. Default value is `OFF`. When Proto is TCP, valid values: `TOA`: Pass the client IP via TOA; `PPV1`: Pass the client IP via Proxy Protocol V1; `PPV2`: Pass the client IP via Proxy Protocol V2; `OFF`: Do not pass the client IP. When Proto=UDP, valid values: `PPV2`: Pass the client IP via Proxy Protocol V2; `OFF`: Do not pass the client IP.
* `session_persist` - (Optional, Bool) Specifies whether to enable session persistence. Default value is false.
* `status` - (Optional, String) Status, the values are: `online`: enabled; `offline`: deactivated; `progress`: being deployed; `stopping`: being deactivated; `fail`: deployment failure/deactivation failure.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.


## Import

teo application_proxy_rule can be imported using the zoneId#proxyId#ruleId, e.g.
```
terraform import tencentcloud_teo_application_proxy_rule.application_proxy_rule zone-2983wizgxqvm#proxy-6972528a-373a-11ed-afca-52540044a456#rule-90b13bb4-373a-11ed-8794-525400eddfed
```

