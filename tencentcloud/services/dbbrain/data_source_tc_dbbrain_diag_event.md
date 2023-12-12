Use this data source to query detailed information of dbbrain diag_event

Example Usage

```hcl
data "tencentcloud_dbbrain_diag_history" "diag_history" {
	instance_id = "%s"
	start_time = "%s"
	end_time = "%s"
	product = "mysql"
}

data "tencentcloud_dbbrain_diag_event" "diag_event" {
  instance_id = "%s"
  event_id = data.tencentcloud_dbbrain_diag_history.diag_history.events.0.event_id
  product = "mysql"
}
```