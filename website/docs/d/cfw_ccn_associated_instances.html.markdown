---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_ccn_associated_instances"
sidebar_current: "docs-tencentcloud-datasource-cfw_ccn_associated_instances"
description: |-
  Use this data source to query detailed information of CFW ccn associated instances
---

# tencentcloud_cfw_ccn_associated_instances

Use this data source to query detailed information of CFW ccn associated instances

## Example Usage

```hcl
data "tencentcloud_cfw_ccn_associated_instances" "example" {
  ccn_id = "ccn-fkb9bo2v"
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String) CCN ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ccn_associated_instances` - Information of instances associated with CCN.
  * `cidr_lst` - List of network segments for the instance.
  * `ins_type` - Instance type.
  * `instance_id` - Instance ID.
  * `instance_name` - Instance name.
  * `instance_region` - Region where the instance belongs.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeCcnAssociatedInstances](https://cloud.tencent.com/document/api/1132/127060) - 查询CCN关联实例信息（不包含云防火墙创建引流实例）

