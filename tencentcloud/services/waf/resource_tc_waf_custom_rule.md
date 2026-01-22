Provides a resource to create a waf custom rule

-> **NOTE:** If `job_type` is `TimedJob`, Then `expire_time` must select the maximum time value of the `end_date_time` in the parameter list `timed`.

Example Usage

Create a standard custom rule

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

Create a timed resource for execution

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

Create a cron resource for execution

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
      w_days = [0, 1, 2, 3, 4, 5, 6]
      start_time = "01:00:00"
      end_time = "03:00:00"
    }
    time_t_zone = "UTC+8"
  }
}
```

Import

waf custom rule can be imported using the id, e.g.

```
terraform import tencentcloud_waf_custom_rule.example test.com#1100310609
```