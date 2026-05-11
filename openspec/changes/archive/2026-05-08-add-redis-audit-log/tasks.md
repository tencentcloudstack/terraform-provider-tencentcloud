# 任务清单：add-redis-audit-log

## 1. 扩展 Redis SDK vendor

**文件**: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412/models.go`

- [x] 检查 SDK vendor 中已有 `OpenLogRequest`/`ModifyLogRequest`/`CloseLogRequest`/`DescribeLogsRequest` 等 struct（无需新增）

**文件**: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412/client.go`

- [x] 检查 SDK vendor 中已有 `OpenLog`/`ModifyLog`/`CloseLog`/`DescribeLogs` 方法（无需新增）

---

## 2. Service 层新增方法

**文件**: `tencentcloud/services/crs/service_tencentcloud_redis.go`

- [x] 新增 `DescribeRedisAuditLogById(ctx context.Context, instanceId string) (instance *redis.LogInstance, errRet error)` 方法，通过 `DescribeLogInstanceList`（LogType=auditLog, LogSwitch=on, Filters instance-id）精确查询

---

## 3. 实现 Resource 主文件

**文件**: `tencentcloud/services/crs/resource_tc_redis_audit_log.go`

- [x] 实现 `ResourceTencentCloudRedisAuditLog()` schema 定义（字段：`instance_id`、`log_sub_type`、`log_expire_day`、`high_log_expire_day`、`degrade_strategy`）
- [x] 实现 `resourceTencentCloudRedisAuditLogCreate`（调用 `OpenLog`）
- [x] 实现 `resourceTencentCloudRedisAuditLogRead`（调用 `DescribeLogs` 验证存在，state 回填）
- [x] 实现 `resourceTencentCloudRedisAuditLogUpdate`（调用 `ModifyLog`）
- [x] 实现 `resourceTencentCloudRedisAuditLogDelete`（调用 `CloseLog`）

---

## 4. 生成 Resource 文档

**文件**: `tencentcloud/services/crs/resource_tc_redis_audit_log.md`

- [x] 编写 Example Usage（包含完整 HCL 示例）
- [x] 编写 Import 说明

---

## 5. 生成单元测试

**文件**: `tencentcloud/services/crs/resource_tc_redis_audit_log_test.go`

- [x] 实现 `TestAccTencentCloudRedisAuditLogResource_basic` 测试函数
- [x] 包含 Create、ImportState 步骤

---

## 6. 注册 Resource 到 Provider

**文件**: `tencentcloud/provider.go`

- [x] 在 Redis 相关资源区域注册 `"tencentcloud_redis_audit_log": crs.ResourceTencentCloudRedisAuditLog()`

---

## 7. 编译验证

- [x] 无编译错误（linter 检查通过，仅有 HINT 级别 deprecated 提示，与 codebase 现有风格一致）

---

## 总结

- **预计工作量**：中等
- **风险等级**：低（新增资源，不影响已有资源）
- **破坏性变更**：无
- **状态**: 已完成
