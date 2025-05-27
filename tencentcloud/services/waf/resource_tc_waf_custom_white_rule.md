Provides a resource to create a WAF custom white rule

-> **NOTE:** If `job_type` is `TimedJob`, Then `expire_time` must select the maximum time value of the `end_date_time` in the parameter list `timed`.

Example Usage

Create a standard custom white rule

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
    field        = "IP_GEO"
    compare_func = "geo_in"
    content = jsonencode(
      {
        "Lang" : "cn",
        "Areas" : [
          { "Country" : "国外" }
        ]
      }
    )
    arg = ""
  }

  status     = "1"
  domain     = "www.demo.com"
  bypass     = "geoip,cc,owasp"
  logical_op = "and"
}
```

Create a timed resource for execution

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
  domain   = "www.demo.com"
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

Create a cron resource for execution

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

  status     = "1"
  domain     = "www.demo.com"
  bypass     = "geoip,cc,owasp"
  job_type   = "CronJob"
  logical_op = "or"
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

Import

WAF custom white rule can be imported using the id, e.g.

```
terraform import tencentcloud_waf_custom_white_rule.example www.demo.com#1100310837
```
