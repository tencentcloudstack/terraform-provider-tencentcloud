Provides a resource to create a WeData resource group

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
}
```
