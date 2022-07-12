---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_api_key_attachment"
sidebar_current: "docs-tencentcloud-resource-api_gateway_api_key_attachment"
description: |-
  Use this resource to API gateway attach access key to usage plan.
---

# tencentcloud_api_gateway_api_key_attachment

Use this resource to API gateway attach access key to usage plan.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_api_key" "key" {
  secret_name = "my_api_key"
  status      = "on"
}

resource "tencentcloud_api_gateway_usage_plan" "plan" {
  usage_plan_name         = "my_plan"
  usage_plan_desc         = "nice plan"
  max_request_num         = 100
  max_request_num_pre_sec = 10
}

resource "tencentcloud_api_gateway_api_key_attachment" "attach" {
  api_key_id    = tencentcloud_api_gateway_api_key.key.id
  usage_plan_id = tencentcloud_api_gateway_usage_plan.plan.id
}
```

## Argument Reference

The following arguments are supported:

* `api_key_id` - (Required, String, ForceNew) ID of API key.
* `usage_plan_id` - (Required, String, ForceNew) ID of the usage plan.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

API gateway attach access key can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_api_key_attachment.attach AKID110b8Rmuw7t0fP1N8bi809n327023Is7xN8f#usagePlan-gyeafpab
```

