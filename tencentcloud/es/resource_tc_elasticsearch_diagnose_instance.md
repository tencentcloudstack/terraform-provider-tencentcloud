Provides a resource to create a elasticsearch diagnose instance

Example Usage

```hcl
resource "tencentcloud_elasticsearch_diagnose_instance" "diagnose_instance" {
  instance_id = "es-xxxxxx"
  diagnose_jobs = ["cluster_health"]
  diagnose_indices = "*"
}
```