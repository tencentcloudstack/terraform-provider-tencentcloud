Provides a resource to create a elasticsearch diagnose

Example Usage

```hcl
resource "tencentcloud_elasticsearch_diagnose" "diagnose" {
  instance_id = "es-xxxxxx"
  cron_time = "15:00:00"
}
```

Import

es diagnose can be imported using the id, e.g.

```
terraform import tencentcloud_elasticsearch_diagnose.diagnose diagnose_id
```