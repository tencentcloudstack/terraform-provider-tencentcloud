## 1. Schema 定义

- [x] 1.1 在 `resource_tc_sms_template.go` 的 Schema map 中添加 `status_code` 字段
- [x] 1.2 设置字段类型为 `schema.TypeInt`
- [x] 1.3 设置字段为 `Computed: true`（只读输出字段，用户不可配置）
- [x] 1.4 设置字段为 `ForceNew: false`
- [x] 1.5 添加清晰的 Description，说明字段含义及可选值（0: 审核通过并已生效, 1: 审核中, 2: 审核通过待生效, -1: 审核未通过或审核失败）
- [x] 1.6 将字段放置在 `remark` 字段之后（保持现有字段顺序不变，新增字段放在最后）

## 2. Read 函数更新

- [x] 2.1 在 `resourceTencentCloudSmsTemplateRead` 函数中，在已有 `d.Set` 调用之后（`international` 字段设置之后），添加对 `status_code` 字段的读取
- [x] 2.2 使用 `if template.StatusCode != nil { _ = d.Set("status_code", template.StatusCode) }` 模式进行 nil 安全检查

## 3. 资源文档更新

- [x] 3.1 更新 `tencentcloud/services/sms/resource_tc_sms_template.md` 文件，在 Example Usage 中展示 `status_code` 作为输出属性
- [x] 3.2 在文档中添加对 `status_code` 字段的描述和可选值说明

## 4. 代码质量检查

- [ ] 4.1 运行 `gofmt` 格式化代码（由 tfpacer-finalize 阶段执行）
- [x] 4.2 验证所有新增代码有适当的 nil 检查和错误处理
- [x] 4.3 确保 `_ = d.Set()` 模式正确处理（忽略未使用的变量错误）

## 5. 生成文档

- [ ] 5.1 运行 `make doc` 生成最终的 `website/docs/` 文档（由 tfpacer-finalize 阶段执行）