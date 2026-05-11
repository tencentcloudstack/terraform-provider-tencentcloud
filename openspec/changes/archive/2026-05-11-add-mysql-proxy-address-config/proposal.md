# 变更提案：新增 tencentcloud_mysql_proxy_address_config 资源

## 变更类型

**新功能** — 新增 `tencentcloud_mysql_proxy_address_config` Config 型 Resource，用于管理腾讯云 MySQL 数据库代理（CDB Proxy）的地址配置参数（权重模式、延迟剔除、故障转移、连接池等）。

## Why

腾讯云 MySQL 数据库代理支持丰富的代理地址级别配置，包括负载均衡权重分配模式、延迟剔除开关、故障转移策略、连接池、事务分离等。当前 terraform-provider-tencentcloud 中已有 `tencentcloud_mysql_proxy` 资源管理代理的创建/删除，但代理地址的精细化配置参数无法独立管理。

通过提供 `tencentcloud_mysql_proxy_address_config` Config 资源，用户可以：
- 独立管理已有代理地址的配置，而无需重建代理
- 调整权重分配模式（system/custom）、只读策略、连接池等参数
- 与代理生命周期解耦，支持持续配置漂移检测

## 接口说明

| 操作 | 接口名 | 说明 |
|------|--------|------|
| Create | `AdjustCdbProxyAddress` | Config 型资源，Create 直接调用 Update（无真实创建） |
| Read   | `DescribeCdbProxyInfo` | 通过 InstanceId + ProxyGroupId 查询，筛选 ProxyAddressId 对应配置 |
| Update | `AdjustCdbProxyAddress` | 先校验 DescribeCdbProxyInfo 能唯一查到目标，再调整地址配置 |
| Delete | 无（no-op） | Config 类型资源，删除时只清除 State，不调用 API |

### AdjustCdbProxyAddress 业务入参

| 参数 | 必选 | 类型 | SDK 字段 |
|------|------|------|----------|
| `ProxyGroupId` | 是 | String | `ProxyGroupId` |
| `ProxyAddressId` | 是 | String | `ProxyAddressId` |
| `WeightMode` | 是 | String | `WeightMode` |
| `IsKickOut` | 是 | Bool | `IsKickOut` |
| `MinCount` | 是 | Integer | `MinCount` |
| `MaxDelay` | 是 | Integer | `MaxDelay` |
| `FailOver` | 是 | Bool | `FailOver` |
| `AutoAddRo` | 是 | Bool | `AutoAddRo` |
| `ReadOnly` | 是 | Bool | `ReadOnly` |
| `TransSplit` | 否 | Bool | `TransSplit` |
| `ConnectionPool` | 否 | Bool | `ConnectionPool` |
| `ProxyAllocation` | 否 | Array | `ProxyAllocation` |
| `AutoLoadBalance` | 否 | Bool | `AutoLoadBalance` |
| `AccessMode` | 否 | String | `AccessMode` |
| `ApNodeAsRoNode` | 否 | String | `ApNodeAsRoNode` |
| `ApQueryToOtherNode` | 否 | String | `ApQueryToOtherNode` |

## 唯一 ID

`InstanceId#ProxyGroupId` — `InstanceId` 为 DescribeCdbProxyInfo 必传参数，`ProxyGroupId` 用于 AdjustCdbProxyAddress。

## What Changes

| 文件 | 变更内容 |
|------|---------|
| `tencentcloud/services/cdb/service_tencentcloud_mysql.go` | 新增 `DescribeMysqlProxyAddressConfig` 方法 |
| `tencentcloud/services/cdb/resource_tc_mysql_proxy_address_config.go` | 新增 Resource 主文件 |
| `tencentcloud/services/cdb/resource_tc_mysql_proxy_address_config.md` | 新增 Resource 文档 |
| `tencentcloud/services/cdb/resource_tc_mysql_proxy_address_config_test.go` | 新增单元测试 |
| `tencentcloud/provider.go` | 注册新 Resource |
