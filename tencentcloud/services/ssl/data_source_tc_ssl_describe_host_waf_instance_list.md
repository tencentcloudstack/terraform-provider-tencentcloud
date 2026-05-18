Use this data source to query detailed information of SSL describe host waf instance list

Example Usage

```hcl
data "tencentcloud_ssl_describe_host_waf_instance_list" "example" {
  certificate_id = "GGQ0tJxn"
  resource_type  = "waf"
}
```