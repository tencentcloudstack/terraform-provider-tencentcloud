---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_provisioned_concurrency_config"
sidebar_current: "docs-tencentcloud-resource-scf_provisioned_concurrency_config"
description: |-
  Provides a resource to create a scf provisioned_concurrency_config
---

# tencentcloud_scf_provisioned_concurrency_config

Provides a resource to create a scf provisioned_concurrency_config

## Example Usage

```hcl
resource "tencentcloud_scf_provisioned_concurrency_config" "provisioned_concurrency_config" {
  function_name                       = "keep-1676351130"
  qualifier                           = "2"
  version_provisioned_concurrency_num = 2
  namespace                           = "default"
  trigger_actions {
    trigger_name                        = "test"
    trigger_provisioned_concurrency_num = 2
    trigger_cron_config                 = "29 45 12 29 05 * 2023"
    provisioned_type                    = "Default"
  }
  provisioned_type = "Default"
  tracking_target  = 0.5
  min_capacity     = 1
  max_capacity     = 2
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String, ForceNew) Name of the function for which to set the provisioned concurrency.
* `qualifier` - (Required, String, ForceNew) Function version number. Note: the $LATEST version does not support provisioned concurrency.
* `version_provisioned_concurrency_num` - (Required, Int, ForceNew) Provisioned concurrency amount. Note: there is an upper limit for the sum of provisioned concurrency amounts of all versions, which currently is the function&amp;#39;s maximum concurrency quota minus 100.
* `max_capacity` - (Optional, Int, ForceNew) The maximum number of instances.
* `min_capacity` - (Optional, Int, ForceNew) The minimum number of instances. It can not be smaller than 1.
* `namespace` - (Optional, String, ForceNew) Function namespace. Default value: default.
* `provisioned_type` - (Optional, String, ForceNew) Specifies the provisioned concurrency type. Default: Static provisioned concurrency. ConcurrencyUtilizationTracking: Scales the concurrency automatically according to the concurrency utilization. If ConcurrencyUtilizationTracking is passed in, TrackingTarget, MinCapacity and MaxCapacity are required, and VersionProvisionedConcurrencyNum must be 0.
* `tracking_target` - (Optional, Float64, ForceNew) The target concurrency utilization. Range: (0,1) (two decimal places).
* `trigger_actions` - (Optional, List, ForceNew) Scheduled provisioned concurrency scaling action.

The `trigger_actions` object supports the following:

* `trigger_cron_config` - (Required, String) Trigger time of the scheduled action in Cron expression. Seven fields are required and should be separated with a space. Note: this field may return null, indicating that no valid values can be obtained.
* `trigger_name` - (Required, String) Scheduled action name Note: this field may return null, indicating that no valid values can be obtained.
* `trigger_provisioned_concurrency_num` - (Required, Int) Target provisioned concurrency of the scheduled scaling action Note: this field may return null, indicating that no valid values can be obtained.
* `provisioned_type` - (Optional, String) The provision type. Value: Default Note: This field may return null, indicating that no valid value can be found.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



