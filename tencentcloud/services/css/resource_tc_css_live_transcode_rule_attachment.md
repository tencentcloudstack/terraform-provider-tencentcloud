Provides a resource to create a css live_transcode_rule_attachment

Example Usage

```hcl
resource "tencentcloud_css_pull_stream_task" "task" {
  source_type = "%s"
  source_urls = ["%s"]
  domain_name = "%s"
  app_name = "%s"
  stream_name = "%s"
  start_time = "%s"
  end_time = "%s"
  operator = "%s"
  comment = "This is a demo."
}

resource "tencentcloud_css_live_transcode_template" "temp" {
  template_name = "xxx"
  acodec = "aac"
  video_bitrate = 100
  vcodec = "origin"
  description = "This_is_a_tf_test_temp."
  need_video = 1
  need_audio = 1
}

resource "tencentcloud_css_live_transcode_rule_attachment" "live_transcode_rule_attachment" {
  domain_name = tencentcloud_css_pull_stream_task.task.domain_name
  app_name = tencentcloud_css_pull_stream_task.task.app_name
  stream_name = tencentcloud_css_pull_stream_task.task.stream_name
  template_id = tencentcloud_css_live_transcode_template.temp.id
}

```
Import

css live_transcode_rule_attachment can be imported using the id, e.g.
```
$ terraform import tencentcloud_css_live_transcode_rule_attachment.live_transcode_rule_attachment liveTranscodeRuleAttachment_id
```