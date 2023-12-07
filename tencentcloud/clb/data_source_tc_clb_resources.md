Use this data source to query detailed information of clb resources

Example Usage

```hcl
data "tencentcloud_clb_resources" "resources" {
  filters {
    name = "isp"
    values = ["BGP"]
  }
}
```