Use this data source to query TEO inference service deployment records, such as the operation type, status, duration and configuration snapshot of each deployment.

Example Usage

```hcl
data "tencentcloud_teo_inference_service_deployment_records" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "sid-xxxxxxxxxxxx"
  sort_by    = "create-time"
  sort_order = "desc"
}
```
