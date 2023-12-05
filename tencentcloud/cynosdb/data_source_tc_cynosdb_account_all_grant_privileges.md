Use this data source to query detailed information of cynosdb account_all_grant_privileges

Example Usage

```hcl
data "tencentcloud_cynosdb_account_all_grant_privileges" "account_all_grant_privileges" {
  cluster_id = "cynosdbmysql-bws8h88b"
  account {
    account_name = "keep_dts"
    host         = "%"
  }
}
```