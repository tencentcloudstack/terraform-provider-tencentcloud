---
subcategory: "Content Delivery Network(CDN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdn_url_purge"
sidebar_current: "docs-tencentcloud-resource-cdn_url_purge"
description: |-
  Provide a resource to invoke a Url Purge Request.
---

# tencentcloud_cdn_url_purge

Provide a resource to invoke a Url Purge Request.

## Example Usage

```hcl
resource "tencentcloud_cdn_url_purge" "foo" {
  urls = [
    "https://www.example.com/a"
  ]
}
```

### argument to request new purge task with same urls

```hcl
resource "tencentcloud_cdn_url_purge" "foo" {
  urls = [
    "https://www.example.com/a"
  ]
  redo = 1
}
```

## Argument Reference

The following arguments are supported:

* `urls` - (Required, List: [`String`], ForceNew) List of url to purge. NOTE: urls need include protocol prefix `http://` or `https://`.
* `area` - (Optional, String) Specify purge area. NOTE: only purge same area cache contents.
* `redo` - (Optional, Int) Change to purge again. NOTE: this argument only works while resource update, if set to `0` or null will not be triggered.
* `url_encode` - (Optional, Bool) Whether to encode urls, if set to `true` will auto encode instead of manual process.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `purge_history` - logs of latest purge task.
  * `create_time` - Purge task create time.
  * `flush_type` - Purge flush type of `flush` or `delete`.
  * `purge_type` - Purge category in of `url` or `path`.
  * `status` - Purge status of `fail`, `done`, `process`.
  * `task_id` - Purge task id.
  * `url` - Purge url.
* `task_id` - Task id of last operation.


