Provides an available image for the user.

The Images data source fetch proper image, which could be one of the private images of the user and images of system resources provided by TencentCloud, as well as other public images and those available on the image market.

~> **NOTE:** This data source will be deprecated, please use `tencentcloud_images` instead.

Example Usage

```hcl
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"

  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}
```