---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_resource_package_list"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_resource_package_list"
description: |-
  Use this data source to query detailed information of cynosdb resource_package_list
---

# tencentcloud_cynosdb_resource_package_list

Use this data source to query detailed information of cynosdb resource_package_list

## Example Usage

```hcl
data "tencentcloud_cynosdb_resource_package_list" "resource_package_list" {
  package_id      = ["package-hy4d2ppl"]
  package_name    = ["keep-package-disk"]
  package_type    = ["DISK"]
  package_region  = ["china"]
  status          = ["using"]
  order_by        = ["startTime"]
  order_direction = "DESC"
}
```

## Argument Reference

The following arguments are supported:

* `order_by` - (Optional, Set: [`String`]) Sorting conditions supported: startTime - effective time, expireTime - expiration time, packageUsedSpec - usage capacity, and packageTotalSpec - total storage capacity. Arrange in array order;.
* `order_direction` - (Optional, String) Sort by, DESC Descending, ASC Ascending.
* `package_id` - (Optional, Set: [`String`]) Resource Package Unique ID.
* `package_name` - (Optional, Set: [`String`]) Resource Package Name.
* `package_region` - (Optional, Set: [`String`]) Resource package usage region China - common in mainland China, overseas - common in Hong Kong, Macao, Taiwan, and overseas.
* `package_type` - (Optional, Set: [`String`]) Resource package type CCU - Compute resource package, DISK - Storage resource package.
* `result_output_file` - (Optional, String) Used to save results.
* `status` - (Optional, Set: [`String`]) Resource package status creating - creating; Using - In use; Expired - has expired; Normal_ Finish - used up; Apply_ Refund - Applying for a refund; Refund - The fee has been refunded.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `resource_package_list` - Resource package details note: This field may return null, indicating that a valid value cannot be obtained.
  * `app_id` - AppID note: This field may return null, indicating that a valid value cannot be obtained.
  * `bind_instance_infos` - Note for binding instance information: This field may return null, indicating that a valid value cannot be obtained.
    * `instance_id` - Instance ID.
    * `instance_region` - Region of instance.
    * `instance_type` - Instance type.
  * `expire_time` - Expiration time: August 1st, 2022 00:00:00 Attention: This field may return null, indicating that a valid value cannot be obtained.
  * `has_quota` - Resource package usage note: This field may return null, indicating that a valid value cannot be obtained.
  * `package_id` - Resource Package Unique ID Note: This field may return null, indicating that a valid value cannot be obtained.
  * `package_name` - Resource package name note: This field may return null, indicating that a valid value cannot be obtained.
  * `package_region` - The resource package is used in China, which is commonly used in mainland China, and in overseas, which is commonly used in Hong Kong, Macao, Taiwan, and overseas. Note: This field may return null, indicating that a valid value cannot be obtained.
  * `package_total_spec` - Attention to the total amount of resource packages: This field may return null, indicating that a valid value cannot be obtained.
  * `package_type` - Resource package type CCU - Compute resource package, DISK - Store resource package Note: This field may return null, indicating that a valid value cannot be obtained.
  * `package_used_spec` - Resource package usage note: This field may return null, indicating that a valid value cannot be obtained.
  * `start_time` - Effective time: July 1st, 2022 00:00:00 Attention: This field may return null, indicating that a valid value cannot be obtained.
  * `status` - Resource package status creating - creating; Using - In use; Expired - has expired; Normal_ Finish - used up; Apply_ Refund - Applying for a refund; Refund - The fee has been refunded. Note: This field may return null, indicating that a valid value cannot be obtained.


