---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_upgrade_proxy_version"
sidebar_current: "docs-tencentcloud-resource-cynosdb_upgrade_proxy_version"
description: |-
  Provides a resource to create a cynosdb upgrade_proxy_version
---

# tencentcloud_cynosdb_upgrade_proxy_version

Provides a resource to create a cynosdb upgrade_proxy_version

## Example Usage

```hcl
resource "tencentcloud_cynosdb_upgrade_proxy_version" "upgrade_proxy_version" {
  cluster_id        = "cynosdbmysql-bws8h88b"
  dst_proxy_version = "1.3.7"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `dst_proxy_version` - (Required, String) Database Agent Upgrade Version.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeFlow](https://cloud.tencent.com/document/api/1003/86162) - 查询任务流信息
* Create: [DescribeProxies](https://cloud.tencent.com/document/api/1003/94132) - 查询数据库代理列表
* Create: [UpgradeProxyVersion](https://cloud.tencent.com/document/api/1003/94126) - 升级数据库代理版本
* Read: [DescribeProxies](https://cloud.tencent.com/document/api/1003/94132) - 查询数据库代理列表
* Update: [DescribeFlow](https://cloud.tencent.com/document/api/1003/86162) - 查询任务流信息
* Update: [UpgradeProxyVersion](https://cloud.tencent.com/document/api/1003/94126) - 升级数据库代理版本

