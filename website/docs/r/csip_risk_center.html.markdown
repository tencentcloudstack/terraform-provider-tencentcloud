---
subcategory: "CSIP"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_csip_risk_center"
sidebar_current: "docs-tencentcloud-resource-csip_risk_center"
description: |-
  Provides a resource to create a csip risk_center
---

# tencentcloud_csip_risk_center

Provides a resource to create a csip risk_center

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `scan_asset_type` - (Required, Int) 0- Full scan, 1- Specify asset scan, 2- Exclude asset scan, 3- Manually fill in the scan. If 1 and 2 are required, the Assets field is required. If 3 is required, SelfDefiningAssets is required.
* `scan_item` - (Required, Set: [`String`]) Scan Project. Example: port/poc/weakpass/webcontent/configrisk/exposedserver.
* `scan_plan_type` - (Required, Int) 0- Periodic task,1- immediate scan,2- periodic scan,3- Custom; 0,2, and 3 are required for ScanPlanContent.
* `task_name` - (Required, String) Task Name.
* `assets` - (Optional, List) Scan the asset information list.
* `scan_from` - (Optional, String) Request origin. The default value vss indicates the vulnerability scanning service. Users of the cloud security center please fill in the csip.
* `scan_plan_content` - (Optional, String) Scan plan details.
* `self_defining_assets` - (Optional, Set: [`String`]) Ip/domain/url array.
* `task_advance_cfg` - (Optional, List) Advanced configuration.
* `task_mode` - (Optional, Int) Physical examination mode, 0-standard mode, 1-fast mode, 2-advanced mode, default standard mode.

The `assets` object supports the following:

* `arn` - (Optional, String) Multi-cloud asset unique idNote: This field may return null, indicating that a valid value cannot be obtained.
* `asset_name` - (Optional, String) Asset nameNote: This field may return null, indicating that a valid value cannot be obtained.
* `asset_type` - (Optional, String) Asset classificationNote: This field may return null, indicating that a valid value cannot be obtained.
* `asset` - (Optional, String) Ip/ domain name/asset id, database id, etc.
* `instance_type` - (Optional, String) Asset typeNote: This field may return null, indicating that a valid value cannot be obtained.
* `region` - (Optional, String) RegionNote: This field may return null, indicating that a valid value cannot be obtained.

The `cfg_risk` object of `task_advance_cfg` supports the following:

* `enable` - (Required, Int) Whether to enable, 0- No, 1- Enable.
* `item_id` - (Required, String) Detection item ID.
* `resource_type` - (Required, String) Resource type.

The `task_advance_cfg` object supports the following:

* `cfg_risk` - (Optional, List) Configure advanced risk Settings.
* `vul_risk` - (Optional, List) Advanced vulnerability risk configuration.
* `weak_pwd_risk` - (Optional, List) Weak password risk advanced configuration.

The `vul_risk` object of `task_advance_cfg` supports the following:

* `enable` - (Required, Int) Whether to enable, 0- No, 1- Enable.
* `risk_id` - (Required, String) Risk ID.

The `weak_pwd_risk` object of `task_advance_cfg` supports the following:

* `check_item_id` - (Required, Int) Detection item ID.
* `enable` - (Required, Int) Whether to enable, 0- No, 1- Enable.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



