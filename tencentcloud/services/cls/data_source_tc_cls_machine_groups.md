Use this data source to query detailed information of cls machine_groups

Example Usage

```hcl
data "tencentcloud_cls_machine_groups" "name" {
  filters {
    name   = "machineGroupName"
    values = ["tf-machine-group"]
  }
}

data "tencentcloud_cls_machine_groups" "all" {}
```
