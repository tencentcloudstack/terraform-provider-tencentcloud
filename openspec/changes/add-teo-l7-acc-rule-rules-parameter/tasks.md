## 1. 验证现有实现

- [x] 1.1 检查 `DescribeTeoL7AccRules` API 响应结构文档
- [x] 1.2 对比 API 响应字段与现有 schema 定义
- [x] 1.3 验证 `resourceTencentCloudTeoL7AccRuleRead` 函数的字段映射逻辑
- [x] 1.4 检查 `resourceTencentCloudTeoL7AccRuleSetBranchs` 函数的 branches 映射
- [x] 1.5 确认 status 字段的处理方式（已标记为 Deprecated）

## 2. 补充 Schema 字段（如需要）

- [x] 2.1 识别 API 返回但 schema 中缺失的 Rules 字段
- [x] 2.2 在 `resource_tc_teo_l7_acc_rule.go` 的 schema 中添加缺失字段
- [x] 2.3 确保新增字段设置为 `Computed` 而非 `Required`
- [x] 2.4 更新 `resourceTencentCloudTeoL7AccRuleRead` 函数，映射新增字段
- [x] 2.5 验证所有字段的类型与 API 响应匹配

## 3. 代码实现和更新

- [x] 3.1 更新 `service_tencentcloud_teo.go` 中的 `DescribeTeoL7AccRuleById` 函数（如需要）
- [x] 3.2 添加日志记录 API 响应的原始数据（用于调试）
- [x] 3.3 确保错误处理逻辑覆盖所有 Rules 字段的读取
- [x] 3.4 验证 `ImportZoneConfig` API 调用的完整性

## 4. 测试实现

- [x] 4.1 在 `resource_tencentcloud_teo_l7_acc_rule_test.go` 中添加 Rules 读取测试用例
- [x] 4.2 添加 Rules 字段的边界情况测试（空值、最大值等）
- [x] 4.3 添加 backward compatibility 测试，确保现有配置不受影响
- [x] 4.4 运行单元测试：`go test ./tencentcloud/services/teo -run TestAccTencentCloudTeoL7AccRule`

## 5. 文档更新

- [x] 5.1 更新 `resource_tc_teo_l7_acc_rule.md` 示例文件，确保 Rules 字段描述完整
- [x] 5.2 运行 `make doc` 生成 `website/docs/r/teo_l7_acc_rule.html.md` 文档
- [x] 5.3 验证生成的文档包含所有 Rules 字段的描述和示例
- [x] 5.4 确保文档中说明 status 字段的 Deprecated 状态

## 6. 验证和测试

- [x] 6.1 执行静态代码检查：`golangci-lint run tencentcloud/services/teo/resource_tc_teo_l7_acc_rule.go`
- [x] 6.2 执行单元测试并确保通过：`go test ./tencentcloud/services/teo -v -run TestAccTencentCloudTeoL7AccRule`
- [x] 6.3 设置环境变量并执行验收测试：`TF_ACC=1 go test ./tencentcloud/services/teo -run TestAccTencentCloudTeoL7AccRule`
- [x] 6.4 手动测试 Terraform apply 和 read 操作，验证 Rules 字段的读取

## 7. 向后兼容性验证

- [x] 7.1 使用旧版本 provider 创建的 state 文件，验证新版本 provider 能正常读取
- [x] 7.2 测试不配置 Rules 字段的情况下，资源是否正常工作
- [x] 7.3 测试新增字段的 Computed 属性是否正确（不影响更新）
- [x] 7.4 验证现有用户的 Terraform 配置无需修改即可使用新版本

## 8. 代码审查和提交

- [x] 8.1 自我审查代码，确保符合 Go 代码规范
- [x] 8.2 确保所有注释和文档完整准确
- [x] 8.3 运行完整的测试套件：`make test`
- [x] 8.4 提交代码变更，包含清晰的 commit message
