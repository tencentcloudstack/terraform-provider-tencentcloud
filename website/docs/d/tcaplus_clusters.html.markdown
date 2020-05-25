---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_clusters"
sidebar_current: "docs-tencentcloud-datasource-tcaplus_clusters"
description: |-
  Use this data source to query tcaplus clusters
---

# tencentcloud_tcaplus_clusters

Use this data source to query tcaplus clusters

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

* `cluster_id` - (Optional) Id of the tcaplus cluster to be query.
* `cluster_name` - (Optional) Name of the tcaplus cluster to be query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of tcaplus cluster. Each element contains the following attributes.
  * `api_access_id` - Access id of the tcaplus cluster.For TcaplusDB SDK connect.
  * `api_access_ip` - Access ip of the tcaplus cluster.For TcaplusDB SDK connect.
  * `api_access_port` - Access port of the tcaplus cluster.For TcaplusDB SDK connect.
  * `cluster_id` - Id of the tcaplus cluster.
  * `cluster_name` - Name of the tcaplus cluster.
  * `create_time` - Create time of the tcaplus cluster.
  * `idl_type` - Idl type of the tcaplus cluster.
  * `network_type` - Network type of the tcaplus cluster.
  * `old_password_expire_time` - This field will display the old password expiration time,if password_status is `unmodifiable` means the old password has not yet expired, otherwise `-`.
  * `password_status` - Password status of the tcaplus cluster.`unmodifiable` means:can not change password now,`modifiable` means:can change password now.
  * `password` - Password of the tcaplus cluster.
  * `subnet_id` - Subnet id of the tcaplus cluster.
  * `vpc_id` - VPC id of the tcaplus cluster.


