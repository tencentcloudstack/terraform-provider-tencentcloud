---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_custom_rule"
sidebar_current: "docs-tencentcloud-resource-waf_custom_rule"
description: |-
  Provides a resource to create a waf custom rule
---

# tencentcloud_waf_custom_rule

Provides a resource to create a waf custom rule

-> **NOTE:** If `job_type` is `TimedJob`, Then `expire_time` must select the maximum time value of the `end_date_time` in the parameter list `timed`.

## Example Usage

### Create a standard custom rule

```hcl
resource "tencentcloud_waf_custom_rule" "example" {
  name        = "tf-example"
  sort_id     = "50"
  redirect    = "/"
  expire_time = "0"

  strategies {
    field        = "IP"
    compare_func = "ipmatch"
    content      = "2.2.2.2"
    arg          = ""
  }

  strategies {
    field        = "QUERY_STRING"
    compare_func = "rematch"
    content      = "need query string"
    arg          = ""
  }

  status      = "1"
  domain      = "test.com"
  action_type = "1"
}
```

### Create a timed resource for execution

```hcl
resource "tencentcloud_waf_custom_rule" "example" {
  name        = "tf-example"
  sort_id     = "50"
  redirect    = "/"
  expire_time = "1740672000"

  strategies {
    field        = "IP"
    compare_func = "ipmatch"
    content      = "2.2.2.2"
    arg          = ""
  }

  strategies {
    field              = "Referer"
    compare_func       = "strprefix"
    content            = "https://www.demo.com"
    arg                = ""
    case_not_sensitive = 1
  }

  status      = "1"
  domain      = "test.com"
  action_type = "1"
  job_type    = "TimedJob"
  job_date_time {
    timed {
      start_date_time = 1740585600
      end_date_time   = 1740672000
    }
    time_t_zone = "UTC+8"
  }
}
```

### Create a cron resource for execution

```hcl
resource "tencentcloud_waf_custom_rule" "example" {
  name        = "tf-example"
  sort_id     = "50"
  redirect    = "/"
  expire_time = "0"

  strategies {
    field        = "IP"
    compare_func = "ipmatch"
    content      = "2.2.2.2"
    arg          = ""
  }

  strategies {
    field              = "Referer"
    compare_func       = "strprefix"
    content            = "https://www.demo.com"
    arg                = ""
    case_not_sensitive = 1
  }

  status      = "1"
  domain      = "test.com"
  action_type = "1"
  job_type    = "CronJob"
  job_date_time {
    cron {
      w_days     = [0, 1, 2, 3, 4, 5, 6]
      start_time = "01:00:00"
      end_time   = "03:00:00"
    }
    time_t_zone = "UTC+8"
  }
}
```

## Argument Reference

The following arguments are supported:

* `action_type` - (Required, String) Action type, 1(Block), 2(Captcha), 3(log), 4(Redirect).
* `domain` - (Required, String) Domain.
* `expire_time` - (Required, String) Expiration time in second-level timestamp, for example, 1677254399 indicates the expiration time is 2023-02-24 23:59:59; 0 indicates it will never expire.
* `name` - (Required, String) Rule Name.
* `sort_id` - (Required, String) Priority, value range 0-100.
* `strategies` - (Required, List) Strategies detail.
* `job_date_time` - (Optional, List) Rule execution time.
* `job_type` - (Optional, String) Rule execution mode: TimedJob indicates scheduled execution. CronJob indicates periodic execution.
* `redirect` - (Optional, String) If the action is a Redirect, it represents the redirect address; Other situations can be left blank.
* `status` - (Optional, String) The status of the rule, 1(open), 0(close).

The `cron` object of `job_date_time` supports the following:

* `days` - (Optional, Set) Days in each month for execution. Note: This field may return null, indicating that no valid values can be obtained.
* `end_time` - (Optional, String) End time. Note: This field may return null, indicating that no valid values can be obtained.
* `start_time` - (Optional, String) Start time. Note: This field may return null, indicating that no valid values can be obtained.
* `w_days` - (Optional, Set) Days of each week for execution. Note: This field may return null, indicating that no valid values can be obtained.

The `job_date_time` object supports the following:

* `cron` - (Optional, List) Time parameters for periodic execution. Note: This field may return null, indicating that no valid values can be obtained.
* `time_t_zone` - (Optional, String) Time zone. Note: This field may return null, indicating that no valid values can be obtained.
* `timed` - (Optional, List) Time parameters for scheduled execution. Note: This field may return null, indicating that no valid values can be obtained.

The `strategies` object supports the following:

* `arg` - (Required, String) Matching parameter
              There are two types of configuration parameters: unsupported parameters and supported parameters.
              The matching parameter can be entered only when the matching field is one of the following four. Otherwise, the parameter is not supported.
                  GET (GET parameter value)		
                  POST (POST parameter value)		
                  ARGS_COOKIE (Cookie parameter value)		
                  ARGS_HEADER (Header parameter value)
          Note: This field may return null, indicating that no valid values can be obtained.
* `compare_func` - (Required, String) Logic symbol
              Logical symbols are divided into the following types:
                  empty (content is empty)
                  null (do not exist)
                  eq (equal to)
                  neq (not equal to)
                  contains (contain)
                  ncontains (do not contain)
                  strprefix (prefix matching)
                  strsuffix (suffix matching)
                  len_eq (length equals to)
                  len_gt (length is greater than)
                  len_lt (length is less than)
                  ipmatch (belong to)
                  ipnmatch (do not belong to)
                  numgt (number greater than)
                  numlt (number less than)
                  geo_in (IP geo belongs to)
                  geo_not_in (IP geo not belongs to)
                  rematch (regex match)
				  numgt (numerically greater than)
				  numlt (numerically less than)
				  numeq (numerically equal to)
				  numneq (numerically not equal to)
				  numle (numerically less than or equal to)
				  numge (numerically greater than or equal to)
              Different matching fields correspond to different logical operators. For details, see the matching field table above.
          Note: This field may return null, indicating that no valid values can be obtained.
* `content` - (Required, String) Matching content
              Currently, when the matching field is COOKIE (cookie), the matching content is not required. In other scenes, the matching content is required.
          Note: This field may return null, indicating that no valid values can be obtained.
* `field` - (Required, String) Matching field
              Different matching fields result in different matching parameters, logical operators, and matching contents. The details are as follows:
			  <table><thead><tr><th>Matching Field</th><th>Matching Parameter</th><th>Logical Symbol</th><th>Matching Content</th></tr></thead><tbody><tr><td>IP (source IP)</td><td>Parameters are not supported.</td><td>ipmatch (match)<br>ipnmatch (mismatch)</td><td>Multiple IP addresses are separated by commas. A maximum of 20 IP addresses are allowed.</td></tr><tr><td>IPv6 (source IPv6)</td><td>Parameters are not supported.</td><td>ipmatch (match)<br>ipnmatch (mismatch)</td><td>A single IPv6 address is supported.</td></tr><tr><td>Referer (referer)</td><td>Parameters are not supported.</td><td>empty (Content is empty.)<br>null (do not exist)<br>eq (equal to)<br>neq (not equal to)<br>contains (contain)<br>ncontains (do not contain)<br>len_eq (length equals to)<br>len_gt (length is greater than)<br>len_lt (length is less than)<br>strprefix (prefix matching)<br>strsuffix (suffix matching)<br>rematch (regular expression matching)</td><td>Enter the content, with a maximum of 512 characters.</td></tr><tr><td>URL (request path)</td><td>Parameters are not supported.</td><td>eq (equal to)<br>neq (not equal to)<br>contains (contain)<br>ncontains (do not contain)<br>len_eq (length equals to)<br>len_gt (length is greater than)<br>len_lt (length is less than)<br>strprefix (prefix matching)<br>strsuffix (suffix matching)<br>rematch (regular expression matching)</td><td>Enter the content starting with /, with a maximum of 512 characters.</td></tr><tr><td>UserAgent (UserAgent)</td><td>Parameters are not supported.</td><td>Same logical symbols as the matching field <font color="Red">Referer</font></td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>HTTP_METHOD (HTTP request method)</td><td>Parameters are not supported.</td><td>eq (equal to)<br>neq (not equal to)</td><td>Enter the method name. The uppercase is recommended.</td></tr><tr><td>QUERY_STRING (request string)</td><td>Parameters are not supported.</td><td>Same logical symbol as the matching field <font color="Red">Request Path</font></td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>GET (GET parameter value)</td><td>Parameter entry is supported.</td><td>contains (contain)<br>ncontains (do not contain)<br>len_eq (length equals to)<br>len_gt (length is greater than)<br>len_lt (length is less than)<br>strprefix (prefix matching)<br>strsuffix (suffix matching)</td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>GET_PARAMS_NAMES (GET parameter name)</td><td>Parameters are not supported.</td><td>exist (Parameter exists.)<br>nexist (Parameter does not exist.)<br>len_eq (length equals to)<br>len_gt (length is greater than)<br>len_lt (length is less than)<br>strprefix (prefix matching)<br>strsuffix (suffix matching)</td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>POST (POST parameter value)</td><td>Parameter entry is supported.</td><td>Same logical symbol as the matching field <font color="Red">GET Parameter Value</font></td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>GET_POST_NAMES (POST parameter name)</td><td>Parameters are not supported.</td><td>Same logical symbol as the matching field <font color="Red">GET Parameter Name</font></td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>POST_BODY (complete body)</td><td>Parameters are not supported.</td><td>Same logical symbol as the matching field <font color="Red">Request Path</font></td><td>Enter the body content with a maximum of 512 characters.</td></tr><tr><td>COOKIE (cookie)</td><td>Parameters are not supported.</td><td>empty (Content is empty.)<br>null (do not exist)<br>rematch (regular expression matching)</td><td><font color="Red">Unsupported currently</font></td></tr><tr><td>GET_COOKIES_NAMES (cookie parameter name)</td><td>Parameters are not supported.</td><td>Same logical symbol as the matching field <font color="Red">GET Parameter Name</font></td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>ARGS_COOKIE (cookie parameter value)</td><td>Parameter entry is supported.</td><td>Same logical symbol as the matching field <font color="Red">GET Parameter Value</font></td><td>Enter the contentwith a maximum of 512 characters.</td></tr><tr><td>GET_HEADERS_NAMES (header parameter name)</td><td>Parameters are not supported.</td><td>exist (Parameter exists.)<br>nexist (Parameter does not exist.)<br>len_eq (length equals to)<br>len_gt (length is greater than)<br>len_lt (length is less than)<br>strprefix (prefix matching)<br>strsuffix (suffix matching)<br>rematch (regular expression matching)</td><td>Enter the content with a maximum of 512 characters. The lowercase is recommended.</td></tr><tr><td>ARGS_HEADER (header parameter value)</td><td>Parameter entry is supported.</td><td>contains (contain)<br>ncontains (do not contain)<br>len_eq (length equals to)<br>len_gt (length is greater than)<br>len_lt (length is less than)<br>strprefix (prefix matching)<br>strsuffix (suffix matching)<br>rematch (regular expression matching)</td><td>Enter the content with a maximum of 512 characters.</td></tr><tr><td>CAPTCHA_RISK (CAPTCHA risk)</td><td>Parameters are not supported.</td><td>eq (equal to)<br>neq (not equal to)<br>belong (belongs to)<br>not_belong (does not belong to)<br>null (does not exist)<br>exist (exists)</td><td>Enter risk level value, supporting numerical range 0-255</td></tr><tr><td>CAPTCHA_DEVICE_RISK (CAPTCHA device risk)</td><td>Parameters are not supported.</td><td>eq (equal to)<br>neq (not equal to)<br>belong (belongs to)<br>not_belong (does not belong to)<br>null (does not exist)<br>exist (exists)</td><td>Enter device risk code, supporting values: 101, 201, 301, 401, 501, 601, 701</td></tr><tr><td>CAPTCHAR_SCORE (CAPTCHA risk assessment score)</td><td>Parameters are not supported.</td><td>numeq (numerically equal to)<br>numgt (numerically greater than)<br>numlt (numerically less than)<br>numle (numerically less than or equal to)<br>numge (numerically greater than or equal to)<br>null (does not exist)<br>exist (exists)</td><td>Enter assessment score, supporting numerical range 0-100</td></tr></tbody></table>
          	  Note: This field may return null, indicating that no valid values can be obtained.
* `case_not_sensitive` - (Optional, Int) 0: case-sensitive, 1: case-insensitive. Note: This field may return null, indicating that no valid values can be obtained.

The `timed` object of `job_date_time` supports the following:

* `end_date_time` - (Optional, Int) End timestamp, in seconds. Note: This field may return null, indicating that no valid values can be obtained.
* `start_date_time` - (Optional, Int) Start timestamp, in seconds. Note: This field may return null, indicating that no valid values can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - rule ID.


## Import

waf custom rule can be imported using the id, e.g.

```
terraform import tencentcloud_waf_custom_rule.example test.com#1100310609
```

