Provides a data source to query the TDSQL Boundless (tdmysql) sale information of regions and zones.

Example Usage

Query all sale info

```hcl
data "tencentcloud_tdmysql_sale_info" "example" {}
```

Query sale info by filters

```hcl
data "tencentcloud_tdmysql_sale_info" "example" {
  src_region    = "ap-guangzhou"
  instance_type = "serverless"
  instance_id   = "tdsql3-gpyfw59om6"
}
```
