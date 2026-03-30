## 1. Pre-implementation Research

- [x] 1.1 查阅腾讯云 TEO CreateFunction API 文档，确认 FunctionId 参数是否为可选输入参数
- [x] 1.2 确认 FunctionId 的格式限制和命名规范（如果有）
- [x] 1.3 检查现有的测试用例，了解测试模式和要求

## 2. Schema Modification

- [x] 2.1 修改 tencentcloud/services/teo/resource_tc_teo_function.go 中的 ResourceTencentCloudTeoFunction 函数
- [x] 2.2 将 function_id 字段从 Computed 改为 Optional + Computed（同时设置两个属性）
- [x] 2.3 更新 function_id 字段的 Description，说明用户可以自定义指定，也可以由系统生成

## 3. Create Function Implementation

- [x] 3.1 修改 resourceTencentCloudTeoFunctionCreate 函数，添加检查用户是否提供 function_id 的逻辑
- [x] 3.2 如果用户提供了 function_id，在 CreateFunctionRequest 中设置 FunctionId 参数
- [x] 3.3 如果用户未提供 function_id，保持原有逻辑（不设置该参数）
- [x] 3.4 确保 resourceTencentCloudTeoFunctionCreate 的错误处理逻辑正确（API 错误正确传递给用户）

## 4. Test Cases Development

- [x] 4.1 在 tencentcloud/services/teo/resource_tc_teo_function_test.go 中添加测试用例：指定自定义 FunctionId 创建资源
- [x] 4.2 在 tencentcloud/services/teo/resource_tc_teo_function_test.go 中添加测试用例：不指定 FunctionId 创建资源（向后兼容测试）
- [x] 4.3 在 tencentcloud/services/teo/resource_tc_teo_function_test.go 中添加测试用例：更新已创建的资源（验证 FunctionId 不可变）
- [x] 4.4 在 tencentcloud/services/teo/resource_tc_teo_function_test.go 中添加测试用例：导入现有资源

## 5. Documentation Update

- [x] 5.1 更新 tencentcloud/services/teo/resource_tc_teo_function.md，添加指定 FunctionId 的示例配置
- [x] 5.2 更新 function_id 字段的文档说明，标明为 Optional 字段
- [ ] 5.3 运行 `make doc` 命令自动生成 website/docs/ 下的 markdown 文档（需要在有 Go 环境的情况下执行）
- [ ] 5.4 验证生成的文档内容正确且完整（需要在有 Go 环境的情况下执行）

## 6. Verification and Testing

- [ ] 6.1 运行 Go build 命令，确保代码编译无错误
- [ ] 6.2 运行现有的验收测试（TF_ACC=1），验证向后兼容性
- [ ] 6.3 运行新的测试用例，验证 FunctionId 输入功能正常工作
- [ ] 6.4 手动测试：创建一个包含自定义 FunctionId 的资源，验证资源创建成功且 ID 正确
- [ ] 6.5 手动测试：创建一个不包含 FunctionId 的资源，验证资源创建成功且系统生成 FunctionId
- [ ] 6.6 手动测试：更新已创建的资源，验证 FunctionId 保持不变
- [ ] 6.7 手动测试：导入现有资源，验证资源状态正确

## 7. Code Review and Cleanup

- [x] 7.1 代码自我审查，确保代码风格符合项目规范
- [x] 7.2 检查是否所有必要的日志和错误处理都已添加
- [x] 7.3 确保没有遗留的调试代码或注释
- [ ] 7.4 验证所有测试用例都能通过（需要在有 Go 环境的情况下执行）
