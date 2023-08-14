---
subcategory: "EventBridge(EB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eb_event_connector"
sidebar_current: "docs-tencentcloud-resource-eb_event_connector"
description: |-
  Provides a resource to create a eb event_connector
---

# tencentcloud_eb_event_connector

Provides a resource to create a eb event_connector

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

resource "tencentcloud_api_gateway_service" "service" {
  service_name = "tf-eb-service"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

locals {
  service_id = tencentcloud_api_gateway_service.service.id
}

resource "tencentcloud_eb_event_connector" "event_connector" {
  event_bus_id    = tencentcloud_eb_event_bus.foo.id
  connection_name = "tf-event-connector"
  description     = "event connector desc1"
  enable          = false
  type            = "apigw"
  connection_description {
    resource_description = "qcs::apigw:ap-guangzhou:uin/100022975249:serviceid/${local.service_id}"
    api_gw_params {
      protocol = "HTTP"
      method   = "GET"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `connection_description` - (Required, List) Connector description.
* `connection_name` - (Required, String) connector name.
* `event_bus_id` - (Required, String, ForceNew) event bus Id.
* `description` - (Optional, String) description.
* `enable` - (Optional, Bool) switch.
* `type` - (Optional, String) type.

The `api_gw_params` object supports the following:

* `method` - (Required, String) POST.
* `protocol` - (Required, String) HTTPS.

The `ckafka_params` object supports the following:

* `offset` - (Required, String) kafka offset.
* `topic_name` - (Required, String) ckafka  topic.

The `connection_description` object supports the following:

* `resource_description` - (Required, String) Resource qcs six-segment style, more reference [resource six-segment style](https://cloud.tencent.com/document/product/598/10606).
* `api_gw_params` - (Optional, List) apigw parameter,Note: This field may return null, indicating that no valid value can be obtained.
* `ckafka_params` - (Optional, List) ckafka parameter, note: this field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

eb event_connector can be imported using the id, e.g.

```
terraform import tencentcloud_eb_event_connector.event_connector event_connector_id
```

