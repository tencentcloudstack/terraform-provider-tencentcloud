Use this data source to query detailed information of private dns records

Example Usage

```hcl
data "tencentcloud_private_dns_records" "private_dns_record" {
  zone_id = "zone-xxxxxx"
  filters {
	name = "Value"
	values = ["8.8.8.8"]
  }
}
```