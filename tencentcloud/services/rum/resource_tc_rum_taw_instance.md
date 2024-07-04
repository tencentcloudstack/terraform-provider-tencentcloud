Provides a resource to create a rum taw_instance

Example Usage

```hcl
resource "tencentcloud_rum_taw_instance" "example" {
  area_id             = "1"
  charge_type         = "1"
  data_retention_days = "30"
  instance_name       = "tf-example"
  instance_desc       = "desc."

  tags = {
    createdBy = "terraform"
  }
}
```
Import

rum taw_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_rum_taw_instance.example rum-WYsjMaiGOVXBeJ
```
