---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_plan"
sidebar_current: "docs-tencentcloud-resource-teo_plan"
description: |-
  Provides a resource to create a TEO plan
---

# tencentcloud_teo_plan

Provides a resource to create a TEO plan

## Example Usage

```hcl
resource "tencentcloud_teo_plan" "example" {
  plan_type = "personal"

  prepaid_plan_param {
    period     = 1
    renew_flag = "off"
  }
}
```

## Argument Reference

The following arguments are supported:

* `plan_type` - (Required, String) The subscription package type, the possible values are: `personal`: personal package, prepaid package; `basic`: basic package, prepaid package; `standard`: standard package, prepaid package; `enterprise`: enterprise package, postpaid package.
* `prepaid_plan_param` - (Optional, List) Subscription prepaid package parameters. When PlanType is personal, basic, or standard, this parameter is optional and is used to enter the subscription duration of the package and whether to enable automatic renewal. If this parameter is not filled in, the default subscription duration is 1 month and automatic renewal is not enabled.

The `prepaid_plan_param` object supports the following:

* `period` - (Optional, Int) The subscription period of the prepaid package, in months, with possible values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. If not filled in, the default value 1 is used.
* `renew_flag` - (Optional, String) The automatic renewal flag of the prepaid package, the values are: `on`: turn on automatic renewal; `off`: do not turn on automatic renewal. If not filled in, the default value off is used. When automatic renewal occurs, the default renewal period is 1 month.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `area` - Service area, possible values are: <li>mainland: Mainland China; </li><li>overseas: Worldwide (excluding Mainland China); </li><li>global: Worldwide (including Mainland China). </li>.
* `enabled_time` - The time when the package takes effect.
* `expired_time` - The expiration date of the package.
* `pay_mode` - Payment type, possible values: <li>0: post-payment; </li><li>1: pre-payment. </li>.
* `plan_id` - Plan ID.
* `status` - Package status, the values are: <li>normal: normal status; </li><li>expiring-soon: about to expire; </li><li>expired: expired; </li><li>isolated: isolated; </li><li>overdue-isolated: overdue isolated. </li>.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreatePlan](https://www.tencentcloud.com/document/api/1145/60514)
* Read: [DescribePlans](https://www.tencentcloud.com/document/api/1145/72310)
* Update:
  - [ModifyPlan](https://www.tencentcloud.com/document/api/1145/60511)
  - [RenewPlan](https://www.tencentcloud.com/document/api/1145/60510)
  - [UpgradePlan](https://www.tencentcloud.com/document/api/1145/60509)
* Delete: [DestroyPlan](https://www.tencentcloud.com/document/api/1145/60513)

## Import

TEO plan can be imported using the id, e.g.

```
terraform import tencentcloud_teo_plan.example edgeone-2o1xvpmq7nn
```

