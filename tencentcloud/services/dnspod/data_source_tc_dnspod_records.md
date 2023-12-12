Use this data source to query dnspod record list.

Example Usage

```hcl
data "tencentcloud_dnspod_records" "record" {
  domain = "example.com"
  subdomain = "www"
}

output "result" {
  value = data.tencentcloud_dnspod_records.record.result
}
```

Use verbose filter

```hcl
data "tencentcloud_dnspod_records" "record" {
  domain = "example.com"
  subdomain = "www"
  limit = 100
  record_type = "TXT"
  sort_field = "updated_on"
  sort_type = "DESC"
}

output "result" {
  value = data.tencentcloud_dnspod_records.record.result
}
```