# 变更提案：新增 tencentcloud_redis_audit_log 资源

## 变更类型

**新功能** — 新增 `tencentcloud_redis_audit_log` Resource，用于管理腾讯云 Redis 实例的审计日志（开启/修改/关闭）。

## Why

腾讯云 Redis 提供了审计日志功能，允许用户对实例的读/写/全部命令进行审计记录。当前 terraform-provider-tencentcloud 中没有对应的 Resource，用户无法通过 Terraform 管理 Redis 审计日志配置。

通过提供 `tencentcloud_redis_audit_log` Resource，用户可以：
- 使用 IaC 方式统一管理 Redis 审计日志开关
- 配置审计日志子类型（读/写/全部）
- 配置日志有效期、高频日志有效期、降级策略阈值

## 接口说明

| 操作 | 接口名 | 说明 |
|------|--------|------|
| Create | `OpenLog` | 开启审计日志 |
| Read   | `DescribeLogs` | 查询日志配置（用于状态读取，判断是否已开启） |
| Update | `ModifyLog` | 修改审计日志配置 |
| Delete | `CloseLog` | 关闭审计日志 |

### OpenLog / ModifyLog 入参

| 参数 | 必选 | 类型 | 说明 |
|------|------|------|------|
| `InstanceId` | 是 | String | 实例 ID |
| `LogType` | 是 | String | 固定为 `auditLog`（代码层面固定，不暴露在 schema 中） |
| `LogSubType` | 是 | String | `write`/`read`/`all` |
| `LogExpireDay` | 是 | Integer | 日志有效期（天），枚举：7 / 30 |
| `HighLogExpireDay` | 是 | Integer | 高频日志有效期（天），枚举：7 |
| `DegradeStrategy` | 否（OpenLog）/ 是（ModifyLog） | Integer | 降级策略阈值（ms），范围 [300, 1000] |

## 唯一 ID

资源唯一 ID 为 `InstanceId`，即 `d.SetId(instanceId)`。

## What Changes

| 文件 | 变更内容 |
|------|---------|
| `vendor/.../redis/v20180412/models.go` | 新增 OpenLog/ModifyLog/CloseLog/DescribeLogs 相关 Request/Response struct |
| `vendor/.../redis/v20180412/client.go` | 新增 4 个接口方法 |
| `tencentcloud/services/crs/service_tencentcloud_redis.go` | 新增 DescribeRedisAuditLogById 方法 |
| `tencentcloud/services/crs/resource_tc_redis_audit_log.go` | 新增 Resource 主文件 |
| `tencentcloud/services/crs/resource_tc_redis_audit_log.md` | 新增 Resource 文档 |
| `tencentcloud/services/crs/resource_tc_redis_audit_log_test.go` | 新增单元测试 |
| `tencentcloud/provider.go` | 注册新 Resource |
