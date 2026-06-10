---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_db_instance_node_property"
sidebar_current: "docs-tencentcloud-datasource-mongodb_db_instance_node_property"
description: |-
  Use this data source to query detailed information of MongoDB (mongodb) DB instance node property
---

# tencentcloud_mongodb_db_instance_node_property

Use this data source to query detailed information of MongoDB (mongodb) DB instance node property

## Example Usage

```hcl
data "tencentcloud_mongodb_db_instance_node_property" "example" {
  instance_id = "cmgo-9d0p6umb"
}
```

### Example Usage with filters

```hcl
data "tencentcloud_mongodb_db_instance_node_property" "example" {
  instance_id = "cmgo-9d0p6umb"
  roles       = ["PRIMARY", "SECONDARY"]
  only_hidden = false
  priority    = 1
  votes       = 1
  tags {
    tag_key   = "env"
    tag_value = "prod"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `node_ids` - (Optional, List: [`String`]) Node ID list.
* `only_hidden` - (Optional, Bool) Whether to query only Hidden nodes. Default is false.
* `priority` - (Optional, Int) Node priority. Value range: [0, 100].
* `result_output_file` - (Optional, String) Used to save results.
* `roles` - (Optional, List: [`String`]) Node role list. Valid values: PRIMARY, SECONDARY, READONLY, ARBITER.
* `tags` - (Optional, List) Node tags for filtering.
* `votes` - (Optional, Int) Node votes. 1: has votes; 0: no votes.

The `tags` object supports the following:

* `tag_key` - (Optional, String) Node tag key.
* `tag_value` - (Optional, String) Node tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `mongos` - Mongos node property list.
  * `address` - Node access address.
  * `hidden` - Whether the node is a Hidden node.
  * `node_name` - Node name.
  * `priority` - Node priority. Value range: [0, 100].
  * `replicate_set_id` - Replica set ID.
  * `role` - Node role. Valid values: PRIMARY, SECONDARY, READONLY, ARBITER.
  * `slave_delay` - Primary-secondary sync delay in seconds.
  * `status` - Node status. Valid values: NORMAL, STARTUP, STARTUP2, RECOVERING, DOWN, UNKNOWN, ROLLBACK, REMOVED.
  * `tags` - Node tags.
    * `tag_key` - Node tag key.
    * `tag_value` - Node tag value.
  * `votes` - Node votes. 1: has votes; 0: no votes.
  * `wan_service_address` - Node public network access address (IP or domain name).
  * `zone` - The availability zone where the node is located.
* `replicate_sets` - Replica set node info list.
  * `nodes` - Node property list in the replica set.
    * `address` - Node access address.
    * `hidden` - Whether the node is a Hidden node.
    * `node_name` - Node name.
    * `priority` - Node priority. Value range: [0, 100].
    * `replicate_set_id` - Replica set ID.
    * `role` - Node role. Valid values: PRIMARY, SECONDARY, READONLY, ARBITER.
    * `slave_delay` - Primary-secondary sync delay in seconds.
    * `status` - Node status. Valid values: NORMAL, STARTUP, STARTUP2, RECOVERING, DOWN, UNKNOWN, ROLLBACK, REMOVED.
    * `tags` - Node tags.
      * `tag_key` - Node tag key.
      * `tag_value` - Node tag value.
    * `votes` - Node votes. 1: has votes; 0: no votes.
    * `wan_service_address` - Node public network access address (IP or domain name).
    * `zone` - The availability zone where the node is located.


