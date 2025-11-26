Use this data source to query detailed information of IGTM instance list

Example Usage

Query all igtm instance list

```hcl
data "tencentcloud_igtm_instance_list" "example" {}
```

Query igtm instance list by filters

```hcl
data "tencentcloud_igtm_instance_list" "example" {
  filters {
    name  = "InstanceId"
    value = ["gtm-uukztqtoaru"]
    fuzzy = true
  }
}
```