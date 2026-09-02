Use this data source to query the list of CAM accounts via the CAM `ListAccounts` API, including sub users, collaborators, message receivers and other account types.

Example Usage

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
