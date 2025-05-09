Provides a resource to create a MQTT authorization policy

Example Usage

```hcl
resource "tencentcloud_mqtt_authorization_policy" "example" {
  instance_id    = "mqtt-g4qgr3gx"
  policy_name    = "tf-example"
  policy_version = 1
  priority       = 10
  effect         = "allow"
  actions        = "connect,pub,sub"
  retain         = 3
  qos            = "0,1,2"
  resources      = "topic-demo"
  username       = "*root*"
  client_id      = "client"
  ip             = "192.168.1.1"
  remark         = "policy remark."
}
```

Or

```
resource "tencentcloud_mqtt_authorization_policy" "example" {
  instance_id    = "mqtt-g4qgr3gx"
  policy_name    = "tf-example"
  policy_version = 1
  priority       = 10
  effect         = "deny"
  actions        = "pub,sub"
  retain         = 3
  qos            = "1,2"
  resources      = "topic-demo"
  username       = "root*"
  client_id      = "*$${Username}*"
  ip             = "192.168.1.0/24"
  remark         = "policy remark."
}
```

Import

MQTT authorization policy can be imported using the id, e.g.

```
terraform import tencentcloud_mqtt_authorization_policy.example mqtt-g4qgr3gx#140
```
