Use this data source to query placement groups.

Example Usage

```hcl
data "tencentcloud_placement_groups" "example" {
  placement_group_id = "ps-bwvst92h"
  name               = "tf_example"
}
```