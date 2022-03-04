---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_cc_policy_v2"
sidebar_current: "docs-tencentcloud-resource-dayu_cc_policy_v2"
description: |-
  Use this resource to create a dayu CC policy
---

# tencentcloud_dayu_cc_policy_v2

Use this resource to create a dayu CC policy

## Example Usage

```hcl
resource "tencentcloud_dayu_cc_policy_v2" "demo" {
  resource_id = "bgpip-000004xf"
  business    = "bgpip"
  thresholds {
    domain    = "12.com"
    threshold = 0
  }
  cc_geo_ip_policys {
    action      = "drop"
    region_type = "china"
    domain      = "12.com"
    protocol    = "http"
  }

  cc_black_white_ips {
    protocol       = "http"
    domain         = "12.com"
    black_white_ip = "1.2.3.4"
    type           = "black"
  }
  cc_precision_policys {
    policy_action = "drop"
    domain        = "1.com"
    protocol      = "http"
    ip            = "162.62.163.34"
    policys {
      field_name     = "cgi"
      field_type     = "value"
      value          = "12123.com"
      value_operator = "equal"
    }
  }
  cc_precision_req_limits {
    domain   = "11.com"
    protocol = "http"
    level    = "loose"
    policys {
      action           = "alg"
      execute_duration = 2
      mode             = "equal"
      period           = 5
      request_num      = 12
      uri              = "15.com"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `business` - (Required) Bussiness of resource instance. bgpip indicates anti-anti-ip ip; bgp means exclusive package; bgp-multip means shared packet; net indicates anti-anti-ip pro version.
* `resource_id` - (Required) The ID of the resource instance.
* `cc_black_white_ips` - (Optional) Blacklist and whitelist.
* `cc_geo_ip_policys` - (Optional) Details of the CC region blocking policy list.
* `cc_precision_policys` - (Optional) CC Precision Protection List.
* `cc_precision_req_limits` - (Optional) CC frequency throttling policy.
* `thresholds` - (Optional) List of protection threshold configurations.

The `cc_black_white_ips` object supports the following:

* `black_white_ip` - (Required) Blacklist and whitelist IP addresses.
* `domain` - (Required) Domain.
* `protocol` - (Required) Protocol.
* `type` - (Required) IP type, value [black(blacklist IP), white (whitelist IP)].
* `create_time` - (Optional) Create time.
* `modify_time` - (Optional) Modify time.

The `cc_geo_ip_policys` object supports the following:

* `action` - (Required) User action, drop or arg.
* `domain` - (Required) domain.
* `protocol` - (Required) Protocol, preferably HTTP, HTTPS.
* `region_type` - (Required) Regional types, divided into china, oversea and customized.
* `area_list` - (Optional) The list of region IDs that the user selects to block.
* `create_time` - (Optional) Create time.
* `modify_time` - (Optional) Modify time.

The `cc_precision_policys` object supports the following:

* `domain` - (Required) Domain.
* `ip` - (Required) Ip address.
* `policy_action` - (Required) Policy mode (discard or captcha).
* `policys` - (Required) A list of policies.
* `protocol` - (Required) Protocol.

The `cc_precision_req_limits` object supports the following:

* `domain` - (Required) Domain.
* `level` - (Required) Protection rating, the optional value of default means default policy, loose means loose, and strict means strict.
* `policys` - (Required) The CC Frequency Limit Policy Item field.
* `protocol` - (Required) Protocol, preferably HTTP, HTTPS.

The `policys` object supports the following:

* `action` - (Required) The frequency limit policy mode, the optional value of arg indicates the verification code, and drop indicates the discard.
* `execute_duration` - (Required) The duration of the frequency limit policy can be taken from 1 to 86400 per second.
* `mode` - (Required) The policy item is compared, and the optional value include indicates inclusion, and equal means equal.
* `period` - (Required) Statistical period, take values 1, 10, 30, 60, in seconds.
* `request_num` - (Required) The number of requests, the value is 1 to 20000.
* `cookie` - (Optional) Cookies, one of the three policy entries can only be filled in.
* `uri` - (Optional) Uri, one of the three policy entries can only be filled in.
* `user_agent` - (Optional) User-Agent, only one of the three policy entries can be filled in.

The `policys` object supports the following:

* `field_name` - (Required) Configuration item types, currently only support value.
* `field_type` - (Required) Configuration fields with the desirable values cgi, ua, cookie, referer, accept, srcip.
* `value_operator` - (Required) Configure the item-value comparison mode, which can be taken as the value of evaluate, not_equal, include.
* `value` - (Required) Configure the value.

The `thresholds` object supports the following:

* `domain` - (Required) domain.
* `threshold` - (Required) Cleaning threshold, -1 indicates that the `default` mode is turned on.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



