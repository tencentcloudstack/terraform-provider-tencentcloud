# 任务清单：fix-identity-center-permission-policy-panic

## 1. 修复 Read 函数 nil pointer dereference

**文件**: `tencentcloud/services/tco/resource_tc_identity_center_role_configuration_permission_policy_attachment.go`

- [x] 在 for 循环中，访问 `*r.RolePolicyId` 前先判断 `r.RolePolicyId != nil`，避免 nil 解引用：
  ```go
  for _, r := range respData.RolePolicies {
      if r.RolePolicyId != nil && *r.RolePolicyId == rolePolicyId {
          rolePolicie = r
          break
      }
  }
  ```
- [x] 循环结束后判断 `rolePolicie == nil`，视为资源已删除，`d.SetId("")` 后 return nil：
  ```go
  if rolePolicie == nil {
      d.SetId("")
      log.Printf("[WARN]%s policy [%d] not found in role configuration [%s], treating as deleted.\n", logId, rolePolicyId, roleConfigurationId)
      return nil
  }
  ```
- [x] 执行 `go fmt ./tencentcloud/services/tco/`

---

## 2. 编译验证

- [x] `go build ./tencentcloud/services/tco/` 确认编译通过

---

## 总结

- **预计工作量**：极小（约 5 分钟）
- **风险等级**：P0，provider 崩溃，必须修复
- **破坏性变更**：无
- **状态**: 已完成

---

## 3. 补充 r != nil 守卫（post-apply 补充）

**背景**: `respData.RolePolicies` 是 `[]*organization.RolePolicie`，slice 中的元素本身也可能为 nil，循环时需先判断 `r != nil` 再访问其成员，否则同样会 panic。

**文件**: `tencentcloud/services/tco/resource_tc_identity_center_role_configuration_permission_policy_attachment.go`

- [x] 将循环条件改为：`if r != nil && r.RolePolicyId != nil && *r.RolePolicyId == rolePolicyId`
- [x] 执行 `go fmt ./tencentcloud/services/tco/`
- [x] `go build ./tencentcloud/services/tco/` 确认编译通过

---

## 4. 修复 role_policy_id 类型不匹配（post-apply 补充）

**背景**: `rolePolicyId` 是 `int64`，但 schema 中 `role_policy_id` 为 `TypeInt`，直接 `d.Set` 会类型不匹配导致静默失败（state 不一致）。

- [x] 将 `d.Set("role_policy_id", rolePolicyId)` 改为 `d.Set("role_policy_id", int(rolePolicyId))`
- [x] 执行 `go fmt ./tencentcloud/services/tco/`
- [x] `go build ./tencentcloud/services/tco/` 确认编译通过
