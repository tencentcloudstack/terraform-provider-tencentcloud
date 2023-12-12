Provides a resource to create a rum taw_instance

Example Usage

```hcl
resource "tencentcloud_rum_taw_instance" "taw_instance" {
  area_id = "1"
  charge_type = "1"
  data_retention_days = "30"
  instance_name = "instanceName-1"
  tags = {
    createdBy = "terraform"
  }
  instance_desc = "instanceDesc-1"
}

```
Import

rum taw_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_rum_taw_instance.taw_instance tawInstance_id
```