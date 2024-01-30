Provides a resource to create a csip risk_center

Example Usage

If task_mode is 0

```hcl
resource "tencentcloud_csip_risk_center" "example" {
  task_name         = "tf_example"
  scan_asset_type   = 1
  scan_item         = ["port", "poc", "weakpass"]
  scan_plan_content = "46 51 16 */1 * * *"
  task_mode         = 0

  assets {
    asset_name    = "iac-test"
    instance_type = "1"
    asset_type    = "PublicIp"
    asset         = "49.232.172.248"
    region        = "ap-beijing"
  }
}
```

If task_mode is 1

```hcl
resource "tencentcloud_csip_risk_center" "example" {
  task_name         = "tf_example"
  scan_asset_type   = 3
  scan_item         = ["port", "poc"]
  scan_plan_type    = 0
  scan_plan_content = "46 51 16 */1 * * *"
  task_mode         = 1
}
```

If task_mode is 2

```hcl
resource "tencentcloud_csip_risk_center" "example" {
  task_name       = "tf_example"
  scan_asset_type = 2
  scan_item       = ["port", "poc"]
  task_mode       = 2

  assets {
    asset_name    = "iac-test"
    instance_type = "1"
    asset_type    = "PublicIp"
    asset         = "49.232.172.248"
    region        = "ap-beijing"
  }

  assets {
    asset_name    = "iac-test"
    instance_type = "POSTGRES"
    asset_type    = "Db"
    asset         = "postgres-fnexv5bj"
    region        = "ap-guangzhou"
  }

  task_advance_cfg {
    port_risk {
      check_type = 0
      detail     = "22、8080、80、443、3380、3389常见流量端口"
      port_sets  = "常见端口"
      enable     = 1
    }
    vul_risk {
      risk_id = "b52a4fcc1f24fa323b87cc41f370aa43"
      enable  = 1
    }
  }
}
```
