---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_accounts"
sidebar_current: "docs-tencentcloud-datasource-cam_accounts"
description: |-
  Use this data source to query the list of CAM accounts via the CAM `ListAccounts` API, including sub users, collaborators, message receivers and other account types. All pages are fetched automatically, so the `users` list contains every account visible to the credentials.
---

# tencentcloud_cam_accounts

Use this data source to query the list of CAM accounts via the CAM `ListAccounts` API, including sub users, collaborators, message receivers and other account types. All pages are fetched automatically, so the `users` list contains every account visible to the credentials.

## Example Usage

```hcl
# query all CAM accounts
data "tencentcloud_cam_accounts" "all" {}

# query accounts filtered by user type
data "tencentcloud_cam_accounts" "sub_users" {
  user_type = "SubUser"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `user_type` - (Optional, String) Account type to filter by. Valid values: `Owner` (master account), `SubUser` (sub user), `CICUser` (CIC sub user), `WechatCorpUser` (WeCom sub user), `AgentIdentity` (AgentIdentity sub user), `Collaborator` (collaborator), `MessageReceiver` (message receiver).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

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


