Use this data source to query detailed information of dnspod domain_log_list

Example Usage

```hcl
data "tencentcloud_dnspod_domain_log_list" "domain_log_list" {
  domain = "iac-tf.cloud"
  domain_id = 123
}
```