Provide a resource to create an EMR boot script.

Example Usage

```hcl
resource "tencentcloud_emr_boot_script" "example" {
  instance_id = "emr-qe336v2e"
  boot_type   = "resourceAfter"
  pre_executed_file_settings {
    path          = "demo.py"
    bucket        = "tf-1309115522"
    cos_file_name = "demo"
    region        = "ap-guangzhou"
  }
}
```

Import

EMR boot script can be imported using the id (`{instance_id}#{boot_type}`), e.g.

```
terraform import tencentcloud_emr_boot_script.example emr-qe336v2e#resourceAfter
```
