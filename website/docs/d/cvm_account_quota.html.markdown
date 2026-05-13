---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_account_quota"
sidebar_current: "docs-tencentcloud-datasource-cvm_account_quota"
description: |-
  Use this data source to query CVM account quota details.
---

# tencentcloud_cvm_account_quota

Use this data source to query CVM account quota details.

## Example Usage

### Basic query without filters

```hcl
data "tencentcloud_cvm_account_quota" "quota" {}

output "app_id" {
  value = data.tencentcloud_cvm_account_quota.quota.app_id
}
```

### Query by availability zone

```hcl
data "tencentcloud_cvm_account_quota" "quota_zone" {
  zone = ["ap-guangzhou-3", "ap-guangzhou-4"]
}
```

### Query by quota type

```hcl
data "tencentcloud_cvm_account_quota" "quota_type" {
  quota_type = "PostPaidQuotaSet"
}
```

### Query with multiple filters

```hcl
data "tencentcloud_cvm_account_quota" "quota_filtered" {
  zone       = ["ap-guangzhou-3"]
  quota_type = "PostPaidQuotaSet"
}
```

### Query with result output file

```hcl
data "tencentcloud_cvm_account_quota" "quota_output" {
  zone               = ["ap-guangzhou-3"]
  result_output_file = "./quota.json"
}
```

## Argument Reference

The following arguments are supported:

* `quota_type` - (Optional, String) Filter by quota type. Valid values: PostPaidQuotaSet, PrePaidQuotaSet, SpotPaidQuotaSet, ImageQuotaSet, DisasterRecoverGroupQuotaSet.
* `result_output_file` - (Optional, String) Used to save results.
* `zone` - (Optional, Set: [`String`]) Filter by availability zone, such as ap-guangzhou-3.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `account_quota_overview` - Account quota overview.
  * `account_quota` - Account quota details.
    * `disaster_recover_group_quota_set` - Disaster recover group quota list.
      * `current_num` - Current number of groups.
      * `cvm_in_host_group_quota` - Maximum instances in host group.
      * `cvm_in_rack_group_quota` - Maximum instances in rack group.
      * `cvm_in_switch_group_quota` - Maximum instances in switch group.
      * `group_quota` - Group quota.
    * `image_quota_set` - Image quota list.
      * `total_quota` - Total quota.
      * `used_quota` - Used quota.
    * `post_paid_quota_set` - Post-paid quota list.
      * `remaining_quota` - Remaining quota.
      * `total_quota` - Total quota.
      * `used_quota` - Used quota.
      * `zone` - Availability zone.
    * `pre_paid_quota_set` - Pre-paid quota list.
      * `once_quota` - Single purchase quota.
      * `remaining_quota` - Remaining quota.
      * `total_quota` - Total quota.
      * `used_quota` - Used quota.
      * `zone` - Availability zone.
    * `spot_paid_quota_set` - Spot instance quota list.
      * `remaining_quota` - Remaining quota.
      * `total_quota` - Total quota.
      * `used_quota` - Used quota.
      * `zone` - Availability zone.
  * `region` - Region.
* `app_id` - User AppId.


