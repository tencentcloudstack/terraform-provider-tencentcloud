Use this data source to query detailed information of oceanus clusters

Example Usage

Query all clusters

```hcl
data "tencentcloud_oceanus_clusters" "example" {}
```

Query the specified cluster

```hcl
data "tencentcloud_oceanus_clusters" "example" {
  cluster_ids = ["cluster-5c42n3a5"]
  order_type  = 1
  filters {
    name   = "name"
    values = ["tf_example"]
  }
  work_space_id = "space-2idq8wbr"
}
```