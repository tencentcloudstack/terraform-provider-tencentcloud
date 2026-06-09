# 任务清单：fix-tcr-instance-endpoint-wait-logic

## 1. 移动等待逻辑

**文件**: `tencentcloud/services/tcr/resource_tc_tcr_instance.go`

- [x] 将 `security_policy` 变更块开头的等待逻辑从 `security_policy` 块中移除
- [x] 将该等待逻辑移至 `open_public_operation` 变更块末尾，且仅在 `operation == true` 时执行
- [x] `security_policy` 变更块中仅保留 add/remove 逻辑（去掉多余的 `var err error` 声明）

---

## 2. 格式化代码

- [x] 执行 `go fmt ./tencentcloud/services/tcr/`
- [x] 执行 `go build ./tencentcloud/services/tcr/` 确认编译通过

---

## 3. 完善等待逻辑（同时处理 Delete 场景）

**文件**: `tencentcloud/services/tcr/resource_tc_tcr_instance.go`

- [x] 去掉 `if operation` 条件限制，改为无论 Create/Delete 都执行等待
- [x] `operation=true` 时：`Opened` → 成功，`Opening` → 重试，其他 → 报错
- [x] `operation=false` 时：`Closed` → 成功，`Deleting` → 重试，其他 → 报错
- [x] 执行 `go fmt` 及编译验证通过

---

## 总结

- **状态**: 🎉 所有任务已完成
- **预计工作量**：极低（纯代码搬移）
- **风险等级**：极低
- **破坏性变更**：无
