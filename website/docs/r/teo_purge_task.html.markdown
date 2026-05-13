---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_purge_task"
sidebar_current: "docs-tencentcloud-resource-teo_purge_task"
description: |-
  Provides a resource to create TEO cache purge task.
---

# tencentcloud_teo_purge_task

Provides a resource to create TEO cache purge task.

## Example Usage

### Purge URLs

```hcl
resource "tencentcloud_teo_purge_task" "purge_url_example" {
  zone_id = "zone-2qtuhspy7cr6"
  type    = "purge_url"
  targets = [
    "https://example.com/path1",
    "https://example.com/path2",
  ]
}
```

### Purge all cache

```hcl
resource "tencentcloud_teo_purge_task" "purge_all_example" {
  zone_id = "zone-2qtuhspy7cr6"
  type    = "purge_all"
}
```

### Purge cache tag

```hcl
resource "tencentcloud_teo_purge_task" "purge_cache_tag_example" {
  zone_id = "zone-2qtuhspy7cr6"
  type    = "purge_cache_tag"
  cache_tag {
    domains = [
      "example.com",
      "www.example.com",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required, String, ForceNew) Purge type. Valid values: purge_url, purge_prefix, purge_host, purge_all, purge_cache_tag.
* `zone_id` - (Required, String, ForceNew) Zone ID.
* `cache_tag` - (Optional, List, ForceNew) Cache tag configuration, required when type is purge_cache_tag.
* `method` - (Optional, String, ForceNew) Purge method. Valid values: invalidate, delete. Default: invalidate.
* `targets` - (Optional, List: [`String`], ForceNew) List of targets to purge.

The `cache_tag` object supports the following:

* `domains` - (Optional, List) Domain list for cache tag.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `job_id` - Task ID returned by CreatePurgeTask.
* `tasks` - List of purge task results.
  * `create_time` - Task creation time.
  * `fail_message` - Failure message.
  * `fail_type` - Failure type.
  * `job_id` - Task ID.
  * `method` - Purge method.
  * `status` - Task status.
  * `target` - Purge target.
  * `type` - Purge type.
  * `update_time` - Task update time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `10m`) Used when creating the resource.

