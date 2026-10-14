---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_l7_rules_v2"
sidebar_current: "docs-tencentcloud-datasource-dayu_l7_rules_v2"
description: |-
  Use this data source to query new dayu layer 7 rules
---

# tencentcloud_dayu_l7_rules_v2

Use this data source to query new dayu layer 7 rules

## Example Usage

```hcl
data "tencentcloud_dayu_l7_rules_v2" "test" {
  business = "bgpip"
  domain   = "qq.com"
  protocol = "https"
}
```

## Argument Reference

The following arguments are supported:

* `business` - (Required, String) Type of the resource that the layer 4 rule works for, valid values are `bgpip`, `bgp`, `bgp-multip` and `net`.
* `domain` - (Optional, String) Domain of resource.
* `ip` - (Optional, String) Ip of the resource.
* `limit` - (Optional, Int, **Deprecated**) It has been deprecated from version 1.81.21. The number of pages, default is `10`.
* `offset` - (Optional, Int, **Deprecated**) It has been deprecated from version 1.81.21. The page start offset, default is `0`.
* `protocol` - (Optional, String) Protocol of resource, value range [`http`, `https`].
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of layer 4 rules. Each element contains the following attributes:
  * `cc_enable` - CC protection status of HTTPS protocol, the value is [0 (off), 1 (on)].
  * `cc_level` - CC protection level of HTTPS protocol.
  * `cc_status` - CC protection status, value [0(off), 1(on)].
  * `cc_threshold` - CC protection threshold of HTTPS protocol.
  * `cert_type` - The source of the certificate.
  * `domain` - Domain of resource.
  * `https_to_http_enable` - Whether to enable the Https protocol to use Http back-to-source, take the value [0 (off), 1 (on)], default is off.
  * `id` - Id of the resource.
  * `ip` - Ip of the resource.
  * `keep_enable` - Session keep switch, value [0 (session keep closed), 1 (session keep open)].
  * `keep_time` - Session hold time, in seconds.
  * `lb_type` - Load balancing mode, the value is [1 (weighted round-robin)].
  * `modify_time` - Modify time of resource.
  * `protocol` - Protocol of resource, value range [`http`, `https`].
  * `region` - The area code.
  * `rule_name` - Rule description.
  * `source_list` - Source list of the rule.
    * `source` - Back-to-source IP or domain name.
    * `weight` - Weight value, take value [0,100].
  * `source_type` - Back-to-origin method, value [1 (domain name back-to-source), 2 (IP back-to-source)].
  * `ssl_id` - SSL id of the resource.
  * `status` - Rule status, value [0 (rule configuration is successful), 1 (rule configuration is in effect), 2 (rule configuration fails), 3 (rule deletion is in effect), 5 (rule deletion fails), 6 (rule is waiting to be configured), 7 (rule pending deletion), 8 (rule pending configuration certificate)].
  * `virtual_port` - Virtual port of resource.


