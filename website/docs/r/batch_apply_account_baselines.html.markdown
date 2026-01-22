---
subcategory: "ControlCenter"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_batch_apply_account_baselines"
sidebar_current: "docs-tencentcloud-resource-batch_apply_account_baselines"
description: |-
  Provides a resource to create a Controlcenter batch apply account baselines
---

# tencentcloud_batch_apply_account_baselines

Provides a resource to create a Controlcenter batch apply account baselines

## Example Usage

```hcl
resource "tencentcloud_batch_apply_account_baselines" "example" {
  member_uin_list = [
    10037652245,
    10037652240,
  ]

  baseline_config_items {
    identifier    = "TCC-AF_SHARE_IMAGE"
    configuration = "{\"Images\":[{\"Region\":\"ap-guangzhou\",\"ImageId\":\"img-mcdsiqrx\",\"ImageName\":\"demo1\"}, {\"Region\":\"ap-guangzhou\",\"ImageId\":\"img-esxgkots\",\"ImageName\":\"demo2\"}]}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `baseline_config_items` - (Required, List, ForceNew) List of baseline item configuration information.
* `member_uin_list` - (Required, Set: [`Int`], ForceNew) Member account UIN, which is also the UIN of the account to which the baseline is applied.

The `baseline_config_items` object supports the following:

* `configuration` - (Optional, String) Account Factory baseline item configuration. Different items have different parameters.Note: This field may return null, indicating that no valid values can be obtained.
* `identifier` - (Optional, String) A unique identifier for an Account Factory baseline item, which can only contain English letters, digits, and @,._[]-:()+=. It must be 2-128 characters long.Note: This field may return null, indicating that no valid values can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



