Use this data source to query detailed information of IGTM instance package list

Example Usage

Query all igtm instance package list

```hcl
data "tencentcloud_igtm_instance_package_list" "example" {}
```

Query igtm instance package list by filter

```hcl
data "tencentcloud_igtm_instance_package_list" "example" {
  filters {
    name  = "InstanceId"
    value = ["gtm-uukztqtoaru"]
    fuzzy = true
  }
}
```
