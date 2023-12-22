Use this data source to query detailed information of eb eb_search

Example Usage

```hcl
resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_eb_put_events" "put_events" {
  event_list {
    source = "ckafka.cloud.tencent"
    data = jsonencode(
      {
        "topic" : "test-topic",
        "Partition" : 1,
        "offset" : 37,
        "msgKey" : "test",
        "msgBody" : "Hello from Ckafka again!"
      }
    )
    type    = "connector:ckafka"
    subject = "qcs::ckafka:ap-guangzhou:uin/1250000000:ckafkaId/uin/1250000000/ckafka-123456"
    time    = 1691572461939

  }
  event_bus_id = tencentcloud_eb_event_bus.foo.id
}

data "tencentcloud_eb_search" "eb_search" {
  start_time   = 1691637288422
  end_time     = 1691648088422
  event_bus_id = "eb-jzytzr4e"
  group_field = "RuleIds"
  filter {
  	type = "OR"
  	filters {
  		key = "status"
  		operator = "eq"
  		value = "1"
  	}
  }

  filter {
  	type = "OR"
  	filters {
  		key = "type"
  		operator = "eq"
  		value = "connector:ckafka"
  	}
  }
  # order_fields = [""]
  order_by = "desc"
}
```
