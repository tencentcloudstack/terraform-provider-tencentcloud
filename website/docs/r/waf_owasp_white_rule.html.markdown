---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_owasp_white_rule"
sidebar_current: "docs-tencentcloud-resource-waf_owasp_white_rule"
description: |-
  Provides a resource to create a WAF owasp white rule
---

# tencentcloud_waf_owasp_white_rule

Provides a resource to create a WAF owasp white rule

## Example Usage

```hcl
resource "tencentcloud_waf_owasp_white_rule" "example" {
  name   = "tf-example"
  domain = "example.qcloud.com"
  strategies {
    field              = "IP"
    compare_func       = "ipmatch"
    content            = "1.1.1.1"
    arg                = ""
    case_not_sensitive = 0
  }
  ids = [
    10000000,
    20000000,
    30000000,
    40000000,
    90000000,
    110000000,
    190000000,
    200000000,
    210000000,
    220000000,
    230000000,
    240000000,
    250000000,
    260000000,
    270000000,
    280000000,
    290000000,
    300000000,
    310000000,
    320000000,
    330000000,
    340000000,
    350000000,
    360000000,
    370000000
  ]
  type     = 1
  job_type = "TimedJob"
  job_date_time {
    timed {
      start_date_time = 0
      end_date_time   = 0
    }

    time_t_zone = "UTC+8"
  }
  expire_time = 0
  status      = 1
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `expire_time` - (Required, Int) If the JobDateTime field is not set, this field is used. 0 means permanent, other values indicate the cutoff time for scheduled effect (unit: seconds).
* `ids` - (Required, Set: [`Int`]) ID list of allowlisted rules.
* `job_date_time` - (Required, List) Scheduled task configuration.
* `job_type` - (Required, String) Rule execution mode: TimedJob indicates scheduled execution. CronJob indicates periodic execution.
* `name` - (Required, String) Rule name.
* `strategies` - (Required, List) Rule-Based matching policy list.
* `type` - (Required, Int) Allowlist type. valid values: 0 (allowlisting by specific rule ID), 1 (allowlisting by rule type).
* `status` - (Optional, Int) Rule status. valid values: 0 (disabled), 1 (enabled). enabled by default.

The `cron` object of `job_date_time` supports the following:

* `days` - (Optional, Set) Execution day of each month.
* `end_time` - (Optional, String) End time.
* `start_time` - (Optional, String) Start time.
* `w_days` - (Optional, Set) Execution day of each week.

The `job_date_time` object supports the following:

* `cron` - (Optional, List) Time parameter for periodic execution.
* `time_t_zone` - (Optional, String) Specifies the time zone.
* `timed` - (Optional, List) Time parameter for scheduled execution.

The `strategies` object supports the following:

* `arg` - (Required, String) Specifies the matching parameter.

Configuration parameters are divided into two data types: parameter not supported and support parameters.
When the match field is one of the following four, the matching parameter can be entered, otherwise not supported.
GET (get parameter value).		
POST (post parameter value).		
ARGS_COOKIE (COOKIE parameter value).		
ARGS_HEADER (HEADER parameter value).
* `compare_func` - (Required, String) Specifies the logic symbol. 

Logical symbols are divided into the following types:.
Empty (content is empty).
null (not found).
Eq (equal to).
neq (not equal to).
contains (contain).
ncontains (do not contain).
strprefix (prefix matching).
strsuffix (suffix matching).
Len_eq (length equals to).
Len_gt (length greater than).
Len_lt (length less than).
ipmatch (belong).
ipnmatch (not_in).
numgt (value greater than).
NumValue smaller than].
Value equal to.
numneq (value not equal to).
numle (less than or equal to).
numge (value is greater than or equal to).
geo_in (IP geographic belong).
geo_not_in (IP geographic not_in).
Specifies different logical operators for matching fields. for details, see the matching field table above.
* `content` - (Required, String) Specifies the match content.

Currently, when the match field is COOKIE (COOKIE), match content is not required. all others are needed.
* `field` - (Required, String) Specifies the matching field.

Different matching fields result in different matching parameters, logical operators, and matching contents. the details are as follows:.
<table><thead><tr><th>Matching Field</th> <th>Matching Parameter</th> <th>Logical Symbol</th> <th>Matching Content</th></tr></thead> <tbody><tr><td>IP (source IP)</td> <td>Parameters are not supported.</td> <td>ipmatch (match)<br/>ipnmatch (mismatch)</td> <td>Multiple IP addresses are separated by commas. A maximum of 20 IP addresses are allowed.</td></tr> <tr><td>IPv6 (source IPv6)</td> <td>Parameters are not supported.</td> <td>ipmatch (match)<br/>ipnmatch (mismatch)</td> <td>A single IPv6 address is supported.</td></tr> <tr><td>Referer (referer)</td> <td>Parameters are not supported.</td> <td>empty (Content is empty.)<br/>null (do not exist)<br/>eq (equal to)<br/>neq (not equal to)<br/>contains (contain)<br/>ncontains (do not contain)<br/>len_eq (length equals to)<br/>len_gt (length is greater than)<br/>len_lt (length is less than)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td> <td>Enter the content, with a maximum of 512 characters.</td></tr> <tr><td>URL (request path)</td> <td>Parameters are not supported.</td> <td>eq (equal to)<br/>neq (not equal to)<br/>contains (contain)<br/>ncontains (do not contain)<br/>len_eq (length equals to)<br/>len_gt (length is greater than)<br/>len_lt (length is 
 less than)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td> <td>Enter the content starting with /, with a maximum of 512 characters.</td></tr> <tr><td>UserAgent (UserAgent)</td> <td>Parameters are not supported.</td><td>Same logical symbols as the matching field <font color="Red">Referer</font></td> <td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>HTTP_METHOD (HTTP request method)</td> <td>Parameters are not supported.</td> <td>eq (equal to)<br/>neq (not equal to)</td> <td>Enter the method name. The uppercase is recommended.</td></tr> <tr><td>QUERY_STRING (request string)</td> <td>Parameters are not supported.</td> <td>Same logical symbol as the matching field <font color="Red">Request Path</font></td><td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>GET (GET parameter value)</td> <td>Parameter entry is supported.</td> <td>contains (contain)<br/>ncontains (do not contain)<br/>len_eq (length equals to)<br/>len_gt (length is greater than)<br/>len_lt (length is less than)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)</td> <td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>GET_PARAMS_NAMES (GET parameter name)</td> <td>Parameters are not supported.</td> <td>exist (Parameter exists.)<br/>nexist (Parameter does not exist.)<br/>len_eq (length equals to)<br/>len_gt (length is greater than)<br/>len_lt (length is less than)<br/>strprefix (prefix matching)<br/>strsuffix (suffix matching)</td><td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>POST (POST parameter value)</td> <td>Parameter entry is supported.</td> <td>Same logical symbol as the matching field <font color="Red">GET Parameter Value</font></td> <td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>GET_POST_NAMES (POST parameter name)</td> <td>Parameters are not supported.</td> <td>Same logical symbol as the matching field <font color="Red">GET Parameter Name</font></td> <td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>POST_BODY (complete body)</td> <td>Parameters are not supported.</td> <td>Same logical symbol as the matching field <font color="Red">Request Path</font></td><td>Enter the body content with a maximum of 512 characters.</td></tr> <tr><td>COOKIE (cookie)</td> <td>Parameters are not supported.</td> <td>empty (Content is empty.)<br/>null (do not exist)<br/>rematch (regular expression matching)</td> <td><font color="Red">Unsupported currently</font></td></tr> <tr><td>GET_COOKIES_NAMES (cookie parameter name)</td> <td>Parameters are not supported.</td> <td>Same logical symbol as the matching field <font color="Red">GET Parameter Name</font></td> <td>Enter the content with a maximum of 512 characters.</td></tr> <tr><td>ARGS_COOKIE (cookie parameter value)</td> <td>Parameter entry is supported.</td> <td>Same logical symbol as the matching field <font color="Red">GET Parameter Value</font></td> <td>Enter the content512 characters limit</td></tr><tr><td>GET_HEADERS_NAMES (Header parameter name)</td><td>parameter not supported</td><td>exsit (parameter exists)<br/>nexsit (parameter does not exist)<br/>len_eq (LENGTH equal)<br/>len_gt (LENGTH greater than)<br/>len_lt (LENGTH less than)<br/>strprefix (prefix match)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td><td>enter CONTENT, lowercase is recommended, up to 512 characters</td></tr><tr><td>ARGS_Header (Header parameter value)</td><td>support parameter entry</td><td>contains (include)<br/>ncontains (does not include)<br/>len_eq (LENGTH equal)<br/>len_gt (LENGTH greater than)<br/>len_lt (LENGTH less than)<br/>strprefix (prefix match)<br/>strsuffix (suffix matching)<br/>rematch (regular expression matching)</td><td>enter CONTENT, up to 512 characters</td></tr><tr><td>CONTENT_LENGTH (CONTENT-LENGTH)</td><td>support parameter entry</td><td>numgt (value greater than)<br/>numlt (value smaller than)<br/>numeq (value equal to)<br/></td><td>enter an integer between 0-9999999999999</td></tr><tr><td>IP_GEO (source IP geolocation)</td><td>support parameter entry</td><td>GEO_in (belong)<br/>GEO_not_in (not_in)<br/></td><td>enter CONTENT, up to 10240 characters, format: serialized JSON, format: [{"Country":"china","Region":"guangdong","City":"shenzhen"}]</td></tr><tr><td>CAPTCHA_RISK (CAPTCHA RISK)</td><td>parameter not supported</td><td>eq (equal)<br/>neq (not equal to)<br/>belong (belong)<br/>not_belong (not belong to)<br/>null (nonexistent)<br/>exist (exist)</td><td>enter RISK level value, value range 0-255</td></tr><tr><td>CAPTCHA_DEVICE_RISK (CAPTCHA DEVICE RISK)</td><td>parameter not supported</td><td>eq (equal)<br/>neq (not equal to)<br/>belong (belong)<br/>not_belong (not belong to)<br/>null (nonexistent)<br/>exist (exist)</td><td>enter DEVICE RISK code, valid values: 101, 201, 301, 401, 501, 601, 701</td></tr><tr><td>CAPTCHAR_SCORE (CAPTCHA RISK assessment SCORE)</td><td>parameter not supported</td><td>numeq (value equal to)<br/>numgt (value greater than)<br/>numlt (value smaller than)<br/>numle (less than or equal to)<br/>numge (value is greater than or equal to)<br/>null (nonexistent)<br/>exist (exist)</td><td>enter assessment SCORE, value range 0-100</td></tr>.
</tbody></table>.
* `case_not_sensitive` - (Optional, Int) Case-Sensitive.
Case-Insensitive.

The `timed` object of `job_date_time` supports the following:

* `end_date_time` - (Optional, Int) End timestamp, in seconds.
* `start_date_time` - (Optional, Int) Start timestamp, in seconds.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.


