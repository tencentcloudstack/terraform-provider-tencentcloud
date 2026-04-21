Provides a resource to identify a zone or subdomain ownership.

Example Usage

```hcl
resource "tencentcloud_teo_identify_zone_operation" "example" {
  zone_name = "example.com"
  domain    = "www.example.com"
}
```

Argument Reference

The following arguments are supported:

* `zone_name` - (Required, ForceNew) Zone name.
* `domain` - (Optional, ForceNew) Subdomain under the zone. Required only when verifying a subdomain.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ascription` - DNS verification information.
  * `subdomain` - DNS record host.
  * `record_type` - DNS record type.
  * `record_value` - DNS record value.
* `file_ascription` - File verification information.
  * `identify_path` - File verification path.
  * `identify_content` - File verification content.
