## 1. Schema 修改

- [x] 1.1 在 tencentcloud_teo_function 资源的 Schema 中新增 FunctionIds 字段（list, 可选）
- [x] 1.2 确保字段类型为 `schema.TypeList` + `schema.TypeString`
- [x] 1.3 确保字段属性为 Optional（非 Computed）

## 2. Read 函数修改

- [x] 2.1 在 Read 函数中添加 FunctionIds 参数的获取逻辑
- [x] 2.2 将 FunctionIds 参数传递给 DescribeFunctions API 请求
- [x] 2.3 确保 FunctionIds 参数为空时，Read 函数仍能正常工作（向后兼容）

## 3. 单元测试

- [x] 3.1 新增测试用例：不使用 FunctionIds 字段时的正常读取
- [x] 3.2 新增测试用例：使用 FunctionIds 字段过滤单个函数 ID
- [x] 3.3 新增测试用例：使用 FunctionIds 字段过滤多个函数 ID
- [x] 3.4 新增测试用例：使用无效的函数 ID 时的错误处理

## 4. 验证和构建

- [x] 4.1 运行资源单元测试：`go test -v ./tencentcloud/services/teo -run TestAccTencentCloudTeoFunction`
- [x] 4.2 运行格式检查：`gofmt -s -w tencentcloud/services/teo/resource_tc_teo_function.go`
- [x] 4.3 运行 lint 检查：`golangci-lint run tencentcloud/services/teo/resource_tc_teo_function.go`

## 5. 文档更新

- [x] 5.1 更新 resource_tc_teo_function.md 示例文件，添加 FunctionIds 字段的用法示例
- [x] 5.2 运行文档生成命令：`make doc`
- [x] 5.3 验证生成的文档包含 FunctionIds 字段的说明
- [ ] 5.2 运行文档生成命令：`make doc`
- [ ] 5.3 验证生成的文档包含 FunctionIds 字段的说明
