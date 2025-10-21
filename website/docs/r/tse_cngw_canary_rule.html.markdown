---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_cngw_canary_rule"
sidebar_current: "docs-tencentcloud-resource-tse_cngw_canary_rule"
description: |-
  Provides a resource to create a tse cngw_canary_rule
---

# tencentcloud_tse_cngw_canary_rule

Provides a resource to create a tse cngw_canary_rule

## Example Usage

```hcl
resource "tencentcloud_tse_cngw_service" "cngw_service" {
  gateway_id = "gateway-ddbb709b"
  name       = "terraform-test"
  path       = "/test"
  protocol   = "http"
  retries    = 5
  tags = {
    "created" = "terraform"
  }
  timeout       = 6000
  upstream_type = "IPList"

  upstream_info {
    algorithm                   = "round-robin"
    auto_scaling_cvm_port       = 80
    auto_scaling_group_id       = "asg-519acdug"
    auto_scaling_hook_status    = "Normal"
    auto_scaling_tat_cmd_status = "Normal"
    port                        = 0
    slow_start                  = 20

    targets {
      health = "HEALTHCHECKS_OFF"
      host   = "192.168.0.1"
      port   = 80
      weight = 100
    }
  }
}

resource "tencentcloud_tse_cngw_canary_rule" "cngw_canary_rule" {
  gateway_id = tencentcloud_tse_cngw_service.cngw_service.gateway_id
  service_id = tencentcloud_tse_cngw_service.cngw_service.service_id
  tags = {
    "created" = "terraform"
  }

  canary_rule {
    enabled  = true
    priority = 100

    balanced_service_list {
      percent      = 100
      service_id   = tencentcloud_tse_cngw_service.cngw_service.service_id
      service_name = tencentcloud_tse_cngw_service.cngw_service.name
    }

    condition_list {
      key      = "test"
      operator = "eq"
      type     = "query"
      value    = "1"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `canary_rule` - (Required, List) canary rule configuration.
* `gateway_id` - (Required, String) gateway ID.
* `service_id` - (Required, String) service ID.
* `tags` - (Optional, Map) Tag description list.

The `balanced_service_list` object of `canary_rule` supports the following:

* `percent` - (Optional, Float64) percent, 10 is 10%, valid values:0 to 100.
* `service_id` - (Optional, String) service ID, required when used as an input parameter.
* `service_name` - (Optional, String) service name, meaningless when used as an input parameter.

The `canary_rule` object supports the following:

* `enabled` - (Required, Bool) the status of canary rule.
* `priority` - (Required, Int, ForceNew) priority. The value ranges from 0 to 100; the larger the value, the higher the priority; the priority cannot be repeated between different rules.
* `balanced_service_list` - (Optional, List) service weight configuration.
* `condition_list` - (Optional, List) parameter matching condition list.
* `service_id` - (Optional, String) service ID.
* `service_name` - (Optional, String) service name.

The `condition_list` object of `canary_rule` supports the following:

* `type` - (Required, String) type.Reference value:`path`,`method`,`query`,`header`,`cookie`,`body`,`system`.
* `delimiter` - (Optional, String) delimiter. valid when operator is in or not in, reference value:`,`, `;`,`\n`.
* `global_config_id` - (Optional, String) global configuration ID.
* `global_config_name` - (Optional, String) global configuration name.
* `key` - (Optional, String) parameter name.
* `operator` - (Optional, String) operator.Reference value:`le`,`eq`,`lt`,`ne`,`ge`,`gt`,`regex`,`exists`,`in`,`not in`,`prefix`,`exact`,`regex`.
* `value` - (Optional, String) parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tse cngw_canary_rule can be imported using the gatewayId#serviceId#priority, e.g.

```
terraform import tencentcloud_tse_cngw_canary_rule.cngw_canary_rule gateway-ddbb709b#b6017eaf-2363-481e-9e93-8d65aaf498cd#100
```

