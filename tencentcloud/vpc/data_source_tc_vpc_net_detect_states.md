Use this data source to query detailed information of vpc net_detect_states

Example Usage

```hcl
data "tencentcloud_vpc_net_detect_states" "net_detect_states" {
  net_detect_ids = ["netd-12345678"]
}
```