Use this data source to query DB Custom node machine types from the TencentCloud DBDC product.

Example Usage

Query all dbdc db custom node types

```hcl
data "tencentcloud_dbdc_db_custom_node_types" "example" {}
```

Query dbdc db custom node types by filters

```hcl
data "tencentcloud_dbdc_db_custom_node_types" "example" {
  filters {
    name   = "zone"
    values = ["ap-guangzhou-6"]
  }

  filters {
    name   = "node-family"
    values = ["DB.SA5"]
  }
}

output "node_type_set" {
  value = data.tencentcloud_dbdc_db_custom_node_types.example.node_type_set
}
```
