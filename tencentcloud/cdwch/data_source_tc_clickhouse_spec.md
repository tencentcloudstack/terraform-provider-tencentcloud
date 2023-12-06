Use this data source to query detailed information of clickhouse spec

Example Usage

```hcl
data "tencentcloud_clickhouse_spec" "spec" {
  zone       = "ap-guangzhou-7"
  pay_mode   = "PREPAID"
  is_elastic = false
}
```