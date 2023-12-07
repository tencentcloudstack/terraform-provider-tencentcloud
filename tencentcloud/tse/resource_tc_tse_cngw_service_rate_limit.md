Provides a resource to create a tse cngw_service_rate_limit

Example Usage

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

Import

tse cngw_service_rate_limit can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_service_rate_limit.cngw_service_rate_limit gatewayId#name
```