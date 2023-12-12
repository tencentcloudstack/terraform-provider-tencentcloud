Use this data source to query detailed information of vpc net_detect_state_check

Example Usage

```hcl
data "tencentcloud_vpc_net_detect_state_check" "net_detect_state_check" {
  net_detect_id         = "netd-12345678"
  detect_destination_ip = [
    "10.0.0.3",
    "10.0.0.2"
  ]
  next_hop_type        = "NORMAL_CVM"
  next_hop_destination = "10.0.0.4"
}
```