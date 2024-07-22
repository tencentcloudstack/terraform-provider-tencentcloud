Provides an available image for the user.

The Images data source fetch proper image, which could be one of the private images of the user and images of system
resources provided by TencentCloud, as well as other public images and those available on the image market.

~> **NOTE:** This data source will be deprecated, please use `tencentcloud_images` instead.

Example Usage

Query image

```hcl
data "tencentcloud_image" "example" {}
```

Query image by filter

```hcl
data "tencentcloud_image" "example" {
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}
```

Query image by os name

```hcl
data "tencentcloud_image" "example" {
  os_name = "centos"
}
```

Query image by image name regex

```hcl
data "tencentcloud_image" "example" {
  image_name_regex = "^Windows\\s.*$"
}
```