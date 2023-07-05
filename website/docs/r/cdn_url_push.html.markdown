---
subcategory: "Content Delivery Network(CDN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdn_url_push"
sidebar_current: "docs-tencentcloud-resource-cdn_url_push"
description: |-
  Provide a resource to invoke a Url Push request.
---

# tencentcloud_cdn_url_push

Provide a resource to invoke a Url Push request.

## Example Usage

```hcl
resource "tencentcloud_cdn_url_push" "foo" {
  urls = ["https://www.example.com/b"]
}
```

### argument to request new push task with same urls

```hcl
resource "tencentcloud_cdn_url_push" "foo" {
  urls = [
    "https://www.example.com/a"
  ]
  redo = 1
}
```

## Argument Reference

The following arguments are supported:

* `urls` - (Required, List: [`String`], ForceNew) List of url to push. NOTE: urls need include protocol prefix `http://` or `https://`.
* `area` - (Optional, String) Specify push area. NOTE: only push same area cache contents.
* `layer` - (Optional, String) Layer to push.
* `parse_m3u8` - (Optional, Bool) Whether to recursive parse m3u8 files.
* `redo` - (Optional, Int) Change to push again. NOTE: this argument only works while resource update, if set to `0` or null will not be triggered.
* `user_agent` - (Optional, String) Specify `User-Agent` HTTP header, default: `TencentCdn`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `push_history` - logs of latest push task.
  * `area` - Push tag area in `mainland`, `overseas` or `global`.
  * `create_time` - Push task create time.
  * `percent` - Push progress in percent.
  * `status` - Push status of `fail`, `done`, `process` or `invalid` (4xx, 5xx response).
  * `task_id` - Push task id.
  * `update_time` - Push task update time.
  * `url` - Push url.
* `task_id` - Push task id.


