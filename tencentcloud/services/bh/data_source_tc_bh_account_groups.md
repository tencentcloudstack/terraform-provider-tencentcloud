Use this data source to query detailed information of BH account groups

Example Usage

Query all bh account groups

```hcl
data "tencentcloud_bh_account_groups" "example" {}
```

Query bh account groups by filter

```hcl
data "tencentcloud_bh_account_groups" "example" {
  deep_in    = 1
  parent_id  = 819729
  group_name = "tf-example"
}
```
