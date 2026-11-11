Use this data source to query detailed information of cls machine_groups

Example Usage

Query all machine groups

```hcl
data "tencentcloud_cls_machine_groups" "example" {}
```

Query machine group by filters

```hcl
data "tencentcloud_cls_machine_groups" "example" {
  filters {
    name   = "machineGroupId"
    values = ["76e09d6f-e0c5-4103-bd6e-22bdbf63a76e"]
  }

  filters {
    name   = "machineGroupName"
    values = ["cls-m26ybna6"]
  }
}
```
