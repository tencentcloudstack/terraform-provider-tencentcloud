Use this data source to query DB Custom node list from TencentCloud DBDC product.

Example Usage

Query all dbdc db custom nodes

```hcl
data "tencentcloud_dbdc_db_custom_nodes" "example" {}
```

Query dbdc db custom nodes by node_ids

```hcl
data "tencentcloud_dbdc_db_custom_nodes" "example" {
  node_ids = [
    "dbcn-abc12345",
    "dbcn-def67890"
  ]
}
```

Query dbdc db custom nodes by filters

```hcl
data "tencentcloud_dbdc_db_custom_nodes" "example" {
  filters {
    name   = "cluster-id"
    values = ["dbcc-nmtmsew8"]
  }

  filters {
    name   = "status"
    values = ["Running"]
  }
}
```

Query dbdc db custom nodes by tags

```hcl
data "tencentcloud_dbdc_db_custom_nodes" "example" {
  tags {
    key   = "env"
    value = "production"
  }
}
```
