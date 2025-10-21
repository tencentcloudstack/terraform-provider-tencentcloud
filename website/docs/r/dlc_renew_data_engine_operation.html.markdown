---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_renew_data_engine_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_renew_data_engine_operation"
description: |-
  Provides a resource to create a DLC renew data engine
---

# tencentcloud_dlc_renew_data_engine_operation

Provides a resource to create a DLC renew data engine

## Example Usage

```hcl
resource "tencentcloud_dlc_renew_data_engine_operation" "example" {
  data_engine_name = "tf-example"
  time_span        = 3600
  pay_mode         = 1
  time_unit        = "m"
  renew_flag       = 1
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_name` - (Required, String, ForceNew) CU queue name.
* `time_span` - (Required, Int, ForceNew) Renewal period in months, which is at least one month.
* `pay_mode` - (Optional, Int, ForceNew) Payment type. It is 1 by default and is prepaid.
* `renew_flag` - (Optional, Int, ForceNew) Auto-renewal flag: 0 means the initial status, and there is no automatic renewal by default. If the user has the privilege to retain services with prepayment, there will be an automatic renewal. 1 means that there is an automatic renewal. 2 means that there is surely no automatic renewal. If it is not specified, the parameter is 0 by default.
* `time_unit` - (Optional, String, ForceNew) Unit. It is m by default, and only m can be filled in.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



