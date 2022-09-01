---
subcategory: "Teo"
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
  zone_id  = tencentcloud_teo_zone.zone.id
  proxy_id = tencentcloud_teo_application_proxy.application_proxy_rule.proxy_id

  forward_client_ip = "TOA"
  origin_type       = "custom"
  origin_value = [
    "1.1.1.1:80",
  ]
  port = [
    "80",
  ]
  proto           = "TCP"
  session_persist = false
}
```

## Argument Reference

The following arguments are supported:

* `origin_type` - (Required, String) Origin server type.- custom: Specified origins.- origins: An origin group.- load_balancing: A load balancer.
* `origin_value` - (Required, Set: [`String`]) Origin server information.When OriginType is custom, this field value indicates multiple origin servers in either of the following formats:- IP:Port- Domain name:Port.When OriginType is origins, it indicates the origin group ID.
* `port` - (Required, Set: [`String`]) Valid values:- port number: `80` means port 80.- port range: `81-90` means port range 81-90.
* `proto` - (Required, String) Protocol. Valid values: `TCP`, `UDP`.
* `proxy_id` - (Required, String, ForceNew) Proxy ID.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `forward_client_ip` - (Optional, String) Passes the client IP.When Proto is TCP, valid values:- TOA: Pass the client IP via TOA.- PPV1: Pass the client IP via Proxy Protocol V1.- PPV2: Pass the client IP via Proxy Protocol V2.- OFF: Do not pass the client IP.When Proto=UDP, valid values:- PPV2: Pass the client IP via Proxy Protocol V2.- OFF: Do not pass the client IP.
* `session_persist` - (Optional, Bool) Specifies whether to enable session persistence.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.


## Import

teo application_proxy_rule can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_application_proxy_rule.application_proxy_rule zoneId#proxyId#ruleId
```

