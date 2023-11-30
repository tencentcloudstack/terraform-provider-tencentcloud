---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_describe_host_update_record"
sidebar_current: "docs-tencentcloud-datasource-ssl_describe_host_update_record"
description: |-
  Use this data source to query detailed information of ssl describe_host_update_record
---

# tencentcloud_ssl_describe_host_update_record

Use this data source to query detailed information of ssl describe_host_update_record

## Example Usage

```hcl
data "tencentcloud_ssl_describe_host_update_record" "describe_host_update_record" {
  old_certificate_id = "8u8DII0l"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Optional, String) New certificate ID.
* `old_certificate_id` - (Optional, String) Original certificate ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `deploy_record_list` - Certificate deployment record listNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `cert_id` - New certificate ID.
  * `create_time` - Deployment time.
  * `id` - Record ID.
  * `old_cert_id` - Original certificate ID.
  * `regions` - List of regional deploymentNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `resource_types` - List of resource types.
  * `status` - Deployment state.
  * `update_time` - Last update time.


