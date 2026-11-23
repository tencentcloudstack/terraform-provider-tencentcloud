---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_program_fpga_image"
sidebar_current: "docs-tencentcloud-resource-cvm_program_fpga_image"
description: |-
  Provides a resource to create a cvm program_fpga_image
---

# tencentcloud_cvm_program_fpga_image

Provides a resource to create a cvm program_fpga_image

## Example Usage

```hcl
resource "tencentcloud_cvm_program_fpga_image" "program_fpga_image" {
  instance_id = "ins-xxxxxx"
  fpga_url    = ""
  dbd_fs      = ""
}
```

## Argument Reference

The following arguments are supported:

* `fpga_url` - (Required, String, ForceNew) COS URL address of the FPGA image file.
* `instance_id` - (Required, String, ForceNew) The ID information of the instance.
* `dbd_fs` - (Optional, Set: [`String`], ForceNew) The DBDF number of the FPGA card on the instance, if left blank, the FPGA image will be burned to all FPGA cards owned by the instance by default.
* `dry_run` - (Optional, Bool, ForceNew) Trial run, will not perform the actual burning action, the default is False.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ProgramFpgaImage](https://www.tencentcloud.com/search/?wd=ProgramFpgaImage&type=document&lang=en)

## Import

cvm program_fpga_image can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_program_fpga_image.program_fpga_image program_fpga_image_id
```

