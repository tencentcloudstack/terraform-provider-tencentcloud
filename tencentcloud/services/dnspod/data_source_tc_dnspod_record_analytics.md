Use this data source to query detailed information of dnspod record_analytics

Example Usage

```hcl
data "tencentcloud_dnspod_record_analytics" "record_analytics" {
  domain = "iac-tf.cloud"
  start_date = "2023-09-07"
  end_date = "2023-11-07"
  subdomain = "www"
  dns_format = "HOUR"
  # domain_id = 123
}
```