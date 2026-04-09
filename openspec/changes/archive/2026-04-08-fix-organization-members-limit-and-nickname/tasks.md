# 任务清单：fix-organization-members-limit-and-nickname

## 1. 调整 DescribeOrganizationMembersByFilter 的 limit 上限

**文件**: `tencentcloud/services/tco/service_tencentcloud_organization.go`

- [x] 将 `limit uint64 = 20` 改为 `limit uint64 = 50`
- [x] 执行 `go fmt ./tencentcloud/services/tco/`

---

## 2. 新增 nick_name 字段到 schema

**文件**: `tencentcloud/services/tco/data_source_tc_organization_members.go`

- [x] 在 `items` 的 schema 中，`permission_status` 字段之后新增 `nick_name` 字段：
  ```go
  "nick_name": {
      Type:        schema.TypeString,
      Computed:    true,
      Description: "Tencent Cloud nickname. Note: This field may return null, indicating that no valid values can be obtained.",
  },
  ```
- [x] 在 read 逻辑中，`permission_status` 赋值之后新增 NickName 赋值：
  ```go
  if orgMember.NickName != nil {
      orgMemberMap["nick_name"] = orgMember.NickName
  }
  ```
- [x] 执行 `go fmt ./tencentcloud/services/tco/`

---

## 3. 编译验证

- [x] `go build ./tencentcloud/services/tco/` 确认编译通过

---

## 总结

- **预计工作量**：小（约 15 分钟）
- **风险等级**：极低（纯增量改动，不破坏已有逻辑）
- **破坏性变更**：无
- **状态**: 已完成
