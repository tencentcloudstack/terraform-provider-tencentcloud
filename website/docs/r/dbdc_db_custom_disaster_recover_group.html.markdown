---
subcategory: "Database Dedicated Cluster(DBDC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbdc_db_custom_disaster_recover_group"
sidebar_current: "docs-tencentcloud-resource-dbdc_db_custom_disaster_recover_group"
description: |-
  Provides a resource to create a DBDC db custom disaster recover group.
---

# tencentcloud_dbdc_db_custom_disaster_recover_group

Provides a resource to create a DBDC db custom disaster recover group.

## Example Usage

```hcl
resource "tencentcloud_dbdc_db_custom_disaster_recover_group" "example" {
  name     = "tf-example1"
  type     = "HOST"
  strategy = "SPREAD"
  affinity = 1

  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Placement group name. Up to 60 characters, only Chinese and English are allowed.
* `affinity` - (Optional, Int) Affinity of the placement group. Instances in the group will be distributed according to this affinity. Valid values: `[1, 10]`. Default value: `1`.
* `strategy` - (Optional, String, ForceNew) Placement group strategy. Valid values: `SPREAD` (spread placement group). Default value: `SPREAD`.
* `tags` - (Optional, Map) Tags.
* `type` - (Optional, String, ForceNew) Placement group type. Valid values: `HOST` (physical machine). Default value: `HOST`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - Creation time.
* `current_num` - Current number of nodes in the placement group.
* `disaster_recover_group_id` - Placement group ID.
* `node_ids` - List of DB Custom node IDs in the placement group.
* `node_quota_total` - Maximum number of nodes that the placement group can hold.
* `status` - Placement group status. Valid values: `Creating`, `Available`, `CreateFailed`, `Deleting`, `Modifying`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `10m`) Used when creating the resource.
* `update` - (Defaults to `10m`) Used when updating the resource.
* `delete` - (Defaults to `10m`) Used when deleting the resource.

## Import

DBDC db custom disaster recover group can be imported using the id, e.g.

```
terraform import tencentcloud_dbdc_db_custom_disaster_recover_group.example dbps-sdcm1hfl
```

