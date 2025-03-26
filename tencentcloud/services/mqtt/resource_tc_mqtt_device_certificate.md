Provides a resource to create a MQTT device certificate

Example Usage

```hcl
resource "tencentcloud_mqtt_device_certificate" "example" {
  instance_id        = "mqtt-zxjwkr98"
  device_certificate = ""
  ca_sn              = ""
  client_id          = ""
  format             = ""
  status             = "ACTIVE"
}
```

Import

MQTT device certificate can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_device_certificate.example mqtt_device_certificate_id
```
