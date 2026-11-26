Provides a data source to query the sellable component specs of TDSQL Boundless (tdmysql).

Example Usage

```hcl
data "tencentcloud_tdmysql_specs" "example" {
  full_replications      = 3
  is_exclusive_instance  = 0
  zones                  = ["ap-guangzhou-3"]
  instance_id            = "tdsqlshard-12345"
}
```
