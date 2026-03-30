Provides a resource to create a teo teo_function

Example Usage

### Create function with API-generated FunctionId (default behavior)

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

### Create function with custom FunctionId

```hcl
resource "tencentcloud_teo_function" "custom_function" {
    function_id = "my-custom-function-id"
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World!!');
          e.respondWith(response);
        });
    EOT
    name        = "custom-function"
    remark      = "test"
    zone_id     = "zone-2qtuhspy7cr6"
}
```

### Import existing function with known FunctionId

```hcl
terraform import tencentcloud_teo_function.custom_function zone_id#function_id
```

Import

teo teo_function can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function.teo_function zone_id#function_id
```
