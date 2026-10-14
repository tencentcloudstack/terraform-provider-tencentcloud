---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_wan_ip_config"
sidebar_current: "docs-tencentcloud-resource-sqlserver_wan_ip_config"
description: |-
  Provides a resource to create a sqlserver wan ip config
---

# tencentcloud_sqlserver_wan_ip_config

Provides a resource to create a sqlserver wan ip config

## Example Usage

### Open/Close wan ip for SQL instance

```hcl
# open
resource "tencentcloud_sqlserver_wan_ip_config" "example" {
  instance_id   = "mssql-gy1lc54f"
  enable_wan_ip = true
}

# close
resource "tencentcloud_sqlserver_wan_ip_config" "example" {
  instance_id   = "mssql-gy1lc54f"
  enable_wan_ip = false
}
```

### Open/Close wan ip for SQL read only group

```hcl
# open
resource "tencentcloud_sqlserver_wan_ip_config" "example" {
  instance_id   = "mssql-gy1lc54f"
  ro_group_id   = "mssqlrg-hyxotm31"
  enable_wan_ip = true
}

# close
resource "tencentcloud_sqlserver_wan_ip_config" "example" {
  instance_id   = "mssql-gy1lc54f"
  ro_group_id   = "mssqlrg-hyxotm31"
  enable_wan_ip = false
}
```

## Argument Reference

The following arguments are supported:

* `enable_wan_ip` - (Required, Bool) Whether to open wan ip, true: enable; false: disable.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `ro_group_id` - (Optional, String, ForceNew) Read only group ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `dns_pod_domain` - Internet address domain name.
* `ro_group` - Read only group.
  * `dns_pod_domain` - Internet address domain name.
  * `tgw_wan_vport` - External port number.
* `tgw_wan_vport` - External port number.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeFlowStatus](https://www.tencentcloud.com/document/product/238/32113) - This API is used to query flow status.
* Create: [ModifyCloseWanIp](https://www.tencentcloud.com/document/product/238/72751) - This API is used to disable the public network for the instance.
* Create: [ModifyOpenWanIp](https://www.tencentcloud.com/document/product/238/72750) - This API is used to enable the public network for the instance.
* Read: [DescribeDBInstances](https://www.tencentcloud.com/document/product/238/32115) - This API is used to query the list of instances.
* Read: [DescribeReadOnlyGroupList](https://www.tencentcloud.com/document/product/238/72738) - This API is used to query the read-only group list.
* Update: [DescribeFlowStatus](https://www.tencentcloud.com/document/product/238/32113) - This API is used to query flow status.
* Update: [ModifyCloseWanIp](https://www.tencentcloud.com/document/product/238/72751) - This API is used to disable the public network for the instance.
* Update: [ModifyOpenWanIp](https://www.tencentcloud.com/document/product/238/72750) - This API is used to enable the public network for the instance.

