---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_upstream"
sidebar_current: "docs-tencentcloud-resource-api_gateway_upstream"
description: |-
  Provides a resource to create a apigateway upstream
---

# tencentcloud_api_gateway_upstream

Provides a resource to create a apigateway upstream

## Example Usage

### Create a basic VPC channel

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cvm"
}

data "tencentcloud_images" "images" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

data "tencentcloud_instance_types" "instance_types" {
  filter {
    name   = "instance-family"
    values = ["S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.3.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_instance" "example" {
  instance_name     = "tf_example"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.3.name
  image_id          = data.tencentcloud_images.images.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  system_disk_size  = 50
  hostname          = "terraform"
  project_id        = 0
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

resource "tencentcloud_api_gateway_upstream" "example" {
  scheme               = "HTTP"
  algorithm            = "ROUND-ROBIN"
  uniq_vpc_id          = tencentcloud_vpc.vpc.id
  upstream_name        = "tf_example"
  upstream_description = "desc."
  upstream_type        = "IP_PORT"
  retries              = 5

  nodes {
    host           = "1.1.1.1"
    port           = 9090
    weight         = 10
    vm_instance_id = tencentcloud_instance.example.id
    tags           = ["tags"]
  }

  tags = {
    "createdBy" = "terraform"
  }
}
```

### Create a complete VPC channel

```hcl
resource "tencentcloud_api_gateway_upstream" "example" {
  scheme               = "HTTP"
  algorithm            = "ROUND-ROBIN"
  uniq_vpc_id          = tencentcloud_vpc.vpc.id
  upstream_name        = "tf_example"
  upstream_description = "desc."
  upstream_type        = "IP_PORT"
  retries              = 5

  nodes {
    host           = "1.1.1.1"
    port           = 9090
    weight         = 10
    vm_instance_id = tencentcloud_instance.example.id
    tags           = ["tags"]
  }

  health_checker {
    enable_active_check    = true
    enable_passive_check   = true
    healthy_http_status    = "200"
    unhealthy_http_status  = "500"
    tcp_failure_threshold  = 5
    timeout_threshold      = 5
    http_failure_threshold = 3
    active_check_http_path = "/"
    active_check_timeout   = 5
    active_check_interval  = 5
    unhealthy_timeout      = 30
  }

  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `algorithm` - (Required, String) Load balancing algorithm, value range: ROUND-ROBIN.
* `scheme` - (Required, String) Backend protocol, value range: HTTP, HTTPS, gRPC, gRPCs.
* `uniq_vpc_id` - (Required, String) VPC Unique ID.
* `health_checker` - (Optional, List) Health check configuration, currently only supports VPC channels.
* `k8s_service` - (Optional, List) Configuration of K8S container service.
* `nodes` - (Optional, List) Backend nodes.
* `retries` - (Optional, Int) Request retry count, default to 3 times.
* `tags` - (Optional, Map) Tag description list.
* `upstream_description` - (Optional, String) Backend channel description.
* `upstream_host` - (Optional, String) Host request header forwarded by gateway to backend.
* `upstream_name` - (Optional, String) Backend channel name.
* `upstream_type` - (Optional, String) Backend access type, value range: IP_PORT, K8S.

The `extra_labels` object of `k8s_service` supports the following:

* `key` - (Required, String) Key of Label.
* `value` - (Required, String) Value of Label.

The `health_checker` object supports the following:

* `enable_active_check` - (Required, Bool) Identify whether active health checks are enabled.
* `enable_passive_check` - (Required, Bool) Identify whether passive health checks are enabled.
* `healthy_http_status` - (Required, String) The HTTP status code that determines a successful request during a health check.
* `http_failure_threshold` - (Required, Int) HTTP continuous error threshold. 0 means HTTP checking is disabled. Value range: [0, 254].
* `tcp_failure_threshold` - (Required, Int) TCP continuous error threshold. 0 indicates disabling TCP checking. Value range: [0, 254].
* `timeout_threshold` - (Required, Int) Continuous timeout threshold. 0 indicates disabling timeout checking. Value range: [0, 254].
* `unhealthy_http_status` - (Required, String) The HTTP status code that determines a failed request during a health check.
* `active_check_http_path` - (Optional, String) Detect the requested path during active health checks. The default is&#39;/&#39;.
* `active_check_interval` - (Optional, Int) The time interval for active health checks is 5 seconds by default.
* `active_check_timeout` - (Optional, Int) The detection request for active health check timed out in seconds. The default is 5 seconds.
* `unhealthy_timeout` - (Optional, Int) The automatic recovery time of abnormal node status, in seconds. When only passive checking is enabled, it must be set to a value&gt;0, otherwise the passive exception node will not be able to recover. The default is 30 seconds.

The `k8s_service` object supports the following:

* `cluster_id` - (Required, String) K8s cluster ID.
* `extra_labels` - (Required, List) Additional Selected Pod Label.
* `namespace` - (Required, String) Container namespace.
* `port` - (Required, Int) Port of service.
* `service_name` - (Required, String) The name of the container service.
* `weight` - (Required, Int) weight.
* `name` - (Optional, String) Customized service name, optional.

The `nodes` object supports the following:

* `host` - (Required, String) IP or domain name.
* `port` - (Required, Int) Port [0, 65535].
* `weight` - (Required, Int) Weight [0, 100], 0 is disabled.
* `cluster_id` - (Optional, String) The ID of the TKE clusterNote: This field may return null, indicating that a valid value cannot be obtained.
* `name_space` - (Optional, String) K8S namespaceNote: This field may return null, indicating that a valid value cannot be obtained.
* `service_name` - (Optional, String) K8S container service nameNote: This field may return null, indicating that a valid value cannot be obtained.
* `source` - (Optional, String) Source of Node, value range: K8SNote: This field may return null, indicating that a valid value cannot be obtained.
* `tags` - (Optional, Set) Dye labelNote: This field may return null, indicating that a valid value cannot be obtained.
* `unique_service_name` - (Optional, String) Unique service name recorded internally by API gatewayNote: This field may return null, indicating that a valid value cannot be obtained.
* `vm_instance_id` - (Optional, String) CVM instance IDNote: This field may return null, indicating that a valid value cannot be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

apigateway upstream can be imported using the id, e.g.

```
terraform import tencentcloud_api_gateway_upstream.upstream upstream_id
```

