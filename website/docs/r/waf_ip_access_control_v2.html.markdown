---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_ip_access_control_v2"
sidebar_current: "docs-tencentcloud-resource-waf_ip_access_control_v2"
description: |-
  Provides a resource to create a waf ip access control v2
---

# tencentcloud_waf_ip_access_control_v2

Provides a resource to create a waf ip access control v2

## Example Usage

```hcl
resource "tencentcloud_waf_ip_access_control_v2" "example" {
  instance_id = "waf_2kxtlbky11bbcr4b"
  domain      = "example.com"
  action_type = 40
  note        = "note."

  ip_list = [
    "10.0.0.10",
    "172.0.0.16",
    "192.168.0.30"
  ]

  job_type = "TimedJob"

  job_date_time {
    time_t_zone = "UTC+8"

    timed {
      end_date_time   = 0
      start_date_time = 0
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `action_type` - (Required, Int, ForceNew) 42: blocklist; 40: allowlist.
* `domain` - (Required, String, ForceNew) Specific domain name, for example, test.qcloudwaf.com.
Global domain name, that is, global.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `ip_list` - (Required, Set: [`String`]) IP parameter list.
* `job_date_time` - (Optional, List) Details of scheduled configuration.
* `job_type` - (Optional, String) Scheduled configuration type.
* `note` - (Optional, String) Remarks.

The `cron` object of `job_date_time` supports the following:

* `days` - (Optional, Set) Days in each month for execution
Note: This field may return null, indicating that no valid values can be obtained.
* `end_time` - (Optional, String) End time

Note: This field may return null, indicating that no valid values can be obtained.
* `start_time` - (Optional, String) Start time

Note: This field may return null, indicating that no valid values can be obtained.
* `w_days` - (Optional, Set) Days of each week for execution
Note: This field may return null, indicating that no valid values can be obtained.

The `job_date_time` object supports the following:

* `cron` - (Optional, List) Time parameters for periodic execution
Note: This field may return null, indicating that no valid values can be obtained.
* `time_t_zone` - (Optional, String) Time zone
Note: This field may return null, indicating that no valid values can be obtained.
* `timed` - (Optional, List) Time parameters for scheduled execution
Note: This field may return null, indicating that no valid values can be obtained.

The `timed` object of `job_date_time` supports the following:

* `end_date_time` - (Optional, Int) End timestamp, in seconds
Note: This field may return null, indicating that no valid values can be obtained.
* `start_date_time` - (Optional, Int) Start timestamp, in seconds
Note: This field may return null, indicating that no valid values can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

waf ip access control v2 can be imported using the id, e.g.

```
terraform import tencentcloud_waf_ip_access_control_v2.example waf_2kxtlbky11bbcr4b#example.com#5503616778
```

