Provides a resource to check CNAME status for TEO domains.

Example Usage

```hcl
resource "tencentcloud_teo_check_cname_status_operation" "example" {
  zone_id = "zone-12345678"
  record_names = [
    "example.com",
    "www.example.com",
  ]
}
```

Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Zone ID.
* `record_names` - (Required, ForceNew) List of record names to check CNAME status.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cname_status` - CNAME status information for each record name.
  * `record_name` - Record name.
  * `cname` - CNAME address. May be null.
  * `status` - CNAME status. Valid values: `active`, `moved`.

Import

The resource can be imported by using the `zone_id`, e.g.

```sh
terraform import tencentcloud_teo_check_cname_status.example zone-12345678
```
