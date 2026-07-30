---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_api_sec_sensitive_custom_event_rule"
sidebar_current: "docs-tencentcloud-resource-waf_api_sec_sensitive_custom_event_rule"
description: |-
  Provides a resource to create a WAF api sec sensitive custom event rule
---

# tencentcloud_waf_api_sec_sensitive_custom_event_rule

Provides a resource to create a WAF api sec sensitive custom event rule

## Example Usage

```hcl
resource "tencentcloud_waf_api_sec_sensitive_custom_event_rule" "example" {
  domain        = "www.example.com"
  rule_name     = "tf-example"
  status        = 1
  description   = "tf example custom event rule"
  req_frequency = [10, 1]
  risk_level    = "100"
  source        = "custom"

  api_name_op {
    op    = "belong"
    value = ["/api/login"]

    api_name_method {
      api_name = "/api/login"
      method   = "POST"
    }
  }

  match_rule_list {
    key     = "get_key"
    operate = "exist"
    value   = ["admin", "root"]
  }

  stat_rule_list {
    key     = "status"
    operate = "num_gt"
    value   = ["50"]
    name    = "200"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `rule_name` - (Required, String, ForceNew) Rule name.
* `status` - (Required, Int) Rule switch, 0: off, 1: on.
* `api_name_op` - (Optional, List) API match list.
* `description` - (Optional, String) Event description.
* `match_rule_list` - (Optional, List) Match rule list.
* `req_frequency` - (Optional, List: [`Int`]) Access frequency, the first field represents the count, the second field represents the minute.
* `risk_level` - (Optional, String) Risk level, the value is 100, 200, 300, respectively representing low, medium, high risk.
* `source` - (Optional, String) Rule source.
* `stat_rule_list` - (Optional, List) Statistics rule list.

The `api_name_method` object of `api_name_op` supports the following:

* `api_name` - (Optional, String) API name.
* `method` - (Optional, String) API request method.

The `api_name_op` object supports the following:

* `api_name_method` - (Optional, List) When manually filtering, this structure should be passed.
* `op` - (Optional, String) Match method, such as belong and regex.
* `value` - (Optional, Set) Match value list.

The `match_rule_list` object supports the following:

* `key` - (Optional, String) Match field.
* `name` - (Optional, String) When the match field is get parameter value, post parameter value, cookie parameter value, header parameter value or rsp parameter value, this field can be filled.
* `operate` - (Optional, String) Operator.
* `value` - (Optional, Set) Match value.

The `stat_rule_list` object supports the following:

* `key` - (Optional, String) Match field.
* `name` - (Optional, String) When the match field is get parameter value, post parameter value, cookie parameter value, header parameter value or rsp parameter value, this field can be filled.
* `operate` - (Optional, String) Operator.
* `value` - (Optional, Set) Match value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `update_time` - Update timestamp.


## Import

WAF api sec sensitive custom event rule can be imported using the domain#ruleName, e.g.

```
terraform import tencentcloud_waf_api_sec_sensitive_custom_event_rule.example www.example.com#tf-example
```

