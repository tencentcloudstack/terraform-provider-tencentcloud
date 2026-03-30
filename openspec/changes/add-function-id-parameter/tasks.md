## 1. 验证 API 兼容性

- [x] 1.1 查阅 CreateFunction API 文档，确认是否支持 FunctionId 参数
- [x] 1.2 如 API 文档不明确，编写简单的测试程序验证 CreateFunction API 是否接受 FunctionId 参数

## 2. 修改资源 Schema

- [x] 2.1 修改 `tencentcloud/services/teo/resource_tc_teo_function.go`，将 function_id 字段的 Schema 定义从 Computed 改为 Computed+Optional
- [x] 2.2 更新 function_id 字段的 Description，说明该字段既可由用户指定，也可由服务端自动生成

## 3. 修改 Create 函数

- [x] 3.1 修改 `resourceTencentCloudTeoFunctionCreate` 函数，在构建 CreateFunctionRequest 前检查用户是否指定了 function_id
- [x] 3.2 如果用户指定了 function_id，将其添加到 request.FunctionId 参数中
- [x] 3.3 确保 CreateFunctionRequest 调用的错误处理逻辑正常工作

## 4. 更新测试用例

- [ ] 4.1 在 `resource_tc_teo_function_test.go` 中添加测试用例：用户创建函数时指定 function_id
- [ ] 4.2 在 `resource_tc_teo_function_test.go` 中添加测试用例：用户创建函数时不指定 function_id（验证向后兼容性）
- [ ] 4.3 在 `resource_tc_teo_function_test.go` 中添加测试用例：读取已存在的函数（验证 function_id 正确返回）

## 5. 更新资源示例文档

- [ ] 5.1 更新 `tencentcloud/services/teo/resource_tc_teo_function.md`，添加包含 function_id 参数的创建示例
- [ ] 5.2 在示例中说明 function_id 参数为可选，不指定时由服务端自动生成

## 6. 生成网站文档

- [ ] 6.1 运行 `make doc` 命令，生成 website/docs/r/teo_function.html.markdown 文档
- [ ] 6.2 验证生成的文档中 function_id 参数的说明正确

## 7. 代码验证

- [ ] 7.1 运行 `go build` 命令，确保代码编译无错误
- [ ] 7.2 运行 `go vet` 命令，确保代码通过静态检查
- [ ] 7.3 运行 `go fmt` 命令，确保代码格式符合规范

## 8. 运行测试

- [ ] 8.1 运行 `go test -v ./tencentcloud/services/teo/... -run TestAccTencentCloudTeoFunction` 命令，运行单元测试
- [ ] 8.2 如需要，运行完整的验收测试（TF_ACC=1），确保新增功能正常工作
- [ ] 8.3 确保所有现有测试通过，验证向后兼容性
