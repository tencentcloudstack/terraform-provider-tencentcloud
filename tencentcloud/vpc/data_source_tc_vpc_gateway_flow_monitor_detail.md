Use this data source to query detailed information of vpc gateway_flow_monitor_detail

Example Usage

```hcl
data "tencentcloud_vpc_gateway_flow_monitor_detail" "gateway_flow_monitor_detail" {
  time_point      = "2023-06-02 12:15:20"
  vpn_id          = "vpngw-gt8bianl"
  order_field     = "OutTraffic"
  order_direction = "DESC"
}
```