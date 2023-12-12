Use this data source to query detailed information of cynosdb resource_package_sale_specs

Example Usage

```hcl
data "tencentcloud_cynosdb_resource_package_sale_specs" "resource_package_sale_specs" {
  instance_type  = "cynosdb-serverless"
  package_region = "china"
  package_type   = "CCU"
}
```