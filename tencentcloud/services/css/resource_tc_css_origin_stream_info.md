Provides a resource to create a CSS origin stream info.

Example Usage

```hcl
resource "tencentcloud_css_origin_stream_info" "example" {
  domain_name             = "www.demo.com"
  origin_stream_play_type = "rtmp"
  cdn_stream_play_type    = ["rtmp"]
  origin_stream_type      = 1
  origin_address_type     = 1
  origin_address          = ["1.1.1.1:8080"]
  origin_timeout          = 10000
  origin_retry_times      = 10
}
```

Import

CSS origin stream info can be imported using the domain name, e.g.

```
terraform import tencentcloud_css_origin_stream_info.example www.demo.com
```
