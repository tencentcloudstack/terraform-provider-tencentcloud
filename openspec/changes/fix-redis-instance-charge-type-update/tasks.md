# 任务清单：fix-redis-instance-charge-type-update

- [x] 检查 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412/` 是否有 `ModifyInstanceChargeType` 相关结构体
- [x] SDK v1.3.73 已包含该接口，vendor 中已有（14 处匹配），无需升级

---

## 1. service 层新增 ModifyInstanceChargeType 方法

**文件**: `tencentcloud/services/crs/service_tencentcloud_redis.go`

- [x] 新增 `ModifyInstanceChargeType(ctx, instanceId, chargeType string, period int) (errRet error)` 方法
- [x] 执行 `go fmt ./tencentcloud/services/crs/`

---

## 2. schema 修改

**文件**: `tencentcloud/services/crs/resource_tc_redis_instance.go`

- [x] `charge_type` 字段：去掉 `ForceNew: true`，更新 Description
- [x] `prepaid_period` 字段：从 `unsupportedUpdateFields` 列表中移除

---

## 3. update 函数新增 charge_type 变更处理

**文件**: `tencentcloud/services/crs/resource_tc_redis_instance.go`

- [x] 新增 `d.HasChange("charge_type")` 处理块
- [x] 转 PREPAID 时读取 `prepaid_period`（必填）
- [x] 调用 `ModifyInstanceChargeType` + `CheckRedisUpdateOk` 等待完成
- [x] 执行 `go fmt ./tencentcloud/services/crs/`

---

## 4. 编译验证

- [x] `go build ./tencentcloud/services/crs/` 确认编译通过

---

## 总结

- **预计工作量**：中等
- **风险等级**：中（计费类型变更是高危操作）
- **破坏性变更**：轻微（去掉 ForceNew，存量用户行为由重建变为原地修改）
- **状态**: 已完成
