---
subcategory: "TcaplusDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_clusters"
sidebar_current: "docs-tencentcloud-datasource-tcaplus_clusters"
description: |-
  Use this data source to query TcaplusDB clusters.
---

# tencentcloud_tcaplus_clusters

Use this data source to query TcaplusDB clusters.

## Example Usage

```hcl
data "tencentcloud_tcaplus_clusters" "name" {
  cluster_name = "cluster"
}
data "tencentcloud_tcaplus_clusters" "id" {
  cluster_id = tencentcloud_tcaplus_cluster.test.id
}
data "tencentcloud_tcaplus_clusters" "idname" {
  cluster_id   = tencentcloud_tcaplus_cluster.test.id
  cluster_name = "cluster"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional, String) ID of the TcaplusDB cluster to be query.
* `cluster_name` - (Optional, String) Name of the TcaplusDB cluster to be query.
* `result_output_file` - (Optional, String) File for saving results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of TcaplusDB cluster. Each element contains the following attributes.
  * `api_access_id` - Access id of the TcaplusDB cluster.For TcaplusDB SDK connect.
  * `api_access_ip` - Access ip of the TcaplusDB cluster.For TcaplusDB SDK connect.
  * `api_access_port` - Access port of the TcaplusDB cluster.For TcaplusDB SDK connect.
  * `cluster_id` - ID of the TcaplusDB cluster.
  * `cluster_name` - Name of the TcaplusDB cluster.
  * `create_time` - Create time of the TcaplusDB cluster.
  * `idl_type` - IDL type of the TcaplusDB cluster.
  * `network_type` - Network type of the TcaplusDB cluster.
  * `old_password_expire_time` - Expiration time of the old password. If `password_status` is `unmodifiable`, it means the old password has not yet expired.
  * `password_status` - Password status of the TcaplusDB cluster. Valid values: `unmodifiable`, `modifiable`. `unmodifiable` means the password can not be changed in this moment; `modifiable` means the password can be changed in this moment.
  * `password` - Access password of the TcaplusDB cluster.
  * `subnet_id` - Subnet ID of the TcaplusDB cluster.
  * `vpc_id` - VPC ID of the TcaplusDB cluster.


