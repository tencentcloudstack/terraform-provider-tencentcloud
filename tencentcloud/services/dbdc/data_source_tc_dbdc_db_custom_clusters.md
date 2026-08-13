Use this data source to query DB Custom cluster list from TencentCloud DBDC product.

Example Usage

Query all dbdc db custom clusters

```hcl
data "tencentcloud_dbdc_db_custom_clusters" "example" {}
```

Query dbdc db custom clusters by cluster_ids

```hcl
data "tencentcloud_dbdc_db_custom_clusters" "example" {
  cluster_ids = [
    "dbcc-nmtmsew8",
    "dbcc-9yui67ac"
  ]
}
```

Query dbdc db custom clusters by filters

```hcl
data "tencentcloud_dbdc_db_custom_clusters" "example" {
  filters {
    name   = "cluster-name"
    values = ["tf-example"]
  }

  filters {
    name   = "cluster-status"
    values = ["Running"]
  }
}
```

Query dbdc db custom clusters by tags

```hcl
data "tencentcloud_dbdc_db_custom_clusters" "example" {
  tags {
    key   = "env"
    value = "production"
  }
}
```
