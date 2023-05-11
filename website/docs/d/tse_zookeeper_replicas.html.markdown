---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_zookeeper_replicas"
sidebar_current: "docs-tencentcloud-datasource-tse_zookeeper_replicas"
description: |-
  Use this data source to query detailed information of tse zookeeper_replicas
---

# tencentcloud_tse_zookeeper_replicas

Use this data source to query detailed information of tse zookeeper_replicas

## Example Usage

```hcl
data "tencentcloud_tse_zookeeper_replicas" "zookeeper_replicas" {
  instance_id = "ins-7eb7eea7"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) engine instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `replicas` - Engine instance replica information.
  * `alias_name` - aliasNote: This field may return null, indicating that a valid value is not available.
  * `name` - name.
  * `role` - role.
  * `status` - status.
  * `subnet_id` - Subnet IDNote: This field may return null, indicating that a valid value is not available.
  * `vpc_id` - VPC IDNote: This field may return null, indicating that a valid value is not available.
  * `zone_id` - Available area IDNote: This field may return null, indicating that a valid value is not available.
  * `zone` - Available area IDNote: This field may return null, indicating that a valid value is not available.


