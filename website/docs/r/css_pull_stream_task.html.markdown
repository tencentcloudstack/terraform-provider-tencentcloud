---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_pull_stream_task"
sidebar_current: "docs-tencentcloud-resource-css_pull_stream_task"
description: |-
  Provides a resource to create a css pull_stream_task
---

# tencentcloud_css_pull_stream_task

Provides a resource to create a css pull_stream_task

## Example Usage

```hcl
resource "tencentcloud_css_pull_stream_task" "pull_stream_task" {
  source_type = "source_type"
  source_urls = ["source_urls"]
  domain_name = "domain_name"
  app_name    = "app_name"
  stream_name = "stream_name"
  start_time  = "2022-11-16T22:09:28Z"
  end_time    = "2022-11-16T22:09:28Z"
  operator    = "admin"
  comment     = "comment."
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Required, String) push app name.
* `domain_name` - (Required, String) push domain name.
* `end_time` - (Required, String) task end time.
* `source_type` - (Required, String) &amp;#39;PullLivePushLive&amp;#39;: SourceUrls live type, &amp;#39;PullVodPushLive&amp;#39;: SourceUrls vod type.
* `source_urls` - (Required, Set: [`String`]) Pull Source media, SourceType=PullLivePushLive only 1 value, SourceType=PullLivePushLive can input multi values.
* `start_time` - (Required, String) task begin time.
* `stream_name` - (Required, String) push stream name.
* `backup_source_type` - (Optional, String) backup pull source type.
* `backup_source_url` - (Optional, String) backup pull source.
* `callback_events` - (Optional, Set: [`String`]) defind the callback event you need, null for all. TaskStart, TaskExit, VodSourceFileStart, VodSourceFileFinish, ResetTaskConfig, PullFileUnstable, PushStreamUnstable, PullFileFailed, PushStreamFailed, FileEndEarly.
* `callback_url` - (Optional, String) task event callback url.
* `comment` - (Optional, String) desc for pull task.
* `extra_cmd` - (Optional, String) ignore_region for ignore the input region and reblance inside the server.
* `file_index` - (Optional, Int) task enable or disable.
* `offset_time` - (Optional, Int) task enable or disable.
* `operator` - (Optional, String) desc operator user name.
* `push_args` - (Optional, String) other pushing args.
* `status` - (Optional, String) task enable or disable.
* `to_url` - (Optional, String) full target push url, DomainName, AppName, StreamName field must be empty.
* `vod_loop_times` - (Optional, Int) loop time for vod.
* `vod_refresh_type` - (Optional, String) vod refresh method. &amp;#39;ImmediateNewSource&amp;#39;: switch to new source at once, &amp;#39;ContinueBreakPoint&amp;#39;: switch to new source while old source finish.
* `watermark_list` - (Optional, List) watermark list, max 4 setting.

The `watermark_list` object supports the following:

* `height` - (Required, Int) pic height.
* `location` - (Required, Int) position type, 0:left top, 1:right top, 2:right bot, 3: left bot.
* `picture_url` - (Required, String) watermark picture url.
* `width` - (Required, Int) pic width.
* `x_position` - (Required, Int) x position.
* `y_position` - (Required, Int) y position.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_by` - desc who create the task.
* `create_time` - create time.
* `region` - task run region.
* `update_by` - desc who update the task.
* `update_time` - update time.


## Import

css pull_stream_task can be imported using the id, e.g.
```
$ terraform import tencentcloud_css_pull_stream_task.pull_stream_task pullStreamTask_id
```

