---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_rule_engine"
sidebar_current: "docs-tencentcloud-resource-teo_rule_engine"
description: |-
  Provides a resource to create a teo rule_engine
---

# tencentcloud_teo_rule_engine

Provides a resource to create a teo rule_engine

~> **NOTE:** The current resource has been deprecated, please use `tencentcloud_teo_l7_acc_rule`.

## Example Usage

```hcl
resource "tencentcloud_teo_rule_engine" "rule1" {
  zone_id   = tencentcloud_teo_zone.example.id
  rule_name = "test-rule"
  status    = "disable"

  rules {
    actions {
      normal_action {
        action = "UpstreamUrlRedirect"
        parameters {
          name = "Type"
          values = [
            "Path",
          ]
        }
        parameters {
          name = "Action"
          values = [
            "addPrefix",
          ]
        }
        parameters {
          name = "Value"
          values = [
            "/sss",
          ]
        }
      }
    }

    or {
      and {
        operator    = "equal"
        target      = "host"
        ignore_case = false
        values = [
          "a.tf-teo-t.xyz",
        ]
      }
      and {
        operator    = "equal"
        target      = "extension"
        ignore_case = false
        values = [
          "jpg",
        ]
      }
    }
    or {
      and {
        operator    = "equal"
        target      = "filename"
        ignore_case = false
        values = [
          "test.txt",
        ]
      }
    }

    sub_rules {
      tags = ["png"]
      rules {
        or {
          and {
            operator    = "notequal"
            target      = "host"
            ignore_case = false
            values = [
              "a.tf-teo-t.xyz",
            ]
          }
          and {
            operator    = "equal"
            target      = "extension"
            ignore_case = false
            values = [
              "png",
            ]
          }
        }
        or {
          and {
            operator    = "notequal"
            target      = "filename"
            ignore_case = false
            values = [
              "test.txt",
            ]
          }
        }
        actions {
          normal_action {
            action = "UpstreamUrlRedirect"
            parameters {
              name = "Type"
              values = [
                "Path",
              ]
            }
            parameters {
              name = "Action"
              values = [
                "addPrefix",
              ]
            }
            parameters {
              name = "Value"
              values = [
                "/www",
              ]
            }
          }
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `rule_name` - (Required, String) The rule name (1 to 255 characters).
* `rules` - (Required, List) Rule items list.
* `status` - (Required, String) Rule status. Values:
  - `enable`: Enabled.
  - `disable`: Disabled.
* `zone_id` - (Required, String, ForceNew) ID of the site.
* `tags` - (Optional, Set: [`String`]) rule tag list.

The `actions` object of `rules` supports the following:

* `code_action` - (Optional, List) Feature operation with a status code. Features of this type include:
  - `ErrorPage`: Custom error page.
  - `StatusCodeCache`: Status code cache TTL.
Note: This field may return null, indicating that no valid values can be obtained.
* `normal_action` - (Optional, List) Common operation. Values:
  - `AccessUrlRedirect`: Access URL rewrite.
  - `UpstreamUrlRedirect`: Origin-pull URL rewrite.
  - `QUIC`: QUIC.
  - `WebSocket`: WebSocket.
  - `VideoSeek`: Video dragging.
  - `Authentication`: Token authentication.
  - `CacheKey`: Custom cache key.
  - `Cache`: Node cache TTL.
  - `MaxAge`: Browser cache TTL.
  - `OfflineCache`: Offline cache.
  - `SmartRouting`: Smart acceleration.
  - `RangeOriginPull`: Range GETs.
  - `UpstreamHttp2`: HTTP/2 forwarding.
  - `HostHeader`: Host header rewrite.
  - `ForceRedirect`: Force HTTPS.
  - `OriginPullProtocol`: Origin-pull HTTPS.
  - `CachePrefresh`: Cache prefresh.
  - `Compression`: Smart compression.
  - `Hsts`.
  - `ClientIpHeader`.
  - `SslTlsSecureConf`.
  - `OcspStapling`.
  - `Http2`: HTTP/2 access.
  - `UpstreamFollowRedirect`: Follow origin redirect.
  - `Origin`: Origin.
Note: This field may return `null`, indicating that no valid value can be obtained.
* `rewrite_action` - (Optional, List) Feature operation with a request/response header. Features of this type include:
  - `RequestHeader`: HTTP request header modification.
  - `ResponseHeader`: HTTP response header modification.
Note: This field may return null, indicating that no valid values can be obtained.

The `and` object of `or` supports the following:

* `operator` - (Required, String) Operator. Valid values:
  - `equal`: Equal.
  - `notEqual`: Does not equal.
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
  - `client_ip`: Client IP.
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

The `and` object of `or` supports the following:

* `operator` - (Required, String) Operator. Valid values:
  - `equal`: Equal.
  - `notEqual`: Does not equal.
  - `exist`: Exists.
  - `notexist`: Does not exist.
* `target` - (Required, String) The match type. Values:
  - `filename`: File name.
  - `extension`: File extension.
  - `host`: Host.
  - `full_url`: Full URL, which indicates the complete URL path under the current site and must contain the HTTP protocol, host, and path.
  - `url`: Partial URL under the current site.  - `client_country`: Country/Region of the client.
  - `query_string`: Query string in the request URL.
  - `request_header`: HTTP request header.
  - `client_ip`: Client IP.
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

The `code_action` object of `actions` supports the following:

* `action` - (Required, String) Feature name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&!document=1) API to view the requirements for entering the feature name.
* `parameters` - (Required, List) Operation parameter.

The `normal_action` object of `actions` supports the following:

* `action` - (Required, String) Feature name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&!document=1) API to view the requirements for entering the feature name.
* `parameters` - (Required, List) Parameter.

The `or` object of `rules` supports the following:

* `and` - (Required, List) AND Conditions list of the rule. Rule would be triggered if all conditions are true.

The `or` object of `rules` supports the following:

* `and` - (Required, List) Rule engine condition. This condition will be considered met if all items in the array are met.

The `parameters` object of `code_action` supports the following:

* `name` - (Required, String) The parameter name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&!document=1) API to view the requirements for entering the parameter name.
* `status_code` - (Required, Int) The status code.
* `values` - (Required, Set) The parameter value.

The `parameters` object of `normal_action` supports the following:

* `name` - (Required, String) Parameter name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&!document=1) API to view the requirements for entering the parameter name.
* `values` - (Required, Set) The parameter value.

The `parameters` object of `rewrite_action` supports the following:

* `action` - (Required, String) Feature parameter name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&!document=1) API to view the requirements for entering the parameter name, which has three values:
  - add: Add the HTTP header.
  - set: Rewrite the HTTP header.
  - del: Delete the HTTP header.
* `name` - (Required, String) Parameter name.
* `values` - (Required, Set) Parameter value.

The `rewrite_action` object of `actions` supports the following:

* `action` - (Required, String) Feature name. You can call the [DescribeRulesSetting](https://tcloud4api.woa.com/document/product/1657/79433?!preview&!document=1) API to view the requirements for entering the feature name.
* `parameters` - (Required, List) Parameter.

The `rules` object of `sub_rules` supports the following:

* `or` - (Required, List) The condition that determines if a feature should run.
Note: If any condition in the array is met, the feature will run.
* `actions` - (Optional, List) The feature to be executed.

The `rules` object supports the following:

* `or` - (Required, List) OR Conditions list of the rule. Rule would be triggered if any of the condition is true.
* `actions` - (Optional, List) Feature to be executed.
* `sub_rules` - (Optional, List) The nested rule.

The `sub_rules` object of `rules` supports the following:

* `rules` - (Required, List) Nested rule settings.
* `tags` - (Optional, Set) Tag of the rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.
* `rule_priority` - Rule priority, the larger the value, the higher the priority, the minimum is 1.


## Import

teo rule_engine can be imported using the id#rule_id, e.g.
```
terraform import tencentcloud_teo_rule_engine.rule_engine zone-297z8rf93cfw#rule-ajol584a
```

