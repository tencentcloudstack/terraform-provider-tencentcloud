Use this data source to query detailed information of teo inference service deployment logs

Example Usage

```hcl
data "tencentcloud_teo_inference_service_deployment_logs" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "sid-2quhspyeq8r6"
  record_id  = "rec-2quhspyeq8r6"
}
```