---
subcategory: "MapReduce(EMR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_emr_cvm_quota"
sidebar_current: "docs-tencentcloud-datasource-emr_cvm_quota"
description: |-
  Use this data source to query detailed information of emr cvm_quota
---

# tencentcloud_emr_cvm_quota

Use this data source to query detailed information of emr cvm_quota

## Example Usage

```hcl
data "tencentcloud_emr_cvm_quota" "cvm_quota" {
  cluster_id = "emr-0ze36vnp"
  zone_id    = 100003
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) EMR cluster ID.
* `result_output_file` - (Optional, String) Used to save results.
* `zone_id` - (Optional, Int) Zone ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `eks_quota_set` - Eks quota Note: This field may return null, indicating that a valid value cannot be obtained.
  * `cpu` - Cpu cores.
  * `memory` - Memory quantity (unit: GB).
  * `node_type` - The specifications of the marketable resource are as follows: `TASK`, `CORE`, `MASTER`, `ROUTER`.
  * `number` - Specifies the maximum number of resources that can be applied for.
* `post_paid_quota_set` - Postpaid quota list Note: This field may return null, indicating that no valid value can be obtained.
  * `remaining_quota` - Residual quota Note: This field may return null, indicating that a valid value cannot be obtained.
  * `total_quota` - Total quota Note: This field may return null, indicating that a valid value cannot be obtained.
  * `used_quota` - Used quota Note: This field may return null, indicating that a valid value cannot be obtained.
  * `zone` - Available area Note: This field may return null, indicating that a valid value cannot be obtained.
* `spot_paid_quota_set` - Biding instance quota list Note: This field may return null, indicating that a valid value cannot be obtained.
  * `remaining_quota` - Residual quota Note: This field may return null, indicating that a valid value cannot be obtained.
  * `total_quota` - Total quota Note: This field may return null, indicating that a valid value cannot be obtained.
  * `used_quota` - Used quota Note: This field may return null, indicating that a valid value cannot be obtained.
  * `zone` - Available area Note: This field may return null, indicating that a valid value cannot be obtained.


