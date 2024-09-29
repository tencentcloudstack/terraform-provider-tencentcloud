---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_rule"
sidebar_current: "docs-tencentcloud-resource-teo_function_rule"
description: |-
  Provides a resource to create a teo teo_function_rule
---

# tencentcloud_teo_function_rule

Provides a resource to create a teo teo_function_rule

## Example Usage

```hcl
resource "tencentcloud_teo_function_rule" "teo_function_rule" {
  function_id = "ef-txx7fnua"
  remark      = "aaa"
  zone_id     = "zone-2qtuhspy7cr6"

  function_rule_conditions {
    rule_conditions {
      ignore_case = false
      name        = null
      operator    = "equal"
      target      = "host"
      values = [
        "aaa.makn.cn",
      ]
    }
    rule_conditions {
      ignore_case = false
      name        = null
      operator    = "equal"
      target      = "extension"
      values = [
        ".txt",
      ]
    }
  }
  function_rule_conditions {
    rule_conditions {
      ignore_case = false
      name        = null
      operator    = "notequal"
      target      = "host"
      values = [
        "aaa.makn.cn",
      ]
    }
    rule_conditions {
      ignore_case = false
      name        = null
      operator    = "equal"
      target      = "extension"
      values = [
        ".png",
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `function_id` - (Required, String, ForceNew) ID of the Function.
* `function_rule_conditions` - (Required, List) The list of rule conditions, where the conditions are connected by an "OR" relationship.
* `zone_id` - (Required, String, ForceNew) ID of the site.
* `remark` - (Optional, String) Rule description, maximum support of 60 characters.

The `function_rule_conditions` object supports the following:

* `rule_conditions` - (Required, List) For edge function trigger rule conditions, if all items in the list are satisfied, then the condition is considered fulfilled.

The `rule_conditions` object of `function_rule_conditions` supports the following:

* `operator` - (Required, String) Operator. Valid values:
  - `equals`: Equals.
  - `notEquals`: Does not equal.
  - `exist`: Exists.
  - `notexist`: Does not exist.
* `target` - (Required, String) The match type. Values:
  - `filename`: File name.
  - `extension`: File extension.
  - `host`: Host.
  - `full_url`: Full URL, which indicates the complete URL path under the current site and must contain the HTTP protocol, host, and path.
  - `url`: Partial URL under the current site.
  - `client_country`: Country/Region of the client.
  - `query_string`: Query string in the request URL.
  - `request_header`: HTTP request header.
* `ignore_case` - (Optional, Bool) Whether the parameter value is case insensitive. Default value: false.
* `name` - (Optional, String) The parameter name of the match type. This field is required only when `Target=query_string/request_header`.
  - `query_string`: Name of the query string, such as "lang" and "version" in "lang=cn&version=1".
  - `request_header`: Name of the HTTP request header, such as "Accept-Language" in the "Accept-Language:zh-CN,zh;q=0.9" header.
* `values` - (Optional, Set) The parameter value of the match type. It can be an empty string only when `Target=query string/request header` and `Operator=exist/notexist`.
  - When `Target=extension`, enter the file extension, such as "jpg" and "txt".
  - When `Target=filename`, enter the file name, such as "foo" in "foo.jpg".
  - When `Target=all`, it indicates any site request.
  - When `Target=host`, enter the host under the current site, such as "www.maxx55.com".
  - When `Target=url`, enter the partial URL path under the current site, such as "/example".
  - When `Target=full_url`, enter the complete URL under the current site. It must contain the HTTP protocol, host, and path, such as "https://www.maxx55.cn/example".
  - When `Target=client_country`, enter the ISO-3166 country/region code.
  - When `Target=query_string`, enter the value of the query string, such as "cn" and "1" in "lang=cn&version=1".
  - When `Target=request_header`, enter the HTTP request header value, such as "zh-CN,zh;q=0.9" in the "Accept-Language:zh-CN,zh;q=0.9" header.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `function_name` - The name of the function.
* `priority` - The priority of the function trigger rule. A higher numerical value indicates a higher priority.
* `rule_id` - ID of the Function Rule.


## Import

teo teo_function_rule can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function_rule.teo_function_rule zone_id#function_id#rule_id
```

