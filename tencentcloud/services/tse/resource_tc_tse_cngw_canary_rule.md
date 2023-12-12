Provides a resource to create a tse cngw_canary_rule

Example Usage

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
  tags       = {
    "created" = "terraform"
  }

  canary_rule {
    enabled  = true
    priority = 100

    balanced_service_list {
      percent       = 100
      service_id    = tencentcloud_tse_cngw_service.cngw_service.service_id
      service_name  = tencentcloud_tse_cngw_service.cngw_service.name
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

Import

tse cngw_canary_rule can be imported using the gatewayId#serviceId#priority, e.g.

```
terraform import tencentcloud_tse_cngw_canary_rule.cngw_canary_rule gateway-ddbb709b#b6017eaf-2363-481e-9e93-8d65aaf498cd#100
```