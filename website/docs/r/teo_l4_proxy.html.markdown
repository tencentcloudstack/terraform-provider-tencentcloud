---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_l4_proxy"
sidebar_current: "docs-tencentcloud-resource-teo_l4_proxy"
description: |-
  Provides a resource to create a teo teo_l4_proxy
---

# tencentcloud_teo_l4_proxy

Provides a resource to create a teo teo_l4_proxy

## Example Usage

```hcl
resource "tencentcloud_teo_l4_proxy" "proxy" {
  accelerate_mainland = "off"
  area                = "overseas"
  ipv6                = "off"
  proxy_name          = "proxy-test"
  static_ip           = "off"
  zone_id             = "zone-2qtuhspy6cr7"
}
```

## Argument Reference

The following arguments are supported:

* `proxy_name` - (Required, String) Layer 4 proxy instance name. You can enter 1-50 characters. Valid characters are a-z, 0-9, and hyphens (-). However, hyphens (-) cannot be used individually or consecutively and should not be placed at the beginning or end of the name. Modifications are not allowed after creation.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `accelerate_mainland` - (Optional, String) Specifies whether to enable network optimization in the Chinese mainland. The default value off is used if left empty. This configuration can only be enabled in certain acceleration zones and security protection configurations. For details, see [Creating an L4 Proxy Instance](https://intl.cloud.tencent.com/document/product/1552/90025?from_cn_redirect=1). Valid values: `on`: Enable; `off`: Disable.
* `area` - (Optional, String) Acceleration zone of the Layer 4 proxy instance. `mainland`: Availability zone in the Chinese mainland; `overseas`: Global availability zone (excluding the Chinese mainland); `global`: Global availability zone.
* `ddos_protection_config` - (Optional, List) Layer 3/Layer 4 DDoS protection. The default protection option of the platform will be used if it is left empty. For details, see [Exclusive DDoS Protection Usage](https://intl.cloud.tencent.com/document/product/1552/95994?from_cn_redirect=1).
* `ipv6` - (Optional, String) Specifies whether to enable IPv6 access. The default value off is used if left empty. This configuration can only be enabled in certain acceleration zones and security protection configurations. For details, see [Creating an L4 Proxy Instance](https://intl.cloud.tencent.com/document/product/1552/90025?from_cn_redirect=1). Valid values: `on`: Enable; `off`: Disable.
* `static_ip` - (Optional, String) Specifies whether to enable the fixed IP address. The default value off is used if left empty. This configuration can only be enabled in certain acceleration zones and security protection configurations. For details, see [Creating an L4 Proxy Instance](https://intl.cloud.tencent.com/document/product/1552/90025?from_cn_redirect=1). Valid values: `on`: Enable; `off`: Disable.

The `ddos_protection_config` object supports the following:

* `level_mainland` - (Optional, String) Exclusive DDoS protection specifications in the Chinese mainland. For details, see [Dedicated DDoS Mitigation Fee (Pay-as-You-Go)] (https://intl.cloud.tencent.com/document/product/1552/94162?from_cn_redirect=1). `PLATFORM`: Default protection of the platform, i.e., Exclusive DDoS protection is not enabled; `BASE30_MAX300`: Exclusive DDoS protection enabled, providing a baseline protection bandwidth of 30 Gbps and an elastic protection bandwidth of up to 300 Gbps; `BASE60_MAX600`: Exclusive DDoS protection enabled, providing a baseline protection bandwidth of 60 Gbps and an elastic protection bandwidth of up to 600 Gbps. If no parameters are filled, the default value PLATFORM is used.
* `level_overseas` - (Optional, String) Exclusive DDoS protection specifications in the worldwide region (excluding the Chinese mainland). `PLATFORM`: Default protection of the platform, i.e., Exclusive DDoS protection is not enabled; `ANYCAST300`: Exclusive DDoS protection enabled, offering a total maximum protection bandwidth of 300 Gbps; `ANYCAST_ALLIN`: Exclusive DDoS protection enabled, utilizing all available protection resources for protection. When no parameters are filled, the default value PLATFORM is used.
* `max_bandwidth_mainland` - (Optional, Int) Configuration of elastic protection bandwidth for exclusive DDoS protection in the Chinese mainland.Valid only when exclusive DDoS protection in the Chinese mainland is enabled (refer to the LevelMainland parameter configuration), and the value has the following limitations: When exclusive DDoS protection is enabled in the Chinese mainland and the 30 Gbps baseline protection bandwidth is used (the LevelMainland parameter value is BASE30_MAX300): the value range is 30 to 300 in Gbps; When exclusive DDoS protection is enabled in the Chinese mainland and the 60 Gbps baseline protection bandwidth is used (the LevelMainland parameter value is BASE60_MAX600): the value range is 60 to 600 in Gbps; When the default protection of the platform is used (the LevelMainland parameter value is PLATFORM): configuration is not supported, and the value of this parameter is invalid.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo teo_l4_proxy can be imported using the id, e.g.

```
terraform import tencentcloud_teo_l4_proxy.teo_l4_proxy teo_l4_proxy_id
```

