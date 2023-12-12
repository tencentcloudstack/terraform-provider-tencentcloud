Use this data source to query detailed information of ses send_tasks

Example Usage

```hcl
data "tencentcloud_ses_send_tasks" "send_tasks" {
  status = 10
  receiver_id = 1063742
  task_type = 1
}
```