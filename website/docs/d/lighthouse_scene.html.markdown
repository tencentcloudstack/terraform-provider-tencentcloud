---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_scene"
sidebar_current: "docs-tencentcloud-datasource-lighthouse_scene"
description: |-
  Use this data source to query detailed information of lighthouse scene
---

# tencentcloud_lighthouse_scene

Use this data source to query detailed information of lighthouse scene

## Example Usage

```hcl
data "tencentcloud_lighthouse_scene" "scene" {
  offset = 0
  limit  = 20
}
```

## Argument Reference

The following arguments are supported:

* `limit` - (Optional, Int) Number of returned results. Default value is 20. Maximum value is 100.
* `offset` - (Optional, Int) Offset. Default value is 0.
* `result_output_file` - (Optional, String) Used to save results.
* `scene_ids` - (Optional, Set: [`String`]) List of scene IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scene_set` - List of scene info.
  * `description` - Use scene description.
  * `display_name` - Use the scene presentation name.
  * `scene_id` - Use scene Id.


