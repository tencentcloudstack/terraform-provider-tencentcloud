Provides a resource to create a eb eb_transform

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

resource "tencentcloud_eb_event_rule" "foo" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  rule_name    = "tf-event_rule"
  description  = "event rule desc"
  enable       = true
  event_pattern = jsonencode(
    {
      source = "apigw.cloud.tencent"
      type = [
        "connector:apigw",
      ]
    }
  )
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_eb_event_transform" "foo" {
    event_bus_id = tencentcloud_eb_event_bus.foo.id
    rule_id      = tencentcloud_eb_event_rule.foo.rule_id

    transformations {

        extraction {
            extraction_input_path = "$"
            format                = "JSON"
        }

        transform {
            output_structs {
                key        = "type"
                value      = "connector:ckafka"
                value_type = "STRING"
            }
            output_structs {
                key        = "source"
                value      = "ckafka.cloud.tencent"
                value_type = "STRING"
            }
            output_structs {
                key        = "region"
                value      = "ap-guangzhou"
                value_type = "STRING"
            }
            output_structs {
                key        = "datacontenttype"
                value      = "application/json;charset=utf-8"
                value_type = "STRING"
            }
            output_structs {
                key        = "status"
                value      = "-"
                value_type = "STRING"
            }
            output_structs {
                key        = "data"
                value      = jsonencode(
                    {
                        Partition = 1
                        msgBody   = "Hello from Ckafka again!"
                        msgKey    = "test"
                        offset    = 37
                        topic     = "test-topic"
                    }
                )
                value_type = "STRING"
            }
        }
    }
}
```

Import

eb eb_transform can be imported using the id, e.g.

```
terraform import tencentcloud_eb_event_transform.eb_transform eb_transform_id
```