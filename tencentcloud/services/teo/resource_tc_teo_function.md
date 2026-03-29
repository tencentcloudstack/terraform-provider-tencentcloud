Provides a resource to create a teo teo_function

Example Usage

**Example 1: Create function with auto-generated function_id**

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

**Example 2: Create function with specified function_id**

```hcl
resource "tencentcloud_teo_function" "teo_function_with_id" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World!!');
          e.respondWith(response);
        });
    EOT
    name        = "aaa-zone-2qtuhspy7cr6-1310708577"
    remark      = "test"
    zone_id     = "zone-2qtuhspy7cr6"
    function_id = "custom-function-id-123"
}
```

> **Note:** The `function_id` parameter is optional. If not specified, the API will generate a function ID automatically.

Import

teo teo_function can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function.teo_function zone_id#function_id
```
