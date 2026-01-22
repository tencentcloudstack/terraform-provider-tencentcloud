Use this data source to query detailed information of Private Dns records

Example Usage

```hcl
data "tencentcloud_private_dns_records" "example" {
  zone_id = "zone-kumt5wos"
}
```

Or

```hcl
data "tencentcloud_private_dns_records" "example" {
  zone_id = "zone-kumt5wos"
  filters {
    name   = "RecordType"
    values = ["A"]
  }
}
```
