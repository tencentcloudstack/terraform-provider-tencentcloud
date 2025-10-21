---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_ddos_ip_attachment_v2"
sidebar_current: "docs-tencentcloud-resource-dayu_ddos_ip_attachment_v2"
description: |-
  Provides a resource to create a antiddos ip. Only support for bgp-multip.
---

# tencentcloud_dayu_ddos_ip_attachment_v2

Provides a resource to create a antiddos ip. Only support for bgp-multip.

## Example Usage

```hcl
resource "tencentcloud_dayu_ddos_ip_attachment_v2" "boundip" {
  id = "bgp-xxxxxx"
  bound_ip_list {
    ip          = "1.1.1.1"
    biz_type    = "public"
    instance_id = "ins-xxx"
    device_type = "cvm"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bgp_instance_id` - (Required, String, ForceNew) Anti-DDoS instance ID.
* `bound_ip_list` - (Optional, List, ForceNew) Array of IPs to bind to the Anti-DDoS instance. For Anti-DDoS Pro Single IP instance, the array contains only one IP. If there are no IPs to bind, it is empty; however, either BoundDevList or UnBoundDevList must not be empty.

The `bound_ip_list` object supports the following:

* `ip` - (Required, String) IP address.
* `biz_type` - (Optional, String) Category of product that can be bound. Valid values: public (CVM and CLB), bm (BM), eni (ENI), vpngw (VPN gateway), natgw (NAT gateway), waf (WAF), fpc (financial products), gaap (GAAP), and other (hosted IP). This field is required when you perform binding.
* `device_type` - (Optional, String) Sub-product category. Valid values: cvm (CVM), lb (Load balancer), eni (ENI), vpngw (VPN gateway), natgw (NAT gateway), waf (WAF), fpc (financial products), gaap (GAAP), eip (BM EIP) and other (managed IP). This field is required when you perform binding.
* `instance_id` - (Optional, String) Anti-DDoS instance ID of the IP. This field is required only when the instance is bound to an IP. For example, this field InstanceId will be eni-* if the instance ID is bound to an ENI IP; none if there is no instance to bind to a managed IP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



