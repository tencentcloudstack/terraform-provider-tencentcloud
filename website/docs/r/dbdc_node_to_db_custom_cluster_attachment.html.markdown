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

## Example Usage

```hcl
# create cluster
resource "tencentcloud_dbdc_db_custom_cluster" "example" {
  cluster_name        = "tf-example"
  cluster_description = "cluster description."

  container_network {
    vpc_id     = "vpc-py7mlxqm"
    subnet_ids = ["subnet-qd4upp83", "subnet-g7vmz57f", "subnet-hqbm5bwx"]
  }

  api_server_network {
    vpc_id    = "vpc-b4zgfr3a"
    subnet_id = "subnet-cp3juq8r"
  }

  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_dbdc_db_custom_node" "example" {
  zone       = "ap-shanghai-5"
  image_id   = "img-rm13akp3"
  vpc_id     = "vpc-py7mlxqm"
  subnet_id  = "subnet-qd4upp83"
  node_type  = "DB.AT5.8XLARGE128"
  period     = 1
  auto_renew = 2
  node_name  = "tf-example"

  login_settings {
    password = "Password@2026"
  }

  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_dbdc_node_to_db_custom_cluster_attachment" "example" {
  cluster_id = tencentcloud_dbdc_db_custom_cluster.example.id
  node_id    = tencentcloud_dbdc_db_custom_node.example.id
  image_id   = "img-rm13akp3"

  login_settings {
    password = "Passw0rd@2026"
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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `1h0m`) Used when creating the resource.
* `delete` - (Defaults to `1h0m`) Used when deleting the resource.

## Import

DBDC node to db custom cluster attachment can be imported using the clusterId#nodeId, e.g.

```
terraform import tencentcloud_dbdc_node_to_db_custom_cluster_attachment.example dbcc-7uh7ludb#dbcn-ttiyh58n
```

