Provides a resource to create a cvm program_fpga_image

Example Usage

```hcl
resource "tencentcloud_cvm_program_fpga_image" "program_fpga_image" {
  instance_id = "ins-xxxxxx"
  fpga_url = ""
  dbd_fs = ""
}
```

Import

cvm program_fpga_image can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_program_fpga_image.program_fpga_image program_fpga_image_id
```
