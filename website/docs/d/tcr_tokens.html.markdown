---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_tokens"
sidebar_current: "docs-tencentcloud-datasource-tcr_tokens"
description: |-
  Use this data source to query detailed information of TCR tokens.
---

# tencentcloud_tcr_tokens

Use this data source to query detailed information of TCR tokens.

## Example Usage

```hcl
data "tencentcloud_tcr_tokens" "name" {
  instance_id = "cls-satg5125"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) ID of the instance that the token belongs to.
* `result_output_file` - (Optional, String) Used to save results.
* `token_id` - (Optional, String) ID of the TCR token to query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `token_list` - Information list of the dedicated TCR tokens.
  * `create_time` - Create time.
  * `description` - Description of the token.
  * `enable` - Indicate that the token is enabled or not.
  * `token_id` - Id of TCR token.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeInstanceToken](https://cloud.tencent.com/document/api/1141/42911) - This API is used to query the information of long-term access credentials.

