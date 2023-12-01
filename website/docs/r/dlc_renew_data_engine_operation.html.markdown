---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_renew_data_engine_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_renew_data_engine_operation"
description: |-
  Provides a resource to create a dlc renew_data_engine
---

# tencentcloud_dlc_renew_data_engine_operation

Provides a resource to create a dlc renew_data_engine

## Example Usage

```hcl
resource "tencentcloud_dlc_renew_data_engine_operation" "renew_data_engine" {
  data_engine_name = "testEngine"
  time_span        = 3600
  pay_mode         = 1
  time_unit        = "m"
  renew_flag       = 1
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_name` - (Required, String, ForceNew) Data engine name.
* `time_span` - (Required, Int, ForceNew) Engine TimeSpan, prePay: minimum of 1, representing one month of purchasing resources, with a maximum of 120, default 3600, postPay: fixed fee of 3600.
* `pay_mode` - (Optional, Int, ForceNew) Engine pay mode type, only support 0: postPay, 1: prePay(default).
* `renew_flag` - (Optional, Int, ForceNew) Automatic renewal flag, 0, initial state, automatic renewal is not performed by default. if the user has prepaid non-stop service privileges, automatic renewal will occur. 1: Automatic renewal. 2: make it clear that there will be no automatic renewal. if this parameter is not passed, the default value is 0.
* `time_unit` - (Optional, String, ForceNew) Engine TimeUnit, prePay: use m(default), postPay: use h.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dlc renew_data_engine can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_renew_data_engine_operation.renew_data_engine renew_data_engine_id
```

