Use this data source to query the list of CAM accounts via the CAM `ListAccounts` API, including sub users, collaborators, message receivers and other account types. All pages are fetched automatically, so the `users` list contains every account visible to the credentials.

Example Usage

```hcl
# query all CAM accounts
data "tencentcloud_cam_accounts" "all" {}

# query accounts filtered by user type
data "tencentcloud_cam_accounts" "sub_users" {
  user_type = "SubUser"
}
```
