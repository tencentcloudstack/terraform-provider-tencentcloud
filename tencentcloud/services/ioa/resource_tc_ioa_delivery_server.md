Provides a resource to create or update an IOA (iOA) log delivery server.

The create and update paths reuse the `ModifyDeliveryServer` API: on create the
server id inside `data` is ignored and the API-generated id is returned and
tracked as the resource `id`; on update the create-time id is sent as `Data.ID`.

Example Usage

Syslog delivery server

```hcl
resource "tencentcloud_ioa_delivery_server" "example" {
  data = {
    name                   = "tf-example"
    type                   = 4
    description            = "syslog delivery server managed by terraform"
    target_syslog_address  = "http://www.example.com"
    target_syslog_protocol = "tcp"
    target_syslog_port     = 601
    syslog_facility        = "local0"
    target_http_timeout    = 60
  }
}
```

Kafka delivery server

```hcl
resource "tencentcloud_ioa_delivery_server" "example" {
  data = {
    name                 = "tf-example"
    type                 = 2
    description          = "kafka delivery server managed by terraform"
    target_kafka_address = "106.55.81.5:50000"
    target_kafka_topic   = "test_topic"
    kafka_version        = "2.8.1"
    kafka_sasl_user_name = "ckafka-9jdm9zxe#root"
    kafka_sasl_password  = "Qwer!234"
    kafka_ssl_enable     = false
    ssl_skip_verify      = false
  }
}
```

HTTP delivery server

```hcl
resource "tencentcloud_ioa_delivery_server" "example" {
  data = {
    name                = "tf-example"
    type                = 3
    description         = "http delivery server managed by terraform"
    target_http_url     = "http://www.example.com"
    target_http_timeout = 60
  }
}
```

Import

IOA delivery server can be imported using the server id, e.g.

```
terraform import tencentcloud_ioa_delivery_server.example 325
```
