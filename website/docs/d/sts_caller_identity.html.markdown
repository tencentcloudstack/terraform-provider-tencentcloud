---
subcategory: "Security Token Service(STS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sts_caller_identity"
sidebar_current: "docs-tencentcloud-datasource-sts_caller_identity"
description: |-
  Use this data source to query detailed information of sts callerIdentity
---

# tencentcloud_sts_caller_identity

Use this data source to query detailed information of sts callerIdentity

## Example Usage

```hcl
data "tencentcloud_sts_caller_identity" "callerIdentity" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `account_id` - The primary account Uin to which the current caller belongs.
* `arn` - Current caller ARN.
* `principal_id` - Account Uin to which the key belongs:- The caller is a cloud account, and the returned current account Uin- The caller is a role, and the returned account Uin that applies for the role key.
* `type` - Identity type.
* `user_id` - Identity:- When the caller is a cloud account, the current account `Uin` is returned.- When the caller is a role, it returns `roleId:roleSessionName`- When the caller is a federated identity, it returns `uin:federatedUserName`.


