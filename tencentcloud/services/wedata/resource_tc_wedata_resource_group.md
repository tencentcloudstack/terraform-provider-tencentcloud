Provides a resource to create a WeData resource group

~> **NOTE:** If an SKDe Error message appears when executing the `terraform destroy` command, please contact Tencent Cloud WeData for consultation.

Example Usage

```hcl
resource "tencentcloud_wedata_resource_group" "example" {
  name = "tf_example"
  type {
    resource_group_type = "Integration"
    integration {
      real_time_data_sync {
        specification = "i32c"
        number        = 1
      }

      offline_data_sync {
        specification = "integrated"
        number        = 2
      }
    }
  }

  auto_renew_enabled = false
  purchase_period    = 1
  vpc_id             = "vpc-ds5rpnxh"
  subnet             = "subnet-fz7rw5zq"
  resource_region    = "ap-beijing-fsi"
  description        = "description."

  lifecycle {
    ignore_changes = [ description, resource_region ]
  }
}
```
