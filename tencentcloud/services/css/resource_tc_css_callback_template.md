Provides a resource to create a css callback_template

Example Usage

```hcl
resource "tencentcloud_css_callback_template" "callback_template" {
  template_name              = "tf-test"
  description                = "this is demo"
  stream_begin_notify_url    = "http://www.yourdomain.com/api/notify?action=streamBegin"
  stream_end_notify_url      = "http://www.yourdomain.com/api/notify?action=streamEnd"
  record_notify_url          = "http://www.yourdomain.com/api/notify?action=record"
  snapshot_notify_url        = "http://www.yourdomain.com/api/notify?action=snapshot"
  porn_censorship_notify_url = "http://www.yourdomain.com/api/notify?action=porn"
  callback_key               = "adasda131312"
  push_exception_notify_url  = "http://www.yourdomain.com/api/notify?action=pushException"
}
```

Import

css callback_template can be imported using the id, e.g.

```
terraform import tencentcloud_css_callback_template.callback_template templateId
```