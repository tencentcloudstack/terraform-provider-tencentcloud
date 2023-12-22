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

* `business` - (Required, String) Business of resource instance. bgpip indicates anti-anti-ip ip; bgp means exclusive package; bgp-multip means shared packet; net indicates anti-anti-ip pro version.
* `resource_id` - (Required, String) The ID of the resource instance.
* `cc_black_white_ips` - (Optional, List) Blacklist and whitelist.
* `cc_geo_ip_policys` - (Optional, List) Details of the CC region blocking policy list.
* `cc_precision_policys` - (Optional, List) CC Precision Protection List.
* `cc_precision_req_limits` - (Optional, List) CC frequency throttling policy.
* `thresholds` - (Optional, List) List of protection threshold configurations.

The `cc_black_white_ips` object supports the following:

* `black_white_ip` - (Required, String) Blacklist and whitelist IP addresses.
* `domain` - (Required, String) Domain.
* `protocol` - (Required, String) Protocol.
* `type` - (Required, String) IP type, value [black(blacklist IP), white (whitelist IP)].
* `create_time` - (Optional, String) Create time.
* `modify_time` - (Optional, String) Modify time.

The `cc_geo_ip_policys` object supports the following:

* `action` - (Required, String) User action, drop or arg.
* `domain` - (Required, String) domain.
* `protocol` - (Required, String) Protocol, preferably HTTP, HTTPS.
* `region_type` - (Required, String) Regional types, divided into china, oversea and customized.
* `area_list` - (Optional, List) The list of region IDs that the user selects to block.
* `create_time` - (Optional, String) Create time.
* `modify_time` - (Optional, String) Modify time.

The `cc_precision_policys` object supports the following:

* `domain` - (Required, String) Domain.
* `ip` - (Required, String) Ip address.
* `policy_action` - (Required, String) Policy mode (discard or captcha).
* `policys` - (Required, List) A list of policies.
* `protocol` - (Required, String) Protocol.

The `cc_precision_req_limits` object supports the following:

* `domain` - (Required, String) Domain.
* `level` - (Required, String) Protection rating, the optional value of default means default policy, loose means loose, and strict means strict.
* `policys` - (Required, List) The CC Frequency Limit Policy Item field.
* `protocol` - (Required, String) Protocol, preferably HTTP, HTTPS.

The `policys` object of `cc_precision_policys` supports the following:

* `field_name` - (Required, String) Configuration item types, currently only support value.
* `field_type` - (Required, String) Configuration fields with the desirable values cgi, ua, cookie, referer, accept, srcip.
* `value_operator` - (Required, String) Configure the item-value comparison mode, which can be taken as the value of evaluate, not_equal, include.
* `value` - (Required, String) Configure the value.

The `policys` object of `cc_precision_req_limits` supports the following:

* `action` - (Required, String) The frequency limit policy mode, the optional value of arg indicates the verification code, and drop indicates the discard.
* `execute_duration` - (Required, Int) The duration of the frequency limit policy can be taken from 1 to 86400 per second.
* `mode` - (Required, String) The policy item is compared, and the optional value include indicates inclusion, and equal means equal.
* `period` - (Required, Int) Statistical period, take values 1, 10, 30, 60, in seconds.
* `request_num` - (Required, Int) The number of requests, the value is 1 to 20000.
* `cookie` - (Optional, String) Cookies, one of the three policy entries can only be filled in.
* `uri` - (Optional, String) Uri, one of the three policy entries can only be filled in.
* `user_agent` - (Optional, String) User-Agent, only one of the three policy entries can be filled in.

The `thresholds` object supports the following:

* `domain` - (Required, String) domain.
* `threshold` - (Required, Int) Cleaning threshold, -1 indicates that the `default` mode is turned on.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



