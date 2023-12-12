Provides a resource to create a eb event_connector

~> **NOTE:** When the type is `apigw`, the import function is not supported.

Example Usage

Create ckafka event connector
```hcl
data "tencentcloud_user_info" "foo" {}

resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_ckafka_instance" "kafka_instance" {
  instance_name      = "ckafka-instance-maz-tf-test"
  zone_id            = 100003
  multi_zone_flag    = true
  zone_ids           = [100003, 100006]
  period             = 1
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  msg_retention_time = 1300
  renew_flag         = 0
  kafka_version      = "1.1.1"
  disk_size          = 500
  disk_type          = "CLOUD_BASIC"

  config {
    auto_create_topic_enable   = true
    default_num_partitions     = 3
    default_replication_factor = 3
  }

  dynamic_retention_config {
    enable = 1
  }
}

locals {
  ckafka_id = tencentcloud_ckafka_instance.kafka_instance.id
  uin = data.tencentcloud_user_info.foo.owner_uin
}

resource "tencentcloud_eb_event_connector" "event_connector" {
  event_bus_id    = tencentcloud_eb_event_bus.foo.id
  connection_name = "tf-event-connector"
  description     = "event connector desc1"
  enable          = true
  type            = "ckafka"
  connection_description {
    resource_description = "qcs::ckafka:ap-guangzhou:uin/${local.uin}:ckafkaId/uin/${local.uin}/${local.ckafka_id}"
    ckafka_params {
      offset     = "latest"
      topic_name = "dasdasd"
    }
  }
}
```

Create api_gateway event connector

```hcl
data "tencentcloud_user_info" "foo" {}

resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_api_gateway_service" "service" {
  service_name = "tf-eb-service"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

locals {
  uin = data.tencentcloud_user_info.foo.owner_uin
  service_id = tencentcloud_api_gateway_service.service.id
}

resource "tencentcloud_eb_event_connector" "event_connector" {
  event_bus_id    = tencentcloud_eb_event_bus.foo.id
  connection_name = "tf-event-connector"
  description     = "event connector desc1"
  enable          = false
  type            = "apigw"

  connection_description {
    resource_description = "qcs::apigw:ap-guangzhou:uin/${local.uin}:serviceid/${local.service_id}"
    api_gw_params {
      protocol = "HTTP"
      method   = "GET"
    }
  }
}
```

Import

eb event_connector can be imported using the id, e.g.

```
terraform import tencentcloud_eb_event_connector.event_connector eventBusId#connectionId
```