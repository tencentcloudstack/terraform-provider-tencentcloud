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

### Query all policies

```hcl
data "tencentcloud_cam_policies" "example" {}
```

### Query policies by filter

```hcl
data "tencentcloud_cam_policies" "example" {
  name        = "tf-example"
  policy_id   = "236215899"
  type        = 1
  create_mode = 2
}
```

## Argument Reference

The following arguments are supported:

* `create_mode` - (Optional, Int) Mode of creation of policy strategy. Valid values: `1`, `2`. `1` means policy was created with console, and `2` means it was created by strategies.
* `description` - (Optional, String) The description of the CAM policy.
* `key_word` - (Optional, String) Match by strategy name.
* `name` - (Optional, String) Name of the CAM policy to be queried.
* `policy_id` - (Optional, String) ID of CAM policy to be queried.
* `result_output_file` - (Optional, String) Used to save results.
* `scope` - (Optional, String) Available values are 'All', 'QCS', and' Local '.' All 'retrieves all policies,' QCS' retrieves preset policies, 'Local' retrieves custom policies, and defaults to 'All'.
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


