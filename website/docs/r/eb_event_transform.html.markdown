---
subcategory: "EventBridge(EB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eb_event_transform"
sidebar_current: "docs-tencentcloud-resource-eb_event_transform"
description: |-
  Provides a resource to create a eb eb_transform
---

# tencentcloud_eb_event_transform

Provides a resource to create a eb eb_transform

## Example Usage

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
        key = "data"
        value = jsonencode(
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

## Argument Reference

The following arguments are supported:

* `event_bus_id` - (Required, String) event bus Id.
* `rule_id` - (Required, String) ruleId.
* `transformations` - (Required, List) A list of transformation rules, currently only one.

The `etl_filter` object supports the following:

* `filter` - (Required, String) Grammatical Rules are consistent.

The `extraction` object supports the following:

* `extraction_input_path` - (Required, String) JsonPath, if not specified, the default value $.
* `format` - (Required, String) Value: `TEXT`, `JSON`.
* `text_params` - (Optional, List) Only Text needs to be passed.

The `output_structs` object supports the following:

* `key` - (Required, String) Corresponding to the key in the output json.
* `value_type` - (Required, String) The data type of value, optional values: `STRING`, `NUMBER`, `BOOLEAN`, `NULL`, `SYS_VARIABLE`, `JSONPATH`.
* `value` - (Required, String) You can fill in the json-path and also support constants or built-in keyword date types.

The `text_params` object supports the following:

* `regex` - (Optional, String) Fill in the regular expression: length 128.
* `separator` - (Optional, String) `Comma`, `|`, `tab`, `space`, `newline`, `%`, `#`, the limit length is 1.

The `transform` object supports the following:

* `output_structs` - (Required, List) Describe how the data is transformed.

The `transformations` object supports the following:

* `etl_filter` - (Optional, List) Describe how to filter data.
* `extraction` - (Optional, List) Describe how to extract data.
* `transform` - (Optional, List) Describe how to convert data.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

eb eb_transform can be imported using the id, e.g.

```
terraform import tencentcloud_eb_event_transform.eb_transform eb_transform_id
```

