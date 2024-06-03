Provides a resource to create a cvm sync_image

Example Usage

```hcl
data "tencentcloud_images" "example" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

resource "tencentcloud_cvm_sync_image" "example" {
  image_id            = data.tencentcloud_images.example.images.0.image_id
  destination_regions = ["ap-guangzhou", "ap-shanghai"]
}
```