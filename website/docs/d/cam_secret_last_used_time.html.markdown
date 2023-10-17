---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_secret_last_used_time"
sidebar_current: "docs-tencentcloud-datasource-cam_secret_last_used_time"
description: |-
  Use this data source to query detailed information of cam secret_last_used_time
---

# tencentcloud_cam_secret_last_used_time

Use this data source to query detailed information of cam secret_last_used_time

## Example Usage

```hcl
data "tencentcloud_cam_secret_last_used_time" "secret_last_used_time" {
  secret_id_list = ["xxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `secret_id_list` - (Required, Set: [`String`]) Query the key ID list. Supports up to 10.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `secret_id_last_used_rows` - Last used time list.
  * `last_secret_used_date` - Last used timestamp.
  * `last_used_date` - Last used date (with 1 day delay).
  * `secret_id` - Secret Id.


