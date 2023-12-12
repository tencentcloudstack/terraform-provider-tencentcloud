Use this data source to query the detail information of CDN domain.

Example Usage

```hcl
data "tencentcloud_cdn_domains" "foo" {
  domain         	   = "xxxx.com"
  service_type   	   = "web"
  full_url_cache 	   = false
  origin_pull_protocol = "follow"
  https_switch		   = "on"
}
```