## 1. 代码修改

- [x] 1.1 在 tencentcloud_teo_function 资源的 Schema 中添加 function_id 字段定义
  - 字段类型：string
  - Computed: true
  - Optional: false
  - Description: 函数 ID，由腾讯云服务端生成

- [x] 1.2 更新 Create 函数，在调用 CreateFunction API 后从响应中读取 FunctionId
  - 从 API 响应中提取 response.FunctionId
  - 使用 `d.Set("function_id", response.FunctionId)` 设置字段值到 Terraform state

- [x] 1.3 更新 Read 函数，在调用 DescribeFunction 或相关 API 后从响应中读取 FunctionId
  - 从 API 响应中提取 response.FunctionId
  - 使用 `d.Set("function_id", response.FunctionId)` 更新字段值到 Terraform state

## 2. 测试更新

- [x] 2.1 在 resource_tencentcloud_teo_function_test.go 中添加单元测试
  - 测试 Create 函数正确设置 function_id 字段
  - 测试 Read 函数正确设置 function_id 字段
  - 覆盖 FunctionId 为有效值的情况
  - 覆盖 API 不返回 FunctionId 的情况

- [x] 2.2 在 resource_tencentcloud_teo_function_test.go 中添加验收测试（TF_ACC=1）
  - 测试创建资源后 state 中包含正确的 function_id
  - 测试 function_id 值与腾讯云服务端返回的一致
  - 测试 refresh 操作后 function_id 保持一致
  - 验证向后兼容性：不包含 function_id 的配置也能正常工作

## 3. 文档更新

- [x] 3.1 更新 resource_tencentcloud_teo_function.md 示例文件
  - 添加 function_id 字段的示例说明（如适用）

- [x] 3.2 运行 make doc 命令生成文档
  - 确认 website/docs/r/teo_function.md 文档包含 function_id 字段的说明

## 4. 构建和验证

- [x] 4.1 运行单元测试验证代码变更
  - 执行 resource_tencentcloud_teo_function_test.go 中的单元测试
  - 确保所有新增的单元测试通过

- [x] 4.2 运行验收测试验证 API 集成
  - 设置 TF_ACC=1 环境变量
  - 运行 resource_tencentcloud_teo_function_test.go 中的验收测试
  - 确保验收测试通过

- [x] 4.3 验证向后兼容性
  - 确认现有的测试用例仍然通过
  - 确认不包含 function_id 字段的配置能够正常工作
  - 确认旧 state 可以正确升级

- [x] 4.4 运行构建和 lint 检查
  - 执行 go build 确认代码编译通过
  - 执行 gofmt 和 golint 确保代码符合规范
