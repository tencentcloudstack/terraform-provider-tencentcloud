Provides a resource to create a tse cngw_strategy

~> **NOTE:** Please pay attention to the correctness of the cycle when modifying the `params` of `cron_config`, otherwise the modification will not be successful.

Example Usage

```hcl
resource "tencentcloud_tse_cngw_strategy" "cngw_strategy" {
  description   = "aaaaa"
  gateway_id    = "gateway-cf8c99c3"
  strategy_name = "test-cron"

  config {
    max_replicas = 2
    behavior {
      scale_down {
        select_policy                = "Max"
        stabilization_window_seconds = 301

        policies {
          period_seconds = 9
          type           = "Pods"
          value          = 1
        }
      }
      scale_up {
        select_policy                = "Max"
        stabilization_window_seconds = 31

        policies {
          period_seconds = 10
          type           = "Pods"
          value          = 1
        }
      }
    }

    metrics {
      resource_name = "cpu"
      target_value  = 1
      type          = "Resource"
    }
  }

  cron_config {
    params {
      crontab         = "0 00 00 * * *"
      period          = "* * *"
      start_at        = "00:00"
      target_replicas = 2
    }
  }
}
```

Import

tse cngw_strategy can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_strategy.cngw_strategy gateway-cf8c99c3#strategy-a6744ff8
```
