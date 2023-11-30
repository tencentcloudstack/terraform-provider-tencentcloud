Use this data source to query detailed information of tcr images

Example Usage

```hcl
data "tencentcloud_tcr_images" "images" {
  registry_id = "tcr-xxx"
  namespace_name = "ns"
  repository_name = "repo"
  image_version = "v1"
  digest = "sha256:xxxxx"
  exact_match = false
  }
```