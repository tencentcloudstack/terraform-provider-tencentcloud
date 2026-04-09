Provides a resource to create a teo teo_function

Example Usage

Basic example (single function query):
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

Example with function_ids parameter (batch query):
```hcl
# Note: The function_ids parameter is used for reading multiple functions at once.
# This is useful for batch operations and data migration scenarios.
# When function_ids is not specified, the default single function query logic is used.
resource "tencentcloud_teo_function" "teo_function_ids" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World');
          e.respondWith(response);
        });
    EOT
    name        = "tf-test-function-ids"
    remark      = "test with function_ids"
    zone_id     = "zone-2qtuhspy7cr6"
    # function_ids = ["function-id-1", "function-id-2"]  # Optional parameter for batch query
}
```

Import

teo teo_function can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function.teo_function zone_id#function_id
```

Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) ID of the site.
* `name` - (Required) Function name. It can only contain lowercase letters, numbers, hyphens, must start and end with a letter or number, and can have a maximum length of 30 characters.
* `remark` - (Optional) Function description, maximum support of 60 characters.
* `content` - (Required) Function content, currently only supports JavaScript code, with a maximum size of 5MB.
* `function_ids` - (Optional) List of function IDs to query. When specified, the read operation will query multiple functions at once. If not specified, the default single function query logic will be used. This parameter is useful for batch operations and data migration scenarios.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `function_id` - ID of the Function.
* `domain` - The default domain name for the function.
* `create_time` - Creation time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.
* `update_time` - Modification time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.

