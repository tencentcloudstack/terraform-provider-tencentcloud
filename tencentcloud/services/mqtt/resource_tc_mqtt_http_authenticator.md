Provides a resource to create a MQTT http authenticator

Example Usage

```hcl
resource "tencentcloud_mqtt_http_authenticator" "example" {
  instance_id     = "mqtt-zxjwkr98"
  endpoint        = "https://example.com"
  concurrency     = 8
  method          = "POST"
  status          = "open"
  remark          = "Remark."
  connect_timeout = 10
  read_timeout    = 10
  header {
    key   = "Content-type"
    value = "application/json"
  }

  body {
    key   = "bodyKey"
    value = "bodyValue"
  }
}
```

Import

MQTT http authenticator can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_http_authenticator.example mqtt-zxjwkr98
```
