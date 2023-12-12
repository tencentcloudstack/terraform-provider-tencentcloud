Use this data source to query monitor events(There is a lot of data and it is recommended to output to a file)

Example Usage

```hcl
data "tencentcloud_monitor_product_event" "cvm_event_data" {
  start_time      = 1588700283
  is_alarm_config = 0
  product_name    = ["cvm"]
}
```