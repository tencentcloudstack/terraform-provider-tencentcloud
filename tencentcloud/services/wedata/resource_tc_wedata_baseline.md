Provides a resource to create a wedata baseline

Example Usage

```hcl
resource "tencentcloud_wedata_baseline" "example" {
  project_id     = "1927766435649077248"
  baseline_name  = "tf_example"
  baseline_type  = "D"
  create_uin     = "100028439226"
  create_name    = "tf_user"
  in_charge_uin  = "tf_user"
  in_charge_name = "100028439226"
  promise_tasks {
    project_id          = "1927766435649077248"
    task_name           = "tf_demo_task"
    task_id             = "20231030145334153"
    task_cycle          = "D"
    workflow_name       = "dataflow_mpp"
    workflow_id         = "e4dafb2e-76eb-11ee-bfeb-b8cef68a6637"
    task_in_charge_name = ";tf_user;"
  }
  promise_time   = "00:00:00"
  warning_margin = 30
  is_new_alarm   = true
  baseline_create_alarm_rule_request {
    alarm_types = [
      "baseLineBroken",
      "baseLineWarning",
      "baseLineTaskFailure"
    ]
    alarm_level = 2
    alarm_ways  = [
      "email",
      "sms"
    ]
    alarm_recipient_type = 1
    alarm_recipients     = [
      "tf_user"
    ]
    alarm_recipient_ids = [
      "100028439226"
    ]
  }
}
```

Import

wedata baseline can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_baseline.example 1927766435649077248#2
```