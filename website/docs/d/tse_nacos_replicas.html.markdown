---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_nacos_replicas"
sidebar_current: "docs-tencentcloud-datasource-tse_nacos_replicas"
description: |-
  Use this data source to query detailed information of tse nacos_replicas
---

# tencentcloud_tse_nacos_replicas

Use this data source to query detailed information of tse nacos_replicas

## Example Usage

```hcl
data "tencentcloud_tse_nacos_replicas" "nacos_replicas" {
  instance_id = "ins-8078da86"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) engine instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `replicas` - Engine instance replica information.
  * `name` - name.
  * `role` - role.
  * `status` - status.
  * `subnet_id` - Subnet IDNote: This field may return null, indicating that a valid value is not available.
  * `vpc_id` - VPC IDNote: This field may return null, indicating that a valid value is not available.
  * `zone_id` - Available area IDNote: This field may return null, indicating that a valid value is not available.
  * `zone` - Available area NameNote: This field may return null, indicating that a valid value is not available.


