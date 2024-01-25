Provides a resource to create a csip risk_center

Example Usage

```hcl
resource "tencentcloud_csip_risk_center" "example" {
  task_name       = "tf_example"
  scan_asset_type = 0
  scan_item       = []
  scan_plan_type  = 1
  assets {
    asset_name    = ""
    instance_type = ""
    asset_type    = ""
    asset         = ""
    region        = ""
    arn           = ""
  }
  scan_plan_content    = ""
  self_defining_assets = []
  scan_from            = ""
  task_advance_cfg {
    vul_risk {
      risk_id = ""
      enable  = 0
    }
    weak_pwd_risk {
      check_item_id = ""
      enable        = 0
    }
    cfg_risk {
      item_id       = ""
      enable        = 0
      resource_type = ""
    }
  }
  task_mode = 0
}
```
