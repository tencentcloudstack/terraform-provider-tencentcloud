Provides a resource to create a css stream_monitor

Example Usage

```hcl
resource "tencentcloud_css_stream_monitor" "stream_monitor" {
  ai_asr_input_index_list = [
    1,
  ]
  ai_format_diagnose = 1
  ai_ocr_input_index_list = [
    1,
  ]
  allow_monitor_report        = 1
  asr_language                = 1
  check_stream_broken         = 1
  check_stream_low_frame_rate = 1
  monitor_name                = "test"
  ocr_language                = 1

  input_list {
    input_app         = "live"
    input_domain      = "177154.push.tlivecloud.com"
    input_stream_name = "ppp"
  }

  notify_policy {
    callback_url       = "http://example.com/test"
    notify_policy_type = 1
  }

  output_info {
    output_domain        = "test122.jingxhu.top"
    output_stream_height = 1080
    output_stream_name   = "afc7847d-1fe1-43bc-b1e4-20d86303c393"
    output_stream_width  = 1920
  }
}
```

Import

css stream_monitor can be imported using the id, e.g.

```
terraform import tencentcloud_css_stream_monitor.stream_monitor stream_monitor_id
```