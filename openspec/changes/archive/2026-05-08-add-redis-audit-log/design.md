# 设计文档：tencentcloud_redis_audit_log

## 资源概览

- **资源名称**: `tencentcloud_redis_audit_log`
- **Package**: `crs`
- **文件位置**: `tencentcloud/services/crs/`
- **唯一 ID**: `InstanceId`（实例级别，一个实例只有一份审计日志配置）

## Schema 设计

```
resource "tencentcloud_redis_audit_log" "example" {
  instance_id        = "crs-xxxxxxxx"
  log_sub_type       = "all"
  log_expire_day     = 7
  high_log_expire_day = 7
  degrade_strategy   = 500
}
```

| 字段 | 类型 | Required/Optional | ForceNew | 说明 |
|------|------|-------------------|----------|------|
| `instance_id` | String | Required | true | Redis 实例 ID |
| `log_sub_type` | String | Required | false | 日志子类型：`write`/`read`/`all` |
| `log_expire_day` | Int | Required | false | 日志有效期（天）：7 或 30 |
| `high_log_expire_day` | Int | Required | false | 高频日志有效期（天）：7 |
| `degrade_strategy` | Int | Optional | false | 降级策略阈值（ms），范围 [300,1000]，默认 500 |

> `LogType` 字段不暴露在 schema 中，代码中固定为 `"auditLog"`。

## CRUD 实现

### Create (`OpenLog`)

1. 从 schema 读取 `instance_id`、`log_sub_type`、`log_expire_day`、`high_log_expire_day`、`degrade_strategy`
2. 构建 `OpenLogRequest`，`LogType` 固定为 `"auditLog"`
3. `resource.Retry(WriteRetryTimeout, ...)` 调用 `OpenLogWithContext`
4. `d.SetId(instanceId)`
5. 调用 `resourceTencentCloudRedisAuditLogRead`

### Read (`DescribeLogs`)

Read 逻辑通过 service 层 `DescribeRedisAuditLogById` 实现：
- 调用 `DescribeLogs`，传入近期时间窗（如过去 1 分钟），`LogType = "auditLog"`
- 若接口无报错且能正常调用，说明审计日志已开启
- **注意**：`DescribeLogs` 返回的是日志记录列表，不是配置信息；Read 主要从本地 state 回填，仅通过接口验证实例存在性
- 实现方案：使用 `DescribeLogs` 并检查是否返回非权限错误（如果实例不存在则 `d.SetId("")`）

### Update (`ModifyLog`)

1. 检查可变字段 `log_sub_type`、`log_expire_day`、`high_log_expire_day`、`degrade_strategy` 是否有变更
2. 若有变更，构建 `ModifyLogRequest`，`LogType` 固定为 `"auditLog"`
3. `resource.Retry(WriteRetryTimeout, ...)` 调用 `ModifyLogWithContext`
4. 调用 `resourceTencentCloudRedisAuditLogRead`

### Delete (`CloseLog`)

1. 从 `d.Id()` 取出 `instanceId`
2. 构建 `CloseLogRequest`，`LogType` 固定为 `"auditLog"`
3. `resource.Retry(WriteRetryTimeout, ...)` 调用 `CloseLogWithContext`

## SDK 扩展

在 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412/` 追加：

### models.go 新增 struct

- `OpenLogRequest` / `OpenLogResponse`
- `ModifyLogRequest` / `ModifyLogResponse`
- `CloseLogRequest` / `CloseLogResponse`
- `DescribeLogsRequest` / `DescribeLogsResponse`
- `LogFilter` (DescribeLogs 的过滤条件)
- `LogResult` (DescribeLogs 返回的日志条目)

### client.go 新增方法

- `NewOpenLogRequest` / `NewOpenLogResponse` / `OpenLog` / `OpenLogWithContext`
- `NewModifyLogRequest` / `NewModifyLogResponse` / `ModifyLog` / `ModifyLogWithContext`
- `NewCloseLogRequest` / `NewCloseLogResponse` / `CloseLog` / `CloseLogWithContext`
- `NewDescribeLogsRequest` / `NewDescribeLogsResponse` / `DescribeLogs` / `DescribeLogsWithContext`

## Service 层

在 `service_tencentcloud_redis.go` 中新增：

```go
func (me *RedisService) DescribeRedisAuditLogById(ctx context.Context, instanceId string) (enabled bool, errRet error)
```

通过调用 `DescribeLogs` 来判断审计日志是否已开启（若调用成功则视为已开启）。

## 代码风格

严格参考 `resource_tc_igtm_strategy.go`：
- 使用 `tccommon.NewResourceLifeCycleHandleFuncContext`
- 使用 `tccommon.LogElapsed` / `tccommon.InconsistentCheck`
- 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)`
- 使用 `helper.String` / `helper.IntInt64` 等辅助函数
- 错误日志使用 `log.Printf("[CRITAL]%s ...")` 格式

## 文件清单

| 文件 | 类型 |
|------|------|
| `vendor/.../redis/v20180412/models.go` | 追加 SDK struct |
| `vendor/.../redis/v20180412/client.go` | 追加 SDK 方法 |
| `tencentcloud/services/crs/service_tencentcloud_redis.go` | 追加 service 方法 |
| `tencentcloud/services/crs/resource_tc_redis_audit_log.go` | 新建 Resource 文件 |
| `tencentcloud/services/crs/resource_tc_redis_audit_log.md` | 新建 Resource 文档 |
| `tencentcloud/services/crs/resource_tc_redis_audit_log_test.go` | 新建测试文件 |
| `tencentcloud/provider.go` | 注册 Resource |
