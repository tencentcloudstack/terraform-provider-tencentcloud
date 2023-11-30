Use this data source to query detailed information of dcdb accounts.

Example Usage

```hcl
data "tencentcloud_dcdb_accounts" "foo" {
  instance_id = tencentcloud_dcdb_account.foo.instance_id
}
```