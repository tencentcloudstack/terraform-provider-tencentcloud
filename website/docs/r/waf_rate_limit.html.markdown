---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_rate_limit"
sidebar_current: "docs-tencentcloud-resource-waf_rate_limit"
description: |-
  Provides a resource to create a WAF rate limit rule
---

# tencentcloud_waf_rate_limit

Provides a resource to create a WAF rate limit rule

## Example Usage

### Create with API path rate limiting

```hcl
resource "tencentcloud_waf_rate_limit" "example" {
  domain         = "example.com"
  name           = "tf-example"
  priority       = 10
  status         = 1
  limit_strategy = 0
  limit_object   = "API"
  block_page     = 209057

  get_params_name {
    content = "get"
    func    = "IN"
  }

  limit_headers {
    key   = "myKey"
    type  = "IN"
    value = "myValue"
  }

  limit_paths {
    path = "/url"
    type = "EXACT"
  }

  limit_window {
    second = 0
    minute = 10
    hour   = 0
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `limit_object` - (Required, String) Rate limit object, supports API or Domain. If based on API, LimitPaths cannot be empty.
* `limit_strategy` - (Required, Int) Rate limit strategy, 0: observe, 1: block, 2: CAPTCHA.
* `limit_window` - (Required, List) Rate limit window configuration.
* `name` - (Required, String) Rule name.
* `priority` - (Required, Int) Rule priority.
* `status` - (Required, Int) Rule switch, 0: off, 1: on.
* `block_page` - (Optional, Int) Block page, 0 means 429, otherwise fill in blockPageID.
* `get_params_name` - (Optional, List) Rate limit based on GET parameter name.
* `get_params_value` - (Optional, List) Rate limit based on GET parameter value.
* `ip_location` - (Optional, List) Rate limit based on IP location.
* `limit_header_name` - (Optional, List) Rate limit based on header parameter name.
* `limit_headers` - (Optional, List) Rate limit headers configuration.
* `limit_method` - (Optional, List) Rate limit method configuration.
* `limit_paths` - (Optional, List) Rate limit path configuration.
* `object_src` - (Optional, Int) Rate limit object source, 0: manual input, 1: API asset.
* `order` - (Optional, Int) Rate limit execution order, 0: default, rate limit first, 1: security protection first.
* `paths_option` - (Optional, List) Path options, can configure request method for each path.
* `post_params_name` - (Optional, List) Rate limit based on POST parameter name.
* `post_params_value` - (Optional, List) Rate limit based on POST parameter value.
* `quota_share` - (Optional, Bool) Whether to share quota. Only valid when object is URL. false: URL exclusive quota, true: all URLs share quota.
* `redirect_info` - (Optional, List) Redirect information. Required when LimitStrategy is redirect.

The `get_params_name` object supports the following:

* `content` - (Optional, String) Match content.
* `func` - (Optional, String) Logic operator.
* `params` - (Optional, String) Match parameter.

The `get_params_value` object supports the following:

* `content` - (Optional, String) Match content.
* `func` - (Optional, String) Logic operator.
* `params` - (Optional, String) Match parameter.

The `ip_location` object supports the following:

* `content` - (Optional, String) Match content.
* `func` - (Optional, String) Logic operator.
* `params` - (Optional, String) Match parameter.

The `limit_header_name` object supports the following:

* `params_name` - (Optional, String) Parameter name.
* `type` - (Optional, String) Operator, supports REGEX, IN, NOT_IN, EACH.

The `limit_headers` object supports the following:

* `key` - (Optional, String) Header key.
* `type` - (Optional, String) Match type, supports EXACT, REGEX, IN, NOT_IN, CONTAINS, NOT_CONTAINS.
* `value` - (Optional, String) Header value.

The `limit_method` object supports the following:

* `method` - (Optional, String) Request method to rate limit.
* `type` - (Optional, String) Match type, supports EXACT, REGEX, IN, NOT_IN, CONTAINS, NOT_CONTAINS.

The `limit_paths` object supports the following:

* `path` - (Optional, String) Rate limit path.
* `type` - (Optional, String) Match type.

The `limit_window` object supports the following:

* `hour` - (Optional, Int) Maximum requests allowed per hour.
* `minute` - (Optional, Int) Maximum requests allowed per minute.
* `quota_share` - (Optional, Bool) Whether to share quota. Only valid when object is URL. false: URL exclusive quota, true: all URLs share quota.
* `second` - (Optional, Int) Maximum requests allowed per second.

The `paths_option` object supports the following:

* `method` - (Optional, String) Request method.
* `path` - (Optional, String) Request path.

The `post_params_name` object supports the following:

* `content` - (Optional, String) Match content.
* `func` - (Optional, String) Logic operator.
* `params` - (Optional, String) Match parameter.

The `post_params_value` object supports the following:

* `content` - (Optional, String) Match content.
* `func` - (Optional, String) Logic operator.
* `params` - (Optional, String) Match parameter.

The `redirect_info` object supports the following:

* `domain` - (Optional, String) Domain.
* `protocol` - (Optional, String) Protocol.
* `url` - (Optional, String) URL path.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `limit_rule_id` - Rate limit rule ID.


## Import

WAF rate limit rule can be imported using the composite id domain#limit_rule_id, e.g.

```
terraform import tencentcloud_waf_rate_limit.example example.com#4000077639
```

