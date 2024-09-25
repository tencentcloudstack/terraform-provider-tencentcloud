Use this data source to query detailed information of Postgresql dedicated clusters

Example Usage

Query all instances

```hcl
data "tencentcloud_postgresql_dedicated_clusters" "example" {}
```

Query instances by filters

```hcl
data "tencentcloud_postgresql_dedicated_clusters" "example" {
  filters {
    name = "dedicated-cluster-id"
    values = ["cluster-262n63e8"]
  }
}
```
