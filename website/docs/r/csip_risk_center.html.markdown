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

### If task_mode is 0

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

### If task_mode is 1

```hcl
resource "tencentcloud_csip_risk_center" "example" {
  task_name       = "tf_example"
  scan_plan_type  = 1
  scan_asset_type = 1
  scan_item       = ["port", "poc"]
  task_mode       = 1
}
```

### If task_mode is 2

```hcl
resource "tencentcloud_csip_risk_center" "example" {
  task_name         = "tf_example"
  scan_plan_type    = 2
  scan_asset_type   = 2
  scan_item         = ["port", "configrisk", "poc", "weakpass"]
  task_mode         = 2
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

## Argument Reference

The following arguments are supported:

* `scan_asset_type` - (Required, Int, ForceNew) 0- Full scan, 1- Specify asset scan, 2- Exclude asset scan, 3- Manually fill in the scan. If 1 and 2 are required while task_mode not 1, the Assets field is required. If 3 is required, SelfDefiningAssets is required.
* `scan_item` - (Required, Set: [`String`], ForceNew) Scan Project. Example: port/poc/weakpass/webcontent/configrisk/exposedserver.
* `scan_plan_type` - (Required, Int, ForceNew) 0- Periodic task, 1- immediate scan, 2- periodic scan, 3- Custom; 0, 2 and 3 are required for scan_plan_content.
* `task_name` - (Required, String) Task Name.
* `assets` - (Optional, List, ForceNew) Scan the asset information list.
* `scan_plan_content` - (Optional, String, ForceNew) Scan plan details.
* `self_defining_assets` - (Optional, Set: [`String`], ForceNew) Ip/domain/url array.
* `task_advance_cfg` - (Optional, List, ForceNew) Advanced configuration.
* `task_mode` - (Optional, Int, ForceNew) Physical examination mode, 0-standard mode, 1-fast mode, 2-advanced mode, default standard mode.

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

The `port_risk` object of `task_advance_cfg` supports the following:

* `check_type` - (Required, Int) Detection item type, 0-system defined, 1-user-defined.
* `detail` - (Required, String) Description of detection items.
* `enable` - (Required, Int) Whether to enable, 0- No, 1- Enable.
* `port_sets` - (Required, String) Port collection, separated by commas.

The `task_advance_cfg` object supports the following:

* `cfg_risk` - (Optional, List) Configure advanced risk Settings.
* `port_risk` - (Optional, List) Advanced Port Risk Configuration.
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
* `scan_from` - Request origin.


