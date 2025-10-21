---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_cngw_service_rate_limit"
sidebar_current: "docs-tencentcloud-resource-tse_cngw_service_rate_limit"
description: |-
  Provides a resource to create a tse cngw_service_rate_limit
---

# tencentcloud_tse_cngw_service_rate_limit

Provides a resource to create a tse cngw_service_rate_limit

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_tse_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_tse_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_tse_cngw_gateway" "cngw_gateway" {
  description                = "terraform test1"
  enable_cls                 = true
  engine_region              = "ap-guangzhou"
  feature_version            = "STANDARD"
  gateway_version            = "2.5.1"
  ingress_class_name         = "tse-nginx-ingress"
  internet_max_bandwidth_out = 0
  name                       = "terraform-gateway1"
  trade_type                 = 0
  type                       = "kong"

  node_config {
    number        = 2
    specification = "1c2g"
  }

  vpc_config {
    subnet_id = tencentcloud_subnet.subnet.id
    vpc_id    = tencentcloud_vpc.vpc.id
  }

  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tse_cngw_service" "cngw_service" {
  gateway_id    = tencentcloud_tse_cngw_gateway.cngw_gateway.id
  name          = "terraform-test"
  path          = "/test"
  protocol      = "http"
  retries       = 5
  timeout       = 60000
  upstream_type = "HostIP"

  upstream_info {
    algorithm             = "round-robin"
    auto_scaling_cvm_port = 0
    host                  = "arunma.cn"
    port                  = 8012
    slow_start            = 0
  }
}

resource "tencentcloud_tse_cngw_service_rate_limit" "cngw_service_rate_limit" {
  gateway_id = tencentcloud_tse_cngw_gateway.cngw_gateway.id
  name       = tencentcloud_tse_cngw_service.cngw_service.name

  limit_detail {
    enabled             = true
    header              = "req"
    hide_client_headers = true
    is_delay            = true
    limit_by            = "header"
    line_up_time        = 15
    policy              = "redis"
    response_type       = "default"

    qps_thresholds {
      max  = 100
      unit = "hour"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String) gateway ID.
* `limit_detail` - (Required, List) rate limit configuration.
* `name` - (Required, String) service name or service ID.

The `external_redis` object of `limit_detail` supports the following:

* `redis_host` - (Required, String) redis ip, maybe null.
* `redis_password` - (Required, String) redis password, maybe null.
* `redis_port` - (Required, Int) redis port, maybe null.
* `redis_timeout` - (Required, Int) redis timeout, unit: `ms`, maybe null.

The `headers` object of `rate_limit_response` supports the following:

* `key` - (Optional, String) key of header.
* `value` - (Optional, String) value of header.

The `limit_detail` object supports the following:

* `enabled` - (Required, Bool) status of service rate limit.
* `hide_client_headers` - (Required, Bool) whether to hide the headers of client.
* `is_delay` - (Required, Bool) whether to enable request queuing.
* `limit_by` - (Required, String) basis for service rate limit.Reference value: `ip`, `service`, `consumer`, `credential`, `path`, `header`.
* `qps_thresholds` - (Required, List) qps threshold.
* `response_type` - (Required, String) response strategy.Reference value: `url`: forward request according to url, `text`: response configuration, `default`: return directly.
* `external_redis` - (Optional, List) external redis information, maybe null.
* `header` - (Optional, String) request headers that require rate limit.
* `line_up_time` - (Optional, Int) queue time.
* `path` - (Optional, String) request paths that require rate limit.
* `policy` - (Optional, String) counter policy.Reference value: `local`, `redis`, `external_redis`.
* `rate_limit_response_url` - (Optional, String) request forwarding address, maybe null.
* `rate_limit_response` - (Optional, List) response configuration, the response strategy is text, maybe null.

The `qps_thresholds` object of `limit_detail` supports the following:

* `max` - (Required, Int) the max threshold.
* `unit` - (Required, String) qps threshold unit.Reference value:`second`, `minute`, `hour`, `day`, `month`, `year`.

The `rate_limit_response` object of `limit_detail` supports the following:

* `body` - (Optional, String) custom response body, maybe bull.
* `headers` - (Optional, List) headrs.
* `http_status` - (Optional, Int) http status code.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tse cngw_service_rate_limit can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit gatewayId#name
```

