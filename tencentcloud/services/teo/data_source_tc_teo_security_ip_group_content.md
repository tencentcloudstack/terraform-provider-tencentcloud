Use this data source to query the IP list within a specified TEO security IP group.

Example Usage

```hcl
data "tencentcloud_teo_security_ip_group_content" "example" {
  zone_id  = "zone-2qtuhspy7cr6"
  group_id = 123
}
```
