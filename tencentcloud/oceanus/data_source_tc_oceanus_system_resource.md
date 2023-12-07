Use this data source to query detailed information of oceanus system_resource

Example Usage

```hcl
data "tencentcloud_oceanus_system_resource" "example" {
  resource_ids = ["resource-abd503yt"]
  filters {
    name   = "Name"
    values = ["tf_example"]
  }
  cluster_id    = "cluster-n8yaia0p"
  flink_version = "Flink-1.11"
}
```