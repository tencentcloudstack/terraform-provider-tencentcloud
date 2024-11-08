Provides a resource to create a waf clb ip access control

Example Usage

```hcl
resource "tencentcloud_waf_clb_ip_access_control" "example" {
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

Import

waf waf_clb_ip_access_control can be imported using the id, e.g.

```
terraform import tencentcloud_waf_clb_ip_access_control.example waf_2kxtlbky11bbcr4b#example.com#5503616582
```
