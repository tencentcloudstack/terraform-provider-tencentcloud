Provides a resource to create a cvm sync_image

Example Usage

```hcl
data "tencentcloud_images" "example" {
  image_type       = ["PRIVATE_IMAGE"]
  image_name_regex = "MyImage"
}

resource "tencentcloud_cvm_sync_image" "example" {
  image_id            = data.tencentcloud_images.example.images.0.image_id
  destination_regions = ["ap-guangzhou", "ap-shanghai"]
  encrypt             = true
  kms_key_id          = "f063c18b-654b-11ef-9d9f-525400d3a886"
}
```