---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_custom_white_rule"
sidebar_current: "docs-tencentcloud-resource-waf_custom_white_rule"
description: |-
  Provides a resource to create a waf custom white rule
---

# tencentcloud_waf_custom_white_rule

Provides a resource to create a waf custom white rule

## Example Usage

### Create a standard custom white rule

```hcl
resource "tencentcloud_waf_custom_white_rule" "example" {
  name        = "tf-example"
  sort_id     = "30"
  expire_time = "0"

  strategies {
    field        = "IP"
    compare_func = "ipmatch"
    content      = "1.1.1.1"
    arg          = ""
  }

  status = "1"
  domain = "test.com"
  bypass = "geoip,cc,owasp"
}
```

### Create a timed resource for execution

```hcl
resource "tencentcloud_waf_custom_white_rule" "example" {
  name        = "tf-example"
  sort_id     = "30"
  expire_time = "1740672000"

  strategies {
    field        = "IP"
    compare_func = "ipmatch"
    content      = "1.1.1.1"
    arg          = ""
  }

  strategies {
    field              = "URL"
    compare_func       = "strprefix"
    content            = "/demo/path"
    arg                = ""
    case_not_sensitive = 1
  }

  status   = "1"
  domain   = "test.com"
  bypass   = "geoip,cc,owasp"
  job_type = "TimedJob"
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
resource "tencentcloud_waf_custom_white_rule" "example" {
  name        = "tf-example"
  sort_id     = "30"
  expire_time = "0"

  strategies {
    field        = "IP"
    compare_func = "ipmatch"
    content      = "1.1.1.1"
    arg          = ""
  }

  strategies {
    field              = "URL"
    compare_func       = "strprefix"
    content            = "/demo/path"
    arg                = ""
    case_not_sensitive = 1
  }

  status   = "1"
  domain   = "www.tencent.com"
  bypass   = "geoip,cc,owasp"
  job_type = "CronJob"
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

* `bypass` - (Required, String) The bypass modules are connected by commas between multiple modules. Supported modules ACL (Custom Rules), OWASP (Rule Engine), Webshell (Malicious File Detection), GeoIP (Geographic Block), BWIP (IP Black and White List), CC, BotRPC (BOT Protection), AntiLeakage (Information Leakage Prevention), API (API Security), AI (AI Engine), ip_outo_deny (IP Block), Applet (Mini Program Traffic Risk Control).
* `domain` - (Required, String) Domain name that needs to add policy.
* `expire_time` - (Required, String) Expiration time in second-level timestamp, for example, 1677254399 indicates the expiration time is 2023-02-24 23:59:59; 0 indicates it will never expire.
* `name` - (Required, String) Rule Name.
* `sort_id` - (Required, String) Priority, value range 1-100, The smaller the number, the higher the execution priority of this rule.
* `strategies` - (Required, List) Strategies detail.
* `job_date_time` - (Optional, List) Rule execution time.
* `job_type` - (Optional, String) Rule execution mode: TimedJob indicates scheduled execution. CronJob indicates periodic execution.
* `status` - (Optional, String) The status of the switch, 1 is on, 0 is off, default 1.

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
                            Different matching fields correspond to different logical operators. For details, see the matching field table above.
                        Note: This field may return null, indicating that no valid values can be obtained.
* `content` - (Required, String) Matching content
                            Currently, when the matching field is COOKIE (cookie), the matching content is not required. In other scenes, the matching content is required.
                        Note: This field may return null, indicating that no valid values can be obtained.
* `field` - (Required, String) Matching field
                            Different matching fields result in different matching parameters, logical operators, and matching contents. The details are as follows:
                        <table><thead><tr><th>Matching Field</th> <th>Matching Parameter</th> <th>Logical Symbol</th> <th>Matching Content</th></tr></thead> <tbody><tr><td>IP (source IP)</td> <td>Parameters are not supported.</td> <td>ipmatch (match)<br/>ipnmatch (mismatch)</td> <td>Multiple IP addresses are separated by commas. A maximum of 20 IP addresses are allowed.</td></tr> <tr><td>IPv6 (source IPv6)</td> <td>Parameters are not supported.</td> <td>ipmatch (match)<br/>ipnmatch (mismatch)</td> <td>A single IPv6 address is supported.</td></tr> <tr><td>Referer (referer)</td> <td>Parameters are not supported.</td> <td>empty (Content is empty.)<br/>null (do not exist)<br/>eq (equal to)<br/>neq (not equal to)<br/>contains (contain)<br/>ncontains (do not contain)<br/>len_eq (length equals to)<br/>len_gt (length is greater than)<br/>len_lt (length is less than)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td> <td>Enter the content, with a maximum of 512 characters.</td></tr> <tr><td>URL (request path)</td> <td>Parameters are not supported.</td> <td>eq (equal to)<br/>neq (not equal to)<br/>contains (contain)<br/>ncontains (do not contain)<br/>len_eq (length equals to)<br/>len_gt (length is greater than)<br/>len_lt (length is 
                         less than)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td> <td>Enter the content starting with /, with a maximum of 512 characters.</td></tr> <tr><td>UserAgent (UserAgent)</td> <td>Parameters are not supported.</td><td>Same logical symbols as the matching field <font color="Red">Referer</font></td> <td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>HTTP_METHOD (HTTP request method)</td> <td>Parameters are not supported.</td> <td>eq (equal to)<br/>neq (not equal to)</td> <td>Enter the method name. The uppercase is recommended.</td></tr> <tr><td>QUERY_STRING (request string)</td> <td>Parameters are not supported.</td> <td>Same logical symbol as the matching field <font color="Red">Request Path</font></td><td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>GET (GET parameter value)</td> <td>Parameter entry is supported.</td> <td>contains (contain)<br/>ncontains (do not contain)<br/>len_eq (length equals to)<br/>len_gt (length is greater than)<br/>len_lt (length is less than)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)</td> <td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>GET_PARAMS_NAMES (GET parameter name)</td> <td>Parameters are not supported.</td> <td>exist (Parameter exists.)<br/>nexist (Parameter does not exist.)<br/>len_eq (length equals to)<br/>len_gt (length is greater than)<br/>len_lt (length is less than)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)</td><td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>POST (POST parameter value)</td> <td>Parameter entry is supported.</td> <td>Same logical symbol as the matching field <font color="Red">GET Parameter Value</font></td> <td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>GET_POST_NAMES (POST parameter name)</td> <td>Parameters are not supported.</td> <td>Same logical symbol as the matching field <font color="Red">GET Parameter Name</font></td> <td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>POST_BODY (complete body)</td> <td>Parameters are not supported.</td> <td>Same logical symbol as the matching field <font color="Red">Request Path</font></td><td>Enter the body content with a maximum of 512 characters.</td></tr> <tr><td>COOKIE (cookie)</td> <td>Parameters are not supported.</td> <td>empty (Content is empty.)<br/>null (do not exist)<br/>rematch (regular expression matching)</td> <td><font color="Red">Unsupported currently</font></td></tr> <tr><td>GET_COOKIES_NAMES (cookie parameter name)</td> <td>Parameters are not supported.</td> <td>Same logical symbol as the matching field <font color="Red">GET Parameter Name</font></td> <td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>ARGS_COOKIE (cookie parameter value)</td> <td>Parameter entry is supported.</td> <td>Same logical symbol as the matching field <font color="Red">GET Parameter Value</font></td> <td>Enter the contentwith a maximum of 512 characters.</td></tr> <tr><td>GET_HEADERS_NAMES (header parameter name)</td> <td>Parameters are not supported.</td> <td>exist (Parameter exists.)<br/>nexist (Parameter does not exist.)<br/>len_eq (length equals to)<br/>len_gt (length is greater than)<br/>len_lt (length is less than)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td> <td>Enter the content with a maximum of 512 characters. The lowercase is recommended.</td> </tr><tr><td>ARGS_HEADER (header parameter value)</td> <td>Parameter entry is supported.</td> <td>contains (contain)<br/>ncontains (do not contain)<br/>len_eq (length equals to)<br/>len_gt (length is greater than)<br/>len_lt (length is less than)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td> <td>Enter the content with a maximum of 512 characters.</td></tr></tbody></table>
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

waf custom white rule can be imported using the id, e.g.

```
terraform import tencentcloud_waf_custom_white_rule.example test.com#1100310837
```

