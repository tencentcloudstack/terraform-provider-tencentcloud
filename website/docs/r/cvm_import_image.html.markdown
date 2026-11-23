---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_import_image"
sidebar_current: "docs-tencentcloud-resource-cvm_import_image"
description: |-
  Provides a resource to create a cvm import_image
---

# tencentcloud_cvm_import_image

Provides a resource to create a cvm import_image

## Example Usage

```hcl
resource "tencentcloud_cvm_import_image" "import_image" {
  architecture      = "x86_64"
  os_type           = "CentOS"
  os_version        = "7"
  image_url         = ""
  image_name        = "sample"
  image_description = "sampleimage"
  dry_run           = false
  force             = false
  tag_specification {
    resource_type = "image"
    tags {
      key   = "tagKey"
      value = "tagValue"
    }

  }
  license_type = "TencentCloud"
  boot_mode    = "Legacy BIOS"
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `architecture` - (Required, String, ForceNew) OS architecture of the image to be imported, `x86_64` or `i386`.
* `image_name` - (Required, String, ForceNew) Image name.
* `image_url` - (Required, String, ForceNew) Address on COS where the image to be imported is stored.
* `os_type` - (Required, String, ForceNew) OS type of the image to be imported. You can call `DescribeImportImageOs` to obtain the list of supported operating systems.
* `os_version` - (Required, String, ForceNew) OS version of the image to be imported. You can call `DescribeImportImageOs` to obtain the list of supported operating systems.
* `boot_mode` - (Optional, String, ForceNew) Boot mode.
* `dry_run` - (Optional, Bool, ForceNew) Dry run to check the parameters without performing the operation.
* `force` - (Optional, Bool, ForceNew) Whether to force import the image.
* `image_description` - (Optional, String, ForceNew) Image description.
* `license_type` - (Optional, String, ForceNew) The license type used to activate the OS after importing an image. Valid values: TencentCloud: Tencent Cloud official license BYOL: Bring Your Own License.
* `tag_specification` - (Optional, List, ForceNew) Tag description list. This parameter is used to bind a tag to a custom image.

The `tag_specification` object supports the following:

* `resource_type` - (Required, String) Resource type. Valid values: instance (CVM), host (CDH), image (for image), and keypair (for key). Note: This field may return null, indicating that no valid values can be obtained.
* `tags` - (Required, List) Tag pairs Note: This field may return null, indicating that no valid values can be obtained.

The `tags` object of `tag_specification` supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ImportImage](https://www.tencentcloud.com/document/api/213/33270)

## Import

cvm import_image can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_import_image.import_image import_image_id
```

