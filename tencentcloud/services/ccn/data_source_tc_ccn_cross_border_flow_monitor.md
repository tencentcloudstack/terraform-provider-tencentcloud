Use this data source to query detailed information of vpc cross_border_flow_monitor

Example Usage

```hcl
data "tencentcloud_ccn_cross_border_flow_monitor" "cross_border_flow_monitor" {
  source_region = "ap-guangzhou"
  destination_region = "ap-singapore"
  ccn_id = "ccn-39lqkygf"
  ccn_uin = "979137"
  period = 60
  start_time = "2023-01-01 00:00:00"
  end_time = "2023-01-01 01:00:00"
}
```