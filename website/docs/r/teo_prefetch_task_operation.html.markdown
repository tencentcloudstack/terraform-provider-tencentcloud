---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_prefetch_task_operation"
sidebar_current: "docs-tencentcloud-resource-teo_prefetch_task_operation"
description: |-
  Provides a resource to create TEO prefetch cache task.
---

# tencentcloud_teo_prefetch_task_operation

Provides a resource to create TEO prefetch cache task.

## Example Usage

```hcl
resource "tencentcloud_teo_prefetch_task_operation" "example" {
  zone_id = "zone-12345678"
  targets = [
    "http://www.example.com/example.txt",
  ]
}
```

### Prefetch with edge mode and headers

```hcl
resource "tencentcloud_teo_prefetch_task_operation" "example" {
  zone_id = "zone-12345678"
  targets = [
    "http://www.example.com/example.txt",
  ]
  mode = "edge"
  headers {
    name  = "X-Custom-Header"
    value = "custom-value"
  }
  prefetch_media_segments = "on"
}
```

## Argument Reference

The following arguments are supported:

* `targets` - (Required, List: [`String`], ForceNew) List of URLs to prefetch. Each element format is like: http://www.example.com/example.txt.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `headers` - (Optional, List, ForceNew) HTTP headers to carry during prefetch.
* `mode` - (Optional, String, ForceNew) Prefetch mode. Valid values: `default` (prefetch to middle layer), `edge` (prefetch to edge and middle layer). Default: `default`.
* `prefetch_media_segments` - (Optional, String, ForceNew) Media segment prefetch control. Valid values: `on` (enable segment prefetch), `off` (only prefetch the submitted description file). Default: `off`.

The `headers` object supports the following:

* `name` - (Required, String) HTTP header name.
* `value` - (Required, String) HTTP header value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `job_id` - Task job ID returned by CreatePrefetchTask.
* `tasks` - Task result list.
  * `create_time` - Task creation time.
  * `fail_message` - Failure reason description.
  * `fail_type` - Failure type.
  * `job_id` - Task ID.
  * `method` - Cache purge method.
  * `status` - Task status. Valid values: processing, success, failed, timeout, canceled.
  * `target` - Resource URL.
  * `type` - Task type.
  * `update_time` - Task completion time.


