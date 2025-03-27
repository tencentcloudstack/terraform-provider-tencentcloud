Provides a resource to create a MQTT jwt authenticator

Example Usage

If algorithm is hmac-based

```hcl
resource "tencentcloud_mqtt_jwt_authenticator" "example" {
  instance_id = "mqtt-zxjwkr98"
  algorithm   = "hmac-based"
  from        = "password"
  secret      = "your secret content"
  remark      = "Remark."
}
```

If algorithm is public-key

```hcl
resource "tencentcloud_mqtt_jwt_authenticator" "example" {
  instance_id = "mqtt-zxjwkr98"
  algorithm   = "public-key"
  from        = "username"
  public_key  = "your public key"
  remark      = "Remark."
}
```

Import

MQTT jwt authenticator can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_jwt_authenticator.example mqtt-zxjwkr98
```
