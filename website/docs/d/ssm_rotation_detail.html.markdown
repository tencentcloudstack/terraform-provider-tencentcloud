---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_rotation_detail"
sidebar_current: "docs-tencentcloud-datasource-ssm_rotation_detail"
description: |-
  Use this data source to query detailed information of ssm rotation_detail
---

# tencentcloud_ssm_rotation_detail

Use this data source to query detailed information of ssm rotation_detail

## Example Usage

```hcl
data "tencentcloud_ssm_rotation_detail" "example" {
  secret_name = "tf_example"
}
```

## Argument Reference

The following arguments are supported:

* `secret_name` - (Required, String) Secret name.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `enable_rotation` - Whether to allow rotation.
* `frequency` - The rotation frequency, in days, defaults to 1 day.
* `latest_rotate_time` - Time of last rotation.
* `next_rotate_begin_time` - The time to start the next rotation.


