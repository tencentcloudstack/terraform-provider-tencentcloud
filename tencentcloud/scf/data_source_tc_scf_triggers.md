Use this data source to query detailed information of scf triggers

Example Usage

```hcl
data "tencentcloud_scf_triggers" "triggers" {
  function_name = "keep-1676351130"
  namespace     = "default"
  order_by      = "add_time"
  order         = "DESC"
}
```