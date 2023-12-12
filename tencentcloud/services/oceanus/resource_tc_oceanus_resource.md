Provides a resource to create a oceanus resource

Example Usage

```hcl
resource "tencentcloud_oceanus_resource" "example" {
  resource_loc {
    storage_type = 1
    param {
      bucket = "keep-terraform-1257058945"
      path   = "OceanusResource/junit-4.13.2.jar"
      region = "ap-guangzhou"
    }
  }

  resource_type          = 1
  remark                 = "remark."
  name                   = "tf_example"
  resource_config_remark = "config remark."
  folder_id              = "folder-7ctl246z"
  work_space_id          = "space-2idq8wbr"
}
```