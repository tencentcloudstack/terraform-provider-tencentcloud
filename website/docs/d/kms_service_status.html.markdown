---
subcategory: "Key Management Service(KMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_service_status"
sidebar_current: "docs-tencentcloud-datasource-kms_service_status"
description: |-
  Use this data source to query detailed information of KMS service_status
---

# tencentcloud_kms_service_status

Use this data source to query detailed information of KMS service_status

## Example Usage

```hcl
data "tencentcloud_kms_service_status" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cmk_limit` - Return KMS user key specification quantity.
* `cmk_user_count` - Return the number of KMS user key usage.
* `exclusive_hsm_enabled` - Whether to activate Exclusive KMS
Note: This field may return `null`, indicating that no valid value can be obtained.
* `exclusive_hsm_list` - Return to Exclusive Cluster Group.
* `exclusive_vsm_enabled` - Whether to activate Managed KMS
Note: This field may return `null`, indicating that no valid value can be obtained.
* `invalid_type` - Service unavailability type. 0: not purchased; 1: normal; 2: suspended due to arrears; 3: resource released.
* `pro_expire_time` - Expiration time of the KMS Ultimate edition. It's represented in a Unix Epoch timestamp.
Note: This field may return null, indicating that no valid values can be obtained.
* `pro_renew_flag` - Whether to automatically renew Ultimate Edition. 0: no, 1: yes
Note: this field may return null, indicating that no valid values can be obtained.
* `pro_resource_id` - Unique ID of the Ultimate Edition purchase record. If the Ultimate Edition is not activated, the returned value will be null.
Note: this field may return null, indicating that no valid values can be obtained.
* `service_enabled` - Whether the KMS service has been activated. true: activated.
* `subscription_info` - KMS subscription information.
Note: This field may return null, indicating that no valid values can be obtained.
* `user_level` - 0: Basic Edition, 1: Ultimate Edition.


