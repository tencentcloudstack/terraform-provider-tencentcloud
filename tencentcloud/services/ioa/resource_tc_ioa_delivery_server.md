---
subcategory: ""
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ioa_delivery_server"
description: |-
  Provides a resource to create or update an IOA log delivery server.
---

# tencentcloud_ioa_delivery_server

Provides a resource to create or update an IOA (iOA) log delivery server.

The create and update paths reuse the `ModifyDeliveryServer` API: on create the
server id inside `data` is ignored and the API-generated id is returned and
tracked as the resource `id`; on update the create-time id is sent as `Data.ID`.

Example Usage

```hcl
resource "tencentcloud_ioa_delivery_server" "example" {
  data {
    name                   = "tf-example-syslog"
    type                   = 4
    description            = "syslog delivery server managed by terraform"
    target_syslog_address  = "10.0.0.1"
    target_syslog_protocol = "udp"
    target_syslog_port     = 514
    syslog_facility        = "local0"
  }
}
```

Kafka delivery server:

```hcl
resource "tencentcloud_ioa_delivery_server" "kafka" {
  data {
    name                 = "tf-example-kafka"
    type                 = 2
    description          = "kafka delivery server managed by terraform"
    target_kafka_address = "10.0.0.2:9092"
    target_kafka_topic   = "ioa-log"
    kafka_version        = "2.8.0"
    kafka_sasl_user_name = "sasl-user"
    kafka_sasl_password  = "sasl-pass"
    kafka_ssl_enable     = true
    ssl_skip_verify      = false
  }
}
```

HTTP delivery server:

```hcl
resource "tencentcloud_ioa_delivery_server" "http" {
  data {
    name                 = "tf-example-http"
    type                 = 3
    description          = "http delivery server managed by terraform"
    target_http_url      = "http://10.0.0.3/collect"
    target_http_timeout  = 3000
  }
}
```

## Import

IOA delivery server can be imported using the server id, e.g.

```
terraform import tencentcloud_ioa_delivery_server.example 12345
```
