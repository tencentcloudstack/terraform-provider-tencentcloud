---
subcategory: "Database Dedicated Cluster(DBDC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbdc_db_custom_cluster"
sidebar_current: "docs-tencentcloud-resource-dbdc_db_custom_cluster"
description: |-
  Provides a resource to create a DBDC db custom cluster.
---

# tencentcloud_dbdc_db_custom_cluster

Provides a resource to create a DBDC db custom cluster.

~> **NOTE:** Both create and destroy operations are asynchronous. The resource waits for the underlying task to reach the `Succeeded` status before returning.

## Example Usage

```hcl
resource "tencentcloud_dbdc_db_custom_cluster" "example" {
  cluster_name        = "tf-example"
  cluster_description = "cluster description."

  container_network {
    vpc_id     = "vpc-xxxxxxxx"
    subnet_ids = ["subnet-xxxxxxxx"]
  }

  api_server_network {
    vpc_id    = "vpc-xxxxxxxx"
    subnet_id = "subnet-xxxxxxxx"
  }

  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `api_server_network` - (Required, List, ForceNew) Network information of the cluster API Server. Must be a network owned by this account, and can be the same as the container network.
* `cluster_name` - (Required, String, ForceNew) Cluster name. Up to 128 characters, only Chinese, English and underscore are allowed.
* `container_network` - (Required, List, ForceNew) Container network. All pods in this cluster are connected to this network.
* `cluster_description` - (Optional, String, ForceNew) Cluster description.
* `tags` - (Optional, Map) Cluster tags.

The `api_server_network` object supports the following:

* `subnet_id` - (Required, String, ForceNew) Subnet ID of the API Server network.
* `vpc_id` - (Required, String, ForceNew) VPC ID of the API Server network.

The `container_network` object supports the following:

* `subnet_ids` - (Required, List, ForceNew) Subnet ID list of the container network.
* `vpc_id` - (Required, String, ForceNew) VPC ID of the container network.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cluster_level` - Cluster level.
* `cluster_node_num` - Number of nodes in the cluster.
* `cluster_status` - DB Custom cluster status. Valid values: `Creating`, `Running`, `Destroying`.
* `cluster_version` - Cluster version.
* `created_time` - Creation time.
* `region` - Region that the cluster belongs to.


## Timeouts

This resource provides the following [Timeouts](https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts) configuration options:

- `create` - (Default `15m`)
- `update` - (Default `15m`)
- `delete` - (Default `15m`)

## Import

DBDC db custom cluster can be imported using the id, e.g.

```
terraform import tencentcloud_dbdc_db_custom_cluster.example dbcc-xxxxxxxx
```

