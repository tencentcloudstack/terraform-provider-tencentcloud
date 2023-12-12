Provides a resource to restart a elasticsearch logstash instance

Example Usage

```hcl
resource "tencentcloud_elasticsearch_restart_logstash_instance_operation" "restart_logstash_instance_operation" {
  instance_id = "ls-xxxxxx"
  type = 0
}
```