---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_accounts"
sidebar_current: "docs-tencentcloud-datasource-cam_accounts"
description: |-
  Use this data source to query the list of CAM accounts via the CAM `ListAccounts` API, including sub users, collaborators, message receivers and other account types.
---

# tencentcloud_cam_accounts

Use this data source to query the list of CAM accounts via the CAM `ListAccounts` API, including sub users, collaborators, message receivers and other account types.

## Example Usage

```hcl
# query all CAM accounts
data "tencentcloud_cam_accounts" "all" {
}

# query accounts filtered by user type
data "tencentcloud_cam_accounts" "sub_users" {
  user_type = "SubUser"
}

# query accounts with paging
data "tencentcloud_cam_accounts" "page" {
  max_items = 50
  marker    = "previous-marker"
}
```

## Argument Reference

The following arguments are supported:

* `marker` - (Optional, String) When the returned result is truncated, use Marker to fetch the content after the current truncation position. Output `marker` carries the next page marker when `is_truncated` is true.
* `max_items` - (Optional, Int) Maximum number of accounts to return per request. Valid range: [1, 100]. When the returned result is truncated due to reaching MaxItems, the output `is_truncated` will be true.
* `result_output_file` - (Optional, String) Used to save results.
* `user_type` - (Optional, String) Account type to filter by. Valid values: `Owner` (master account), `SubUser` (sub user), `CICUser` (CIC sub user), `WechatCorpUser` (WeCom sub user), `AgentIdentity` (AgentIdentity sub user), `Collaborator` (collaborator), `MessageReceiver` (message receiver).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `is_truncated` - Whether the returned result is truncated.
* `users` - A list of CAM accounts. Each element contains the following attributes:
  * `console_login` - Whether the account can log in to the console. Returned as the raw int value from the API.
  * `country_code` - Country code.
  * `create_time` - Creation time. Format: YYYY-MM-DD hh:mm:ss.
  * `email` - Email.
  * `name` - Account name.
  * `phone_num` - Phone number.
  * `remark` - Account remark.
  * `uid` - Account UID.
  * `uin` - Account ID (Uin).
  * `user_type` - Account type. Valid values: `Owner`, `SubUser`, `CICUser`, `WechatCorpUser`, `AgentIdentity`, `Collaborator`, `MessageReceiver`, `Unknown`.


