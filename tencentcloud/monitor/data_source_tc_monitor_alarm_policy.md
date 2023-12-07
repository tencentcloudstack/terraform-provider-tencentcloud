Use this data source to query detailed information of monitor alarm_policy

Example Usage

```hcl
data "tencentcloud_monitor_alarm_policy" "alarm_policy" {
  module        = "monitor"
  policy_name   = "terraform"
  monitor_types = ["MT_QCE"]
  namespaces    = ["cvm_device"]
  project_ids   = [0]
  notice_ids    = ["notice-f2svbu3w"]
  rule_types    = ["STATIC"]
  enable        = [1]
}
```