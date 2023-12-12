Use this data source to query images.

Example Usage

```hcl
data "tencentcloud_images" "foo" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos 7.5"
}
```