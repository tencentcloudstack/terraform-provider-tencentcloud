Use this data source to query TEO inference hardware specifications by zone ID, returning the list of available hardware specifications and their resource configuration (CPU, memory, GPU, etc.) for creating inference services.

Example Usage

```hcl
data "tencentcloud_teo_inference_hardware_specifications" "specifications" {
  zone_id = "zone-2qtuhspy7cr6"
}
```