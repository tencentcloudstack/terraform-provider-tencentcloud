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

~> **NOTE:** The current resource has been deprecated, please use `tencentcloud_teo_l4_proxy_rule`.

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

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateApplicationProxyRule](https://cloud.tencent.com/document/api/1552/80693) - This API is on an earlier version. If you want to call it, please switch to the latest version. For details, see [CreateL4ProxyRules] (https://intl.cloud.tencent.com/document/product/1552/103416?from_cn_redirect=1).
* Read: [DescribeApplicationProxies](https://cloud.tencent.com/document/api/1552/80690) - This API is on an earlier version. If you want to call it, please switch to the latest version. In the latest version, this API has been split into two APIs: one for querying the Layer 4 proxy instance list and the other for querying Layer ...
* Update: [ModifyApplicationProxyRule](https://cloud.tencent.com/document/api/1552/80688) - This API is on an earlier version. If you want to call it, please switch to the latest version. For details, see [ModifyL4ProxyRules] (https://intl.cloud.tencent.com/document/product/1552/103410?from_cn_redirect=1).
* Update: [ModifyApplicationProxyRuleStatus](https://cloud.tencent.com/document/api/1552/80687) - This API is on an earlier version. If you want to call it, please switch to the latest version. For details, see [ModifyL4ProxyRulesStatus
* Delete: [DeleteApplicationProxyRule](https://cloud.tencent.com/document/api/1552/80691) - This API is on an earlier version. If you want to call it, please switch to the latest version. For details, see [DeleteL4ProxyRules] (https://intl.cloud.tencent.com/document/product/1552/103414?from_cn_redirect=1).
* Delete: [ModifyApplicationProxyRuleStatus](https://cloud.tencent.com/document/api/1552/80687) - This API is on an earlier version. If you want to call it, please switch to the latest version. For details, see [ModifyL4ProxyRulesStatus

## Import

teo application_proxy_rule can be imported using the zoneId#proxyId#ruleId, e.g.
```
terraform import tencentcloud_teo_application_proxy_rule.application_proxy_rule zone-2983wizgxqvm#proxy-6972528a-373a-11ed-afca-52540044a456#rule-90b13bb4-373a-11ed-8794-525400eddfed
```

