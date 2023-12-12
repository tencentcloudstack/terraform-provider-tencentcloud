Provides a resource to create a elasticsearch index

Example Usage

```hcl
resource "tencentcloud_elasticsearch_index" "index" {
  instance_id = "es-xxxxxx"
  index_type = "normal"
  index_name = "test-es-index"
  index_meta_json = "{\"mappings\":{},\"settings\":{\"index.number_of_replicas\":1,\"index.number_of_shards\":1,\"index.refresh_interval\":\"30s\"}}"
}
```

Import

elasticsearch index can be imported using the id, e.g.

```
terraform import tencentcloud_elasticsearch_index.index index_id
```