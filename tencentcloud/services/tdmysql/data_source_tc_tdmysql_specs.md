Provides a data source to query the sellable component specs of TDSQL Boundless (tdmysql).

Example Usage

Query all specs

```hcl
data "tencentcloud_tdmysql_specs" "example" {}
```

Query specs by filters

```hcl
data "tencentcloud_tdmysql_specs" "example" {
  full_replications     = 3
  is_exclusive_instance = 0
  zones                 = ["ap-guangzhou-6", "ap-guangzhou-7"]
  instance_id           = "tdsql3-gpyfw59om6"
}
```
