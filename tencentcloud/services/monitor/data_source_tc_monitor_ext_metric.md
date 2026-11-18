Use this data source to query Cloud Monitor external metrics of an external namespace.

Example Usage

```hcl
data "tencentcloud_monitor_ext_metric" "example" {
  ext_namespace = "qcs::monitor:my-namespace"
}
```
