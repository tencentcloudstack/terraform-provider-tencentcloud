# 任务清单

本修复涉及简单的字段读取补充,任务清单较短且直接。

## 1. 更新资源 Schema 定义

**文件**: `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go`

- [x] 在 `auto_renew_flag` 字段的 schema 定义中添加 `Computed: true` 属性
- [x] 确保 description 准确描述该字段的行为

**验证**: ✅ 编译通过,schema 定义正确

---

## 2. 在 Read 函数中设置 auto_renew_flag

**文件**: `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go`

- [x] 在 `resourceTencentCloudPostgresqlInstanceRead` 函数中,设置 `charge_type` 之后
- [x] 添加代码读取 `instance.AutoRenew` 并设置到 state 中
- [x] 添加 nil 检查以避免空指针异常
- [x] 正确处理类型转换 (`*uint64` → `int`)

**验证**: 
- ✅ 编译通过
- ✅ Linter 检查通过 (无新增错误)
- ✅ 逻辑与只读实例的实现保持一致

---

## 3. 编译和 Linter 验证

- [x] 运行 `go build` 确保代码编译通过
- [x] 运行 linter 检查没有新增警告或错误
- [x] 确认代码风格符合项目规范

**命令**:
```bash
go build -o /dev/null ./tencentcloud/services/postgresql/
```

**结果**: ✅ 编译成功,无新增 linter 错误

---

## 4. 手动功能测试 (可选但推荐)

- [ ] 创建包含 `auto_renew_flag` 的预付费实例
- [ ] 执行 `terraform refresh` 验证字段正确读取
- [ ] 在控制台修改自动续费设置
- [ ] 执行 `terraform plan` 验证能检测配置漂移
- [ ] 测试 `terraform import` 导入现有实例

**测试配置示例**:
```hcl
resource "tencentcloud_postgresql_instance" "test" {
  name              = "test-auto-renew"
  charge_type       = "PREPAID"
  period            = 1
  auto_renew_flag   = 1
  # ... 其他必需字段
}
```

---

## 5. 更新 Changelog (可选)

**文件**: `.changelog/<PR_NUMBER>.txt`

- [ ] 创建 changelog 条目记录此 bug 修复

**格式**:
```
```bug
postgresql: Fix `auto_renew_flag` not being read in `tencentcloud_postgresql_instance` resource
```
```

---

## 总结

- **预计工作量**: ✅ 完成 (约 30 分钟)
- **风险等级**: 极低 (纯粹是补充缺失的字段读取)
- **破坏性变更**: 无
- **测试需求**: ✅ 基础编译测试通过
- **状态**: 🎉 **所有核心任务已完成!**

### 实施的修改

1. ✅ Schema 定义添加 `Computed: true` (line 60)
2. ✅ Read 函数添加 `auto_renew_flag` 设置逻辑 (line 906-908)
3. ✅ 编译验证通过
4. ✅ Linter 验证通过 (无新增错误)
