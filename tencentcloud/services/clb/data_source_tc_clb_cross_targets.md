Use this data source to query detailed information of clb cross_targets

Example Usage

```hcl
data "tencentcloud_clb_cross_targets" "cross_targets" {
  filters {
    name = "vpc-id"
    values = ["vpc-4owdpnwr"]
  }
}
```