Use this data source to query the available database specifications for different regions. And a maximum of 20 requests can be initiated per second for this query.

Example Usage

```hcl
data "tencentcloud_mysql_zone_config" "mysql" {
  region             = "ap-guangzhou"
  result_output_file = "mytestpath"
}
```