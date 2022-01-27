---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_l7_rule_v2"
sidebar_current: "docs-tencentcloud-resource-dayu_l7_rule_v2"
description: |-
  Use this resource to create dayu new layer 7 rule
---

# tencentcloud_dayu_l7_rule_v2

Use this resource to create dayu new layer 7 rule

~> **NOTE:** This resource only support resource Anti-DDoS of type `bgpip`

## Example Usage

```hcl
resource "tencentcloud_dayu_l7_rule_v2" "tencentcloud_dayu_l7_rule_v2" {
  resource_type = "bgpip"
  resource_id   = "bgpip-000004xe"
  resource_ip   = "119.28.217.162"
  rule {
    keep_enable = false
    keeptime    = 0
    source_list {
      source = "1.2.3.5"
      weight = 100
    }
    source_list {
      source = "1.2.3.6"
      weight = 100
    }
    lb_type     = 1
    protocol    = "http"
    source_type = 2
    domain      = "github.com"
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, ForceNew) ID of the resource that the layer 7 rule works for.
* `resource_ip` - (Required, ForceNew) Ip of the resource that the layer 7 rule works for.
* `resource_type` - (Required, ForceNew) Type of the resource that the layer 7 rule works for, valid value is `bgpip`.
* `rule` - (Required) A list of layer 7 rules. Each element contains the following attributes:

The `rule` object supports the following:

* `domain` - (Required) Domain of the rule.
* `keep_enable` - (Required) session hold switch.
* `keeptime` - (Required) The keeptime of the layer 4 rule.
* `lb_type` - (Required) LB type of the rule, `1` for weight cycling and `2` for IP hash.
* `protocol` - (Required) Protocol of the rule.
* `source_list` - (Required) 
* `source_type` - (Required) Source type, `1` for source of host, `2` for source of IP.
* `cc_enable` - (Optional) HTTPS protocol CC protection status, value [0 (off), 1 (on)], defaule is 0.
* `cert_type` - (Optional) The source of the certificate must be filled in when the forwarding protocol is https, the value [2 (Tencent Cloud Hosting Certificate)], and 0 when the forwarding protocol is http.
* `https_to_http_enable` - (Optional) Whether to enable the Https protocol to use Http back-to-source, take the value [0 (off), 1 (on)], do not fill in the default is off, defaule is 0.
* `ssl_id` - (Optional) When the certificate source is a Tencent Cloud managed certificate, this field must be filled in with the managed certificate ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



