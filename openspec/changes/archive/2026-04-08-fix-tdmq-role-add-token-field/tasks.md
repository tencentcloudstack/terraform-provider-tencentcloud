# 任务清单：fix-tdmq-role-add-token-field

## 1. schema 新增 token 字段

**文件**: `tencentcloud/services/tpulsar/resource_tc_tdmq_role.go`

- [x] 在 `remark` 字段之后新增 `token` 字段：
  ```go
  "token": {
      Type:        schema.TypeString,
      Computed:    true,
      Sensitive:   true,
      Description: "Role token. This field is returned by the API and used for authentication.",
  },
  ```

---

## 2. Read 模块新增 token 赋值

**文件**: `tencentcloud/services/tpulsar/resource_tc_tdmq_role.go`

- [x] 在 `d.Set("remark", info.Remark)` 之后新增：
  ```go
  _ = d.Set("token", info.Token)
  ```
- [x] 执行 `go fmt ./tencentcloud/services/tpulsar/`

---

## 3. 编译验证

- [x] `go build ./tencentcloud/services/tpulsar/` 确认编译通过

---

## 总结

- **预计工作量**：极小（约 5 分钟）
- **风险等级**：极低（纯增量 Computed 字段）
- **破坏性变更**：无
- **状态**: 已完成

---

## 4. ~~支持 Import 功能~~（已回滚，不可实施）

**决策**: 修改资源 ID 格式属于破坏性变更，会导致所有存量用户 state 不兼容，严禁实施。已完整回滚至原始代码。

该资源 ID 为 `role_name`，Read 依赖 state 中的 `cluster_id`，在不改变 ID 格式的前提下**无法支持 Import**。如需支持，需等待下一个 major 版本。
