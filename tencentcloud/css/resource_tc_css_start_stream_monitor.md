Provides a resource to create a css start_stream_monitor

Example Usage

```hcl
resource "tencentcloud_css_start_stream_monitor" "start_stream_monitor" {
  monitor_id               = "3d5738dd-1ca2-4601-a6e9-004c5ec75c0b"
  audible_input_index_list = [1]
}
```

Import

css start_stream_monitor can be imported using the id, e.g.

```
terraform import tencentcloud_css_start_stream_monitor.start_stream_monitor start_stream_monitor_id
```