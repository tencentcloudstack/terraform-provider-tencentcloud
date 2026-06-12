---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cluster_read_only_instance_group_acces_operation"
sidebar_current: "docs-tencentcloud-resource-cynosdb_cluster_read_only_instance_group_acces_operation"
description: |-
  Provides a resource to open CynosDB (TDSQL-C) cluster read-only instance group access
---

# tencentcloud_cynosdb_cluster_read_only_instance_group_acces_operation

Provides a resource to open CynosDB (TDSQL-C) cluster read-only instance group access

## Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_read_only_instance_group_acces_operation" "example" {
  cluster_id         = "cynosdbmysql-8rn1byp7"
  port               = "3306"
  security_group_ids = ["sg-4rd5741x"]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `port` - (Optional, String, ForceNew) Port.
* `security_group_ids` - (Optional, List: [`String`], ForceNew) Security group IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `5m`) Used when creating the resource.

