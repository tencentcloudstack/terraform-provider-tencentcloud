Provides a resource to create a cvm export_images

Example Usage

```hcl
resource "tencentcloud_cvm_export_images" "export_images" {
  bucket_name = "xxxxxx"
  image_id = "img-xxxxxx"
  file_name_prefix = "test-"
}
```