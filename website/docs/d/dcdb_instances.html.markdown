---
subcategory: "TDSQL for MySQL(dcdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_instances"
sidebar_current: "docs-tencentcloud-datasource-dcdb_instances"
description: |-
  Use this data source to query detailed information of dcdb instances
---

# tencentcloud_dcdb_instances

Use this data source to query detailed information of dcdb instances

## Example Usage

```hcl
data "tencentcloud_dcdb_instances" "instances1" {
  instance_ids        = "your_dcdb_instance1_id"
  search_name         = "instancename"
  search_key          = "search_key"
  project_ids         = [0]
  excluster_type      = 0
  is_filter_excluster = true
  excluster_type      = 0
  is_filter_vpc       = true
  vpc_id              = "your_vpc_id"
  subnet_id           = "your_subnet_id"
}

data "tencentcloud_dcdb_instances" "instances2" {
  instance_ids = ["your_dcdb_instance2_id"]
}

data "tencentcloud_dcdb_instances" "instances3" {
  search_name         = "instancename"
  search_key          = "instances3"
  is_filter_excluster = false
  excluster_type      = 2
}
```

## Argument Reference

The following arguments are supported:

* `excluster_type` - (Optional, Int) cluster excluster type.
* `instance_ids` - (Optional, Set: [`String`]) instance ids.
* `is_filter_excluster` - (Optional, Bool) search according to the cluster excluter type.
* `is_filter_vpc` - (Optional, Bool) search according to the vpc.
* `project_ids` - (Optional, Set: [`Int`]) project ids.
* `result_output_file` - (Optional, String) Used to save results.
* `search_key` - (Optional, String) search key, support fuzzy query.
* `search_name` - (Optional, String) search name, support instancename, vip, all.
* `subnet_id` - (Optional, String) subnet id, valid when IsFilterVpc is true.
* `vpc_id` - (Optional, String) vpc id, valid when IsFilterVpc is true.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - instance list.
  * `app_id` - app id.
  * `auto_renew_flag` - auto renew flag.
  * `create_time` - create time.
  * `db_engine` - db engine.
  * `db_version` - db engine version.
  * `instance_id` - instance id.
  * `instance_name` - instance name.
  * `instance_type` - instance type.
  * `is_audit_supported` - aduit support, 0:support, 1:unsupport.
  * `is_tmp` - tmp instance mark.
  * `isolated_timestamp` - isolated time.
  * `memory` - memory, the unit is GB.
  * `node_count` - node count.
  * `paymode` - pay mode.
  * `period_end_time` - expired time.
  * `project_id` - project id.
  * `region` - region.
  * `resource_tags` - resource tags.
    * `tag_key` - tag key.
    * `tag_value` - tag value.
  * `shard_count` - shard count.
  * `shard_detail` - shard detail.
    * `cpu` - cpu cores.
    * `createtime` - shard create time.
    * `memory` - memory.
    * `node_count` - node count.
    * `shard_id` - shard id.
    * `shard_instance_id` - shard instance id.
    * `shard_serial_id` - shard serial id.
    * `status` - shard status.
    * `storage` - storage.
  * `status_desc` - status description.
  * `status` - status.
  * `storage` - memory, the unit is GB.
  * `subnet_id` - subnet id.
  * `uin` - account uin.
  * `update_time` - update time.
  * `vip` - vip.
  * `vpc_id` - vpc id.
  * `vport` - vport.
  * `wan_domain` - wan domain.
  * `wan_port` - wan port.
  * `wan_status` - wan status, 0:nonactivated, 1:activated, 2:closed, 3:activating.
  * `wan_vip` - wan vip.


