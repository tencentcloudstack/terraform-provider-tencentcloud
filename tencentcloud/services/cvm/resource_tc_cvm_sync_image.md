Provides a resource to create a cvm sync_image

Example Usage

```hcl
resource "tencentcloud_cvm_sync_image" "sync_image" {
  image_id = "img-xxxxxx"
  destination_regions =["ap-guangzhou", "ap-shanghai"]
}
```