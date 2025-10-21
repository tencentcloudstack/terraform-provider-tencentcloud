---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_l4_proxy_rule"
sidebar_current: "docs-tencentcloud-resource-teo_l4_proxy_rule"
description: |-
  Provides a resource to create a teo teo_l4_proxy_rule
---

# tencentcloud_teo_l4_proxy_rule

Provides a resource to create a teo teo_l4_proxy_rule

## Example Usage

```hcl
resource "tencentcloud_teo_l4_proxy_rule" "teo_l4_proxy_rule" {
  proxy_id = "sid-38hbn26osico"
  zone_id  = "zone-36bjhygh1bxe"

  l4_proxy_rules {
    client_ip_pass_through_mode = "OFF"
    origin_port_range           = "1212"
    origin_type                 = "IP_DOMAIN"
    origin_value = [
      "www.aaa.com",
    ]
    port_range = [
      "1212",
    ]
    protocol             = "TCP"
    rule_tag             = "aaa"
    session_persist      = "off"
    session_persist_time = 3600
  }
}
```

## Argument Reference

The following arguments are supported:

* `l4_proxy_rules` - (Required, List) List of forwarding rules. Note: When L4ProxyRule is used here, Protocol, PortRange, OriginType, OriginValue, and OriginPortRange are required fields; ClientIPPassThroughMode, SessionPersist, SessionPersistTime, and RuleTag are optional fields; do not fill in RuleId and Status.
* `proxy_id` - (Required, String, ForceNew) Layer 4 proxy instance ID.
* `zone_id` - (Required, String, ForceNew) Zone ID.

The `l4_proxy_rules` object supports the following:

* `client_ip_pass_through_mode` - (Optional, String) Transmission of the client IP address. Valid values:<li>TOA: Available only when Protocol=TCP;</li> 
<li>PPV1: Transmission via Proxy Protocol V1. Available only when Protocol=TCP;</li>
<li>PPV2: Transmission via Proxy Protocol V2;</li> 
<li>SPP: Transmission via Simple Proxy Protocol. Available only when Protocol=UDP;</li> 
<li>OFF: No transmission.</li>
Note: This parameter is optional when L4ProxyRule is used as an input parameter in Createl4ProxyRule, and if not specified, the default value OFF will be used; it is optional when L4ProxyRule is used as an input parameter in Modifyl4ProxyRule. If not specified, it will retain its existing value.
* `origin_port_range` - (Optional, String) Origin server port, which can be set as follows:<li>A single port, such as 80;</li>
<li>A range of ports, such as 81-85, representing ports 81, 82, 83, 84, 85. When inputting a range of ports, ensure that the length corresponds with that of the forwarding port range. For example, if the forwarding port range is 80-90, this port range should be 90-100.</li>
Note: This parameter must be filled in when L4ProxyRule is used as an input parameter in Createl4ProxyRule; it is optional when L4ProxyRule is used as an input parameter in Modifyl4ProxyRule. If not specified, it will retain its existing value.
* `origin_type` - (Optional, String) Origin server type. Valid values:
<li>IP_DOMAIN: IP/Domain name origin server;</li>
<li>ORIGIN_GROUP: Origin server group;</li>
<li>LB: Cloud Load Balancer, currently only open to the allowlist.</li>
Note: This parameter must be filled in when L4ProxyRule is used as an input parameter in Createl4ProxyRule; it is optional when L4ProxyRule is used as an input parameter in Modifyl4ProxyRule. If not specified, it will retain its existing value.
* `origin_value` - (Optional, Set) Origin server address.
<li>When OriginType is set to IP_DOMAIN, enter the IP address or domain name, such as 8.8.8.8 or test.com;</li>
<li>When OriginType is set to ORIGIN_GROUP, enter the origin server group ID, such as og-537y24vf5b41;</li>
<li>When OriginType is set to LB, enter the Cloud Load Balancer instance ID, such as lb-2qwk30xf7s9g.</li>
Note: This parameter must be filled in when L4ProxyRule is used as an input parameter in Createl4ProxyRule; it is optional when L4ProxyRule is used as an input parameter in Modifyl4ProxyRule. If not specified, it will retain its existing value.
* `port_range` - (Optional, Set) Forwarding port, which can be set as follows:
<li>A single port, such as 80;</li>
<li>A range of ports, such as 81-85, representing ports 81, 82, 83, 84, 85.</li>
Note: This parameter must be filled in when L4ProxyRule is used as an input parameter in Createl4ProxyRule; it is optional when L4ProxyRule is used as an input parameter in Modifyl4ProxyRule. If not specified, it will retain its existing value.
* `protocol` - (Optional, String) Forwarding protocol. Valid values:
<li>TCP: TCP protocol;</li>
<li>UDP: UDP protocol.</li>
Note: This parameter must be filled in when L4ProxyRule is used as an input parameter in Createl4ProxyRule; it is optional when L4ProxyRule is used as an input parameter in Modifyl4ProxyRule. If not specified, it will retain its existing value.
* `rule_tag` - (Optional, String) Rule tag. Accepts 1-50 arbitrary characters.
Note: This parameter is optional when L4ProxyRule is used as an input parameter in Createl4ProxyRule; it is optional when L4ProxyRule is used as an input parameter in Modifyl4ProxyRule. If not specified, it will retain its existing value.
* `session_persist_time` - (Optional, Int) Session persistence period, with a range of 30-3600, measured in seconds.
Note: This parameter is optional when L4ProxyRule is used as an input parameter in Createl4ProxyRule. It is valid only when SessionPersist is set to on and defaults to 3600 if not specified. It is optional when L4ProxyRule is used as an input parameter in Modifyl4ProxyRule. If not specified, it will retain its existing value.
* `session_persist` - (Optional, String) Specifies whether to enable session persistence. Valid values:
<li>on: Enable;</li>
<li>off: Disable.</li>
Note: This parameter is optional when L4ProxyRule is used as an input parameter in Createl4ProxyRule, and if not specified, the default value off will be used; it is optional when L4ProxyRule is used as an input parameter in Modifyl4ProxyRule. If not specified, it will retain its existing value.
* `status` - (Optional, String) Rule status. Valid values:<li>online: Enabled;</li>
<li>offline: Disabled;</li>
<li>progress: Deploying;</li>
<li>stopping: Disabling;</li>
<li>fail: Failed to deploy or disable.</li>.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo teo_l4_proxy can be imported using the id, e.g.

```
terraform import tencentcloud_teo_l4_proxy_rule.teo_l4_proxy_rule zoneId#proxyId#ruleId
```

