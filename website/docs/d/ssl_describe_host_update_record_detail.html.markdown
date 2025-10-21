---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_describe_host_update_record_detail"
sidebar_current: "docs-tencentcloud-datasource-ssl_describe_host_update_record_detail"
description: |-
  Use this data source to query detailed information of ssl describe_host_update_record_detail
---

# tencentcloud_ssl_describe_host_update_record_detail

Use this data source to query detailed information of ssl describe_host_update_record_detail

## Example Usage

```hcl
data "tencentcloud_ssl_describe_host_update_record_detail" "describe_host_update_record_detail" {
  deploy_record_id = "35364"
}
```

## Argument Reference

The following arguments are supported:

* `deploy_record_id` - (Required, String) One -click update record ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `failed_total_count` - Total number of failuresNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `record_detail_list` - Certificate deployment record listNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `list` - List of deployment resources details.
    * `bucket` - BUCKET name (COS dedicated)Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `cert_id` - New certificate ID.
    * `create_time` - Deployment time.
    * `domains` - List of deployment domainNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `env_id` - Environment IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `error_msg` - Deployment error messageNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `id` - Detailed record ID.
    * `instance_id` - Deployment instance IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `instance_name` - Deployment example nameNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `listener_id` - Deploy listener ID (CLB for CLB)Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `listener_name` - Deploy listener name (CLB for CLB)Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `namespace` - Naming Space (TKE)Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `old_cert_id` - Old certificate ID.
    * `port` - portNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `protocol` - protocolNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `region` - DeploymentNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `resource_type` - Deploy resource type.
    * `secret_name` - Secret Name (TKE for TKE)Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `sni_switch` - Whether to turn on SNI (CLB dedicated)Note: This field may return NULL, indicating that the valid value cannot be obtained.
    * `status` - Deployment state.
    * `t_c_b_type` - TCB deployment typeNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `update_time` - Last update time.
  * `resource_type` - Deploy resource type.
  * `total_count` - The total number of deployment resources.
* `running_total_count` - Total number of deploymentNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `success_total_count` - Total successNote: This field may return NULL, indicating that the valid value cannot be obtained.


