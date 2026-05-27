Use this data source to query GA2 (Global Accelerator 2.0) cross-border settlement traffic usage information.

Example Usage

Query cross-border settlement traffic

```hcl
data "tencentcloud_ga2_cross_border_settlement" "example" {
  global_accelerator_id = "ga2-xxxxxxxx"
  accelerate_region     = "ap-guangzhou"
  endpoint_group_region = "ap-singapore"
  settlement_month      = 202501
}
```
