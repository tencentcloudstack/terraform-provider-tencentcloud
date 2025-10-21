---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_policies"
sidebar_current: "docs-tencentcloud-datasource-cam_policies"
description: |-
  Use this data source to query detailed information of CAM policies
---

# tencentcloud_cam_policies

Use this data source to query detailed information of CAM policies

## Example Usage

```hcl
# query by policy_id
data "tencentcloud_cam_policies" "foo" {
  policy_id = tencentcloud_cam_policy.foo.id
}

# query by policy_id and name
data "tencentcloud_cam_policies" "bar" {
  policy_id = tencentcloud_cam_policy.foo.id
  name      = "tf-auto-test"
}
```

## Argument Reference

The following arguments are supported:

* `create_mode` - (Optional, Int) Mode of creation of policy strategy. Valid values: `1`, `2`. `1` means policy was created with console, and `2` means it was created by strategies.
* `description` - (Optional, String) The description of the CAM policy.
* `name` - (Optional, String) Name of the CAM policy to be queried.
* `policy_id` - (Optional, String) ID of CAM policy to be queried.
* `result_output_file` - (Optional, String) Used to save results.
* `type` - (Optional, Int) Type of the policy strategy. Valid values: `1`, `2`. `1` means customer strategy and `2` means preset strategy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `policy_list` - A list of CAM policies. Each element contains the following attributes:
  * `attachments` - Number of attached users.
  * `create_mode` - Mode of creation of policy strategy. `1` means policy was created with console, and `2` means it was created by strategies.
  * `create_time` - Create time of the CAM policy.
  * `description` - Description of CAM policy.
  * `name` - Name of CAM policy.
  * `policy_id` - ID of the policy strategy.
  * `service_type` - Name of attached products.
  * `type` - Type of the policy strategy. `1` means customer strategy and `2` means preset strategy.


