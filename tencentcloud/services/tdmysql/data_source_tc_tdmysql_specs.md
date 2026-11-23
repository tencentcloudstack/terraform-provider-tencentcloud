Use this data source to query TDSQL Boundless (tdmysql) salable component specs.

Example Usage

```hcl
data "tencentcloud_tdmysql_specs" "example" {
  full_replications     = 2
  is_exclusive_instance = 1
  zones                 = ["ap-guangzhou-1"]
  instance_id           = "tdsql-o2mtest"
}

output "first_hybrid_spec_code" {
  value = data.tencentcloud_tdmysql_specs.example.hybrid_node_specs[0].spec_code
}
```
