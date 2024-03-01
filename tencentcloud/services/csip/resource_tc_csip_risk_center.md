Provides a resource to create a csip risk_center

Example Usage

If task_mode is 0

```hcl
resource "tencentcloud_csip_risk_center" "example" {
  task_name         = "tf_example"
  scan_plan_type    = 0
  scan_asset_type   = 2
  scan_item         = ["port", "poc", "weakpass"]
  scan_plan_content = "0 0 0 */1 * * *"
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
  task_name       = "tf_example"
  scan_plan_type  = 1
  scan_asset_type = 1
  scan_item       = ["port", "poc"]
  task_mode       = 1
}
```

If task_mode is 2

```hcl
resource "tencentcloud_csip_risk_center" "example" {
  task_name       = "tf_example"
  scan_plan_type  = 2
  scan_asset_type = 2
  scan_item       = ["port", "configrisk", "poc", "weakpass"]
  task_mode       = 2
  scan_plan_content = "0 0 0 20 3 * 2024"

  assets {
    asset_name    = "sub machine of tke"
    instance_type = "Instance"
    asset_type    = "CVM"
    asset         = "ins-9p3dkkwy"
    region        = "ap-guangzhou"
  }

  task_advance_cfg {
    port_risk {
      check_type = 0
      detail     = "22、8080、80、443、3380、3389常见流量端"
      port_sets  = "常见端口"
      enable     = 1
    }
    vul_risk {
      risk_id = "f79e371ce5f644f0fdc72a143144c4b2"
      enable  = 1
    }
    weak_pwd_risk {
      check_item_id = 50
      enable        = 1
    }
    cfg_risk {
      item_id       = "02c9337f-a6da-49b4-8858-64663a02b79f"
      enable        = 1
      resource_type = "cdb;rds"
    }
  }
}
```
