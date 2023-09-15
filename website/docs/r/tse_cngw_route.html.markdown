---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_cngw_route"
sidebar_current: "docs-tencentcloud-resource-tse_cngw_route"
description: |-
  Provides a resource to create a tse cngw_route
---

# tencentcloud_tse_cngw_route

Provides a resource to create a tse cngw_route

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

resource "tencentcloud_tse_cngw_route" "cngw_route" {
  destination_ports = []
  gateway_id        = tencentcloud_tse_cngw_gateway.cngw_gateway.id
  hosts = [
    "192.168.0.1:9090",
  ]
  https_redirect_status_code = 426
  paths = [
    "/user",
  ]
  headers {
    key   = "req"
    value = "terraform"
  }
  preserve_host = false
  protocols = [
    "http",
    "https",
  ]
  route_name = "terraform-route"
  service_id = tencentcloud_tse_cngw_service.cngw_service.service_id
  strip_path = true
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String) gateway ID.
* `service_id` - (Required, String) ID of the service which the route belongs to.
* `destination_ports` - (Optional, Set: [`Int`]) destination port for Layer 4 matching.
* `force_https` - (Optional, Bool, **Deprecated**) This field has been deprecated and will be deleted in subsequent versions. whether to enable forced HTTPS, no longer use.
* `headers` - (Optional, List) the headers of route.
* `hosts` - (Optional, Set: [`String`]) host list.
* `https_redirect_status_code` - (Optional, Int) https redirection status code.
* `methods` - (Optional, Set: [`String`]) route methods. Reference value:`GET`,`POST`,`DELETE`,`PUT`,`OPTIONS`,`PATCH`,`HEAD`,`ANY`,`TRACE`,`COPY`,`MOVE`,`PROPFIND`,`PROPPATCH`,`MKCOL`,`LOCK`,`UNLOCK`.
* `paths` - (Optional, Set: [`String`]) path list.
* `preserve_host` - (Optional, Bool) whether to keep the host when forwarding to the backend.
* `protocols` - (Optional, Set: [`String`]) the protocol list of route.Reference value:`https`,`http`.
* `route_name` - (Optional, String) the name of the route, unique in the instance.
* `strip_path` - (Optional, Bool) whether to strip path when forwarding to the backend.

The `headers` object supports the following:

* `key` - (Optional, String) key of header.
* `value` - (Optional, String) value of header.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `route_id` - the id of the route, unique in the instance.


