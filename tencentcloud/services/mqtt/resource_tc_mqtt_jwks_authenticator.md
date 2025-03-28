Provides a resource to create a MQTT jwks authenticator

Example Usage

```hcl
resource "tencentcloud_mqtt_jwks_authenticator" "example" {
  instance_id      = "mqtt-zxjwkr98"
  from             = "username"
  endpoint         = "https://example.com"
  refresh_interval = 60
  remark           = "Remark."
}
```

Or

```hcl
resource "tencentcloud_mqtt_jwks_authenticator" "example" {
  instance_id = "mqtt-zxjwkr98"
  from        = "password"
  text        = "your text content"
  remark      = "Remark."
}
```

Import

MQTT jwks authenticator can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_jwks_authenticator.example mqtt-zxjwkr98
```
