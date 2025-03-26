Provides a resource to create a MQTT user

Example Usage

```hcl
resource "tencentcloud_mqtt_user" "example" {
  instance_id = "mqtt-zxjwkr98"
  username    = "tf-example"
  password    = "Password@123"
  remark      = "Remark."
}
```

Import

MQTT user can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_user.example mqtt-zxjwkr98#tf-example
```
