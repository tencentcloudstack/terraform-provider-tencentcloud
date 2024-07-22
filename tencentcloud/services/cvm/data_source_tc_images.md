Use this data source to query images.

Example Usage

Query all images

```hcl
data "tencentcloud_images" "example" {}
```

Query images by image ID

```hcl
data "tencentcloud_images" "example" {
  image_id = "img-9qrfy1xt"
}
```

Query images by os name

```hcl
data "tencentcloud_images" "example" {
  os_name = "TencentOS Server 3.2 (Final)"
}
```

Query images by image name regex

```hcl
data "tencentcloud_images" "example" {
  image_name_regex = "^TencentOS"
}
```

Query images by image type

```hcl
data "tencentcloud_images" "example" {
  image_type = ["PUBLIC_IMAGE"]
}
```

Query images by instance type

```hcl
data "tencentcloud_images" "example" {
  instance_type = "S1.SMALL1"
}
```
