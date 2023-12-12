Provides a resource to create a cvm launch template

Example Usage

```hcl
data "tencentcloud_images" "my_favorite_image" {
  image_type = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

resource "tencentcloud_cvm_launch_template" "demo" {
  launch_template_name = "test"
  placement {
    zone = "ap-guangzhou-6"
    project_id = 0
  }
  image_id = data.tencentcloud_images.my_favorite_image.images.0.image_id
}
```