---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_reload_proxy_node"
sidebar_current: "docs-tencentcloud-resource-cynosdb_reload_proxy_node"
description: |-
  Provides a resource to create a cynosdb reload_proxy_node
---

# tencentcloud_cynosdb_reload_proxy_node

Provides a resource to create a cynosdb reload_proxy_node

## Example Usage

```hcl
resource "tencentcloud_cynosdb_reload_proxy_node" "reload_proxy_node" {
  cluster_id     = "cynosdbmysql-cgd2gpwr"
  proxy_group_id = "cynosdbmysql-proxy-8lqtl8pk"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) cluster id.
* `proxy_group_id` - (Required, String, ForceNew) proxy group id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeFlow](https://cloud.tencent.com/document/api/1003/86162) - 查询任务流信息
* Create: [ReloadBalanceProxyNode](https://cloud.tencent.com/document/api/1003/94128) - 负载均衡数据库代理

## Import

cynosdb reload_proxy_node can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_reload_proxy_node.reload_proxy_node reload_proxy_node_id
```

