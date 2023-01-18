---
subcategory: "Cloud File Storage(CFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfs_file_system_clients"
sidebar_current: "docs-tencentcloud-datasource-cfs_file_system_clients"
description: |-
  Use this data source to query detailed information of cfs file_system_clients
---

# tencentcloud_cfs_file_system_clients

Use this data source to query detailed information of cfs file_system_clients

## Example Usage

```hcl
data "tencentcloud_cfs_file_system_clients" "file_system_clients" {
  file_system_id = "cfs-iobiaxtj"
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, String) File system ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `client_list` - Client list.
  * `cfs_vip` - IP address of the file system.
  * `client_ip` - Client IP.
  * `mount_directory` - Path in which the file system is mounted to the client.
  * `vpc_id` - File system VPCID.
  * `zone_name` - AZ name.
  * `zone` - Name of the availability zone, e.g. ap-beijing-1. For more information, see regions and availability zones in the Overview document.


