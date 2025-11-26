Use this data source to query detailed information of IGTM monitors

Example Usage

Query all igtm monitors

```hcl
data "tencentcloud_igtm_monitors" "example" {}
```

Query igtm monitors by filter

```hcl
data "tencentcloud_igtm_monitors" "example" {
  filters {
    name  = "MonitorId"
    value = ["12383"]
    fuzzy = true
  }
}
```
