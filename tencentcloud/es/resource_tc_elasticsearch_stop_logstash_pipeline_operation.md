Provides a resource to stop a elasticsearch logstash pipeline

Example Usage

```hcl
resource "tencentcloud_elasticsearch_stop_logstash_pipeline_operation" "stop_logstash_pipeline_operation" {
  instance_id = "ls-xxxxxx"
  pipeline_id = "xxxxxx"
}
```