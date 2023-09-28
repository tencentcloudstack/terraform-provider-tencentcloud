---
subcategory: "Performance Testing Service(PTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_pts_tmp_key_generate"
sidebar_current: "docs-tencentcloud-resource-pts_tmp_key_generate"
description: |-
  Provides a resource to create a pts tmp_key
---

# tencentcloud_pts_tmp_key_generate

Provides a resource to create a pts tmp_key

## Example Usage

```hcl
resource "tencentcloud_pts_tmp_key_generate" "tmp_key" {
  project_id  = "project-1b0zqmhg"
  scenario_id = "scenario-abc"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String, ForceNew) Project ID.
* `scenario_id` - (Optional, String, ForceNew) Scenario ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `credentials` - Temporary access credentials.
  * `tmp_secret_id` - Temporary secret ID.
  * `tmp_secret_key` - Temporary secret key.
  * `token` - Temporary token.
* `expired_time` - Timestamp of temporary access credential timeout (in seconds).
* `start_time` - The timestamp of the moment when the temporary access credential was obtained (in seconds).


