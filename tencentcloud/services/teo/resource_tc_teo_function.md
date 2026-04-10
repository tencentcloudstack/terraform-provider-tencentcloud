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

    environment_variables {
        key   = "API_KEY"
        value = "your-api-key"
        type  = "string"
    }

    rules {
        function_rule_conditions {
            rule_conditions {
                target   = "url"
                operator = "equals"
                values   = ["/api"]
            }
        }
    }
}
```

Example Usage with Multiple Fields

```hcl
resource "tencentcloud_teo_function" "teo_function" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World!!');
          e.respondWith(response);
        });
    EOT
    name        = "aaa-zone-2qtuhspy7cr6-1310708577"
    remark      = "test with all features"
    zone_id     = "zone-2qtuhspy7cr6"

    environment_variables {
        key   = "ENV1"
        value = "value1"
        type  = "string"
    }

    environment_variables {
        key   = "ENV2"
        value = "{\"config\":\"nested\"}"
        type  = "json"
    }

    rules {
        function_rule_conditions {
            rule_conditions {
                target   = "url"
                operator = "equals"
                values   = ["/api/v1"]
            }
        }
    }

    region_selection = ["CN", "US"]
}
```

Import

teo teo_function can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function.teo_function zone_id#function_id
```
