---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_describe_host_deploy_record"
sidebar_current: "docs-tencentcloud-datasource-ssl_describe_host_deploy_record"
description: |-
  Use this data source to query detailed information of ssl describe_host_deploy_record
---

# tencentcloud_ssl_describe_host_deploy_record

Use this data source to query detailed information of ssl describe_host_deploy_record

## Example Usage

```hcl
data "tencentcloud_ssl_describe_host_deploy_record" "describe_host_deploy_record" {
  certificate_id = "8u8DII0l"
  resource_type  = "ddos"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required, String) Certificate ID to be deployed.
* `resource_type` - (Optional, String) Resource Type.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `deploy_record_list` - Certificate deployment record listNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `cert_id` - Deployment certificate ID.
  * `create_time` - Deployment time.
  * `id` - Deployment record ID.
  * `region` - Deployment.
  * `resource_type` - Deploy resource type.
  * `status` - Deployment state.
  * `update_time` - Recent update time.


