# 任务清单：add-bh-acl-resource

## 1. 新增 service 层方法

**文件**: `tencentcloud/services/bh/service_tencentcloud_bh.go`

- [x] 在文件末尾新增 `DescribeBhAclById(ctx, aclId string) (*bhv20230418.Acl, error)` 方法
- [x] 执行 `go fmt ./tencentcloud/services/bh/` 格式化

---

## 2. 新增资源主文件

**文件**: `tencentcloud/services/bh/resource_tc_bh_acl.go`

- [x] 实现 `ResourceTencentCloudBhAcl()` 函数，完整 Schema
- [x] 实现 `resourceTencentCloudBhAclCreate`
- [x] 实现 `resourceTencentCloudBhAclRead`
- [x] 实现 `resourceTencentCloudBhAclUpdate`
- [x] 实现 `resourceTencentCloudBhAclDelete`
- [x] 执行 `go fmt ./tencentcloud/services/bh/` 格式化

---

## 3. 新增资源文档示例

**文件**: `tencentcloud/services/bh/resource_tc_bh_acl.md`

- [x] 编写资源示例 HCL，包含必填字段和主要可选字段

---

## 4. 新增单元测试文件

**文件**: `tencentcloud/services/bh/resource_tc_bh_acl_test.go`

- [x] 实现 `TestAccTencentCloudBhAclResource_basic`
- [x] 编写 `testAccBhAcl` 和 `testAccBhAclUpdate` HCL 常量

---

## 5. 注册资源

**文件**: `tencentcloud/provider.go`

- [x] 在 `ResourcesMap` 中追加 `"tencentcloud_bh_acl": bh.ResourceTencentCloudBhAcl()`

---

## 6. 编译验证

- [x] 运行 `go build ./tencentcloud/services/bh/` 确认编译通过
- [x] 运行 `go build ./tencentcloud/` 确认 provider 编译通过

---

## 7. 修正 int/bool 字段取值方式

**文件**: `tencentcloud/services/bh/resource_tc_bh_acl.go`

- [x] create 函数中 `max_file_up_size`、`max_file_down_size`、`max_access_credential_duration` 由 `d.GetOk` 改为 `d.GetOkExists`
- [x] update 函数中相同的 3 个 int 字段同步修改
- [x] 确保 create/update 函数中所有 bool/int 类型字段均使用 `d.GetOkExists`
- [x] 执行 `go fmt` 及编译验证通过

---

## 8. 新增 DescribeBhDepartments service 方法

**文件**: `tencentcloud/services/bh/service_tencentcloud_bh.go`

- [x] `DescribeBhDepartments` 方法已存在（由其他变更预先实现），无需新增

---

## 9. Read 函数中 department_id 添加 Enabled 判断

**文件**: `tencentcloud/services/bh/resource_tc_bh_acl.go`

- [x] 在 `resourceTencentCloudBhAclRead` 中调用 `service.DescribeBhDepartments`
- [x] 仅当 `Departments != nil && Departments.Enabled != nil && *Departments.Enabled == true` 时才执行 `d.Set("department_id", ...)`

---

## 10. 编译验证 + go fmt

- [x] `go fmt ./tencentcloud/services/bh/`
- [x] `go build ./tencentcloud/services/bh/`

---

## 总结

- **状态**: 🎉 所有任务已完成
- **风险等级**：低
- **破坏性变更**：无
