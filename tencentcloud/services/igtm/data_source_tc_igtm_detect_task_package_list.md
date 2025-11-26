Use this data source to query detailed information of IGTM detect task package list

Example Usage

Query all igtm detect task package list

```hcl
data "tencentcloud_igtm_detect_task_package_list" "example" {}
```

Query igtm detect task package list by filter

```hcl
data "tencentcloud_igtm_detect_task_package_list" "example" {
  filters {
    name  = "ResourceId"
    value = ["task-qqcoptejbwbf"]
    fuzzy = true
  }
}
```
