---
subcategory: "SSL Certificates"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_describe_host_deploy_record_detail"
sidebar_current: "docs-tencentcloud-datasource-ssl_describe_host_deploy_record_detail"
description: |-
  Use this data source to query detailed information of ssl describe_host_deploy_record_detail
---

# tencentcloud_ssl_describe_host_deploy_record_detail

Use this data source to query detailed information of ssl describe_host_deploy_record_detail

## Example Usage

```hcl
data "tencentcloud_ssl_describe_host_deploy_record_detail" "describe_host_deploy_record_detail" {
  deploy_record_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `deploy_record_id` - (Required, String) Deployment record ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `deploy_record_detail_list` - Certificate deployment record listNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `bucket` - COS storage barrel nameNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `cert_id` - Deployment certificate ID.
  * `create_time` - Deployment record details Create time.
  * `domains` - List of deployment domain.
  * `env_id` - TCB environment IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `error_msg` - Deployment error messageNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `id` - Deployment record details ID.
  * `instance_id` - Deployment instance ID.
  * `instance_name` - Deployment example name.
  * `listener_id` - Deployment monitor IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `listener_name` - Delicate monitor name.
  * `namespace` - Named space nameNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `old_cert_id` - Original binding certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `port` - portNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `protocol` - Deployment monitoring protocolNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `region` - Deployed TCB regionNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `secret_name` - Secret nameNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `sni_switch` - Whether to turn on SNI.
  * `status` - Deployment state.
  * `tcb_type` - Deployed TCB typeNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `update_time` - Deployment record details last update time.
* `failed_total_count` - Total number of failuresNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `running_total_count` - Total number of deploymentNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `success_total_count` - Total successNote: This field may return NULL, indicating that the valid value cannot be obtained.


