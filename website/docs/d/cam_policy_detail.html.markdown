---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_policy_detail"
sidebar_current: "docs-tencentcloud-datasource-cam_policy_detail"
description: |-
  Use this data source to query the detail of a CAM policy by its ID.
---

# tencentcloud_cam_policy_detail

Use this data source to query the detail of a CAM policy by its ID.

## Example Usage

```hcl
data "tencentcloud_cam_policy_detail" "example" {
  policy_id = 236245899
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, Int) Policy ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `policy_info` - Policy detail information.
  * `add_time` - Time the policy was created.
  * `description` - Policy description.
  * `is_service_linked_role_policy` - Whether the policy is a service-linked role policy. 0 means no, 1 means yes.
  * `policy_document` - Policy document.
  * `policy_name` - Policy name.
  * `preset_alias` - Preset policy alias. Note: this field may return null.
  * `tags` - Tags associated with the policy.
    * `key` - Tag key.
    * `value` - Tag value.
  * `type` - Policy type. 1 means custom policy, 2 means preset policy.
  * `update_time` - Time the policy was last updated.


