---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_secrets"
sidebar_current: "docs-tencentcloud-datasource-ssm_secrets"
description: |-
  Use this data source to query detailed information of SSM secret
---

# tencentcloud_ssm_secrets

Use this data source to query detailed information of SSM secret

## Example Usage

```hcl
data "tencentcloud_ssm_secrets" "foo" {
  secret_name = "test"
  order_type  = 1
  state       = 1
}
```

## Argument Reference

The following arguments are supported:

* `order_type` - (Optional, Int) The order to sort the create time of secret. `0` - desc, `1` - asc. Default value is `0`.
* `result_output_file` - (Optional, String) Used to save results.
* `secret_name` - (Optional, String) Secret name used to filter result.
* `state` - (Optional, Int) Filter by state of secret. `0` - all secrets are queried, `1` - only Enabled secrets are queried, `2` - only Disabled secrets are queried, `3` - only PendingDelete secrets are queried.
* `tags` - (Optional, Map) Tags to filter secret.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `secret_list` - A list of SSM secrets.
  * `create_time` - Create time of secret.
  * `create_uin` - Uin of Creator.
  * `delete_time` - Delete time of CMK.
  * `description` - Description of secret.
  * `kms_key_id` - KMS keyId used to encrypt secret.
  * `secret_name` - Name of secret.
  * `status` - Status of secret.


