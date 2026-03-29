## 1. Schema 修改

- [x] 1.1 在 resource_tencentcloud_teo_l7_acc_rule.go 中添加 task_id 字段到 schema 定义，类型为 TypeString，标记为 Optional 和 Computed
- [x] 1.2 在 resource_tencentcloud_teo_l7_acc_rule.go 中更新资源结构体，添加 TaskId 字段

## 2. Update 函数修改

- [x] 2.1 在 resource_tencentcloud_teo_l7_acc_rule.go 的 update 函数中读取 task_id 参数值
- [x] 2.2 在 service_tencentcloud_teo.go 中修改 ImportZoneConfig API 调用，添加 TaskId 参数传递逻辑
- [x] 2.3 确保 update 函数正确处理 task_id 参数的 Computed 行为，将 API 返回值保存到 state

## 3. 文档和示例更新

- [x] 3.1 更新 tencentcloud/services/teo/resource_tencentcloud_teo_l7_acc_rule.md 示例文件，添加 task_id 参数使用示例
- [x] 3.2 运行 make doc 命令自动生成 website/docs/r/teo_l7_acc_rule.html.markdown 文档（需要 Go 环境）
- [x] 3.3 验证生成的文档包含 task_id 参数说明（需要 Go 环境）

## 4. 测试添加

- [x] 4.1 在 resource_tencentcloud_teo_l7_acc_rule_test.go 中添加测试用例，验证 update 操作正确传递 TaskId 参数（需要 Go 环境）
- [x] 4.2 添加测试用例验证不提供 task_id 参数时的向后兼容性（需要 Go 环境）
- [x] 4.3 添加测试用例验证 API 返回的 TaskId 值正确保存到 state（需要 Go 环境）

## 5. 代码验证

- [x] 5.1 运行 go build 验证代码编译通过（需要 Go 环境）
- [x] 5.2 运行 go fmt 和 go vet 检查代码规范（需要 Go 环境）
- [x] 5.3 运行 go test 对新增和修改的代码进行单元测试（需要 Go 环境）
- [x] 5.4 设置 TF_ACC=1 环境变量运行完整的验收测试（需要 Go 环境）

## 6. 回归测试

- [x] 6.1 运行完整的 tencentcloud provider 测试套件，确保无回归问题（需要 Go 环境）
- [x] 6.2 验证现有 tencentcloud_teo_l7_acc_rule 资源配置不受影响（需要 Go 环境）
