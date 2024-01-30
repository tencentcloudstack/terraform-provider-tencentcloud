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

### If task_mode is 1

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

### If task_mode is 2

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

## Argument Reference

The following arguments are supported:

* `scan_asset_type` - (Required, Int) 0- Full scan, 1- Specify asset scan, 2- Exclude asset scan, 3- Manually fill in the scan. If 1 and 2 are required, the Assets field is required. If 3 is required, SelfDefiningAssets is required.
* `scan_item` - (Required, Set: [`String`]) Scan Project. Example: port/poc/weakpass/webcontent/configrisk/exposedserver.
* `task_name` - (Required, String) Task Name.
* `assets` - (Optional, List) Scan the asset information list.
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


