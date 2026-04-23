Use this data source to query detailed information of TEO IP region

Example Usage

Query IP region info

```hcl
data "tencentcloud_teo_ip_region" "example" {
  ips = [
    "1.1.1.1",
    "2.2.2.2"
  ]
}
```
