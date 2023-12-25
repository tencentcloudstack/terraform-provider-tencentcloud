Provides a resource to create a cvm import_image

Example Usage

```hcl
resource "tencentcloud_cvm_import_image" "import_image" {
  architecture = "x86_64"
  os_type = "CentOS"
  os_version = "7"
  image_url = ""
  image_name = "sample"
  image_description = "sampleimage"
  dry_run = false
  force = false
  tag_specification {
		resource_type = "image"
		tags {
			key = "tagKey"
			value = "tagValue"
		}

  }
  license_type = "TencentCloud"
  boot_mode = "Legacy BIOS"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

cvm import_image can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_import_image.import_image import_image_id
```
