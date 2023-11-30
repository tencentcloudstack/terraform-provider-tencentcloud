Use this data source to query detailed information of tcr image_manifests

Example Usage

```hcl
data "tencentcloud_tcr_image_manifests" "image_manifests" {
	registry_id = "%s"
	namespace_name = "%s"
	repository_name = "%s"
	image_version = "v1"
}
```