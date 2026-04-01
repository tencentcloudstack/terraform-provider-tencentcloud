Provides a resource to create a teo teo_function

Example Usage

```hcl
resource "tencentcloud_teo_function" "teo_function" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World!!');
          e.respondWith(response);
        });
    EOT
    name        = "aaa-zone-2qtuhspy7cr6-1310708577"
    remark      = "test"
    zone_id     = "zone-2qtuhspy7cr6"
}
```

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `function_id` - ID of the Function.
* `domain` - The default domain name for the function.
* `create_time` - Creation time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.
* `update_time` - Modification time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.
* `functions` - List of all function information in the zone. Each object contains:
    * `function_id` - Function ID.
    * `zone_id` - Site ID.
    * `name` - Function name.
    * `remark` - Function description.
    * `content` - Function content.
    * `domain` - Function default domain.
    * `create_time` - Creation time (UTC, ISO 8601 format).
    * `update_time` - Modification time (UTC, ISO 8601 format).

## Import

teo teo_function can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function.teo_function zone_id#function_id
```
