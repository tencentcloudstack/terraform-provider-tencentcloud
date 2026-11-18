Use this data source to query the build log of a Makers deployment.

## Example Usage

```hcl
data "tencentcloud_makers_deployment_log" "example" {
  zone_id       = "zone-123456"
  project_id    = "project-123456"
  deployment_id = "deployment-123456"
}

output "log_url" {
  value = data.tencentcloud_makers_deployment_log.example.log_url
}
```
