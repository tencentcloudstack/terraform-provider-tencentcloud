Use this data source to query the reference information of a specified TEO (EdgeOne) IP group.

Example Usage

```hcl
data "tencentcloud_teo_ip_group_references" "example" {
  zone_id  = "zone-3fkff38fyw8s"
  group_id = 33711
}
```
