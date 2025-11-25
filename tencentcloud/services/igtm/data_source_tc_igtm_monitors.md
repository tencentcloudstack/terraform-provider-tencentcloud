Use this data source to query detailed information of IGTM monitors

Example Usage

```hcl
data "tencentcloud_igtm_monitors" "example" {
  filters {
    name  = "MonitorId"
    value = ["12383"]
    fuzzy = true
  }
}
```
