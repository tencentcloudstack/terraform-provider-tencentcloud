Use this data source to query detailed information of antiddos basic_device_status

Example Usage

```hcl
data "tencentcloud_antiddos_basic_device_status" "basic_device_status" {
  ip_list = [
    "127.0.0.1"
  ]
  filter_region = 1
}
```