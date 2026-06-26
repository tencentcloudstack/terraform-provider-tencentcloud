---
subcategory: "Database Dedicated Cluster(DBDC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbdc_node_to_db_custom_cluster_attachment"
sidebar_current: "docs-tencentcloud-resource-dbdc_node_to_db_custom_cluster_attachment"
description: |-
  Provides a resource to create a DBDC node to db custom cluster attachment.
---

# tencentcloud_dbdc_node_to_db_custom_cluster_attachment

Provides a resource to create a DBDC node to db custom cluster attachment.

~> **NOTE:** Both create and delete operations are asynchronous. The resource waits for the underlying task to reach the `Succeeded` status before returning.

## Example Usage

```hcl
resource "tencentcloud_dbdc_node_to_db_custom_cluster_attachment" "example" {
  cluster_id = "dbcc-xxxxxxxx"
  node_id    = "dbcn-xxxxxxxx"
  image_id   = "img-xxxxxxxx"

  login_settings {
    password = "Passw0rd@2024"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) DB Custom cluster ID.
* `node_id` - (Required, String, ForceNew) DB Custom node ID to add to the cluster.
* `image_id` - (Optional, String, ForceNew) OS image ID to reset the node to after it is added to the cluster.
* `login_settings` - (Optional, List, ForceNew) Instance login settings. You can set the login method to password, key, or keep the original image login settings. Only one method can be set.

The `login_settings` object supports the following:

* `keep_image_login` - (Optional, String, ForceNew) Whether to keep the original login settings of the image. Valid values: `true`, `false`. Cannot be specified together with Password or KeyIds.
* `key_ids` - (Optional, List, ForceNew) Key pair ID list. Only a single ID is supported currently. Password and key cannot be specified at the same time.
* `password` - (Optional, String, ForceNew) Instance login password. Password complexity limits vary by operating system type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `lan_ip` - Intranet IP address of the node.
* `node_name` - Node name.
* `node_type` - Node spec.
* `ssh_endpoint` - SSH endpoint to access the node, in the format `IP:Port`.
* `status` - Instance status of the node in the cluster.
* `zone` - Availability zone that the node belongs to.


## Timeouts

This resource provides the following [Timeouts](https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts) configuration options:

- `create` - (Default `30m`)
- `delete` - (Default `30m`)

## Import

DBDC node to db custom cluster attachment can be imported using the id, e.g.

```
terraform import tencentcloud_dbdc_node_to_db_custom_cluster_attachment.example dbcc-xxxxxxxx#dbcn-xxxxxxxx
```

