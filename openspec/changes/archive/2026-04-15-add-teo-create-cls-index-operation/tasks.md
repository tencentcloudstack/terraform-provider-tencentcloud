## 1. 代码实现

- [x] 1.1 创建资源文件 `tencentcloud/services/teo/resource_tc_teo_create_cls_index_operation.go`，定义 Resource 函数和 Schema
- [x] 1.2 实现资源 Create 函数，调用 CreateCLSIndex API
- [x] 1.3 实现资源 Read 函数（空实现）
- [x] 1.4 实现资源 Delete 函数（空实现）
- [x] 1.5 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中注册新资源

## 2. 单元测试

- [x] 2.1 创建测试文件 `tencentcloud/services/teo/resource_tc_teo_create_cls_index_operation_test.go`
- [x] 2.2 实现 Create 操作成功场景的单元测试
- [x] 2.3 实现 API 错误处理的单元测试
- [x] 2.4 实现缺少必需参数的单元测试

## 3. 文档更新

- [x] 3.1 创建资源示例文件 `tencentcloud/services/teo/resource_tc_teo_create_cls_index_operation.md`
- [x] 3.2 在示例文件中添加使用说明和 HCL 示例

## 4. 代码格式化和文档生成

- [x] 4.1 使用 gofmt 格式化新增的 Go 代码文件（已执行，无需修改）
- [x] 4.2 运行 `make doc` 命令生成 website/docs/ 下的文档（已执行）

## 5. 验证

- [ ] 5.1 验证新增的 Go 代码文件语法正确性（将在其他验证流程执行）
- [ ] 5.2 验证单元测试可以正常运行（将在其他验证流程执行）
- [ ] 5.3 验证资源注册正确（resource 函数可以正确返回）（将在其他验证流程执行）
- [ ] 5.4 验证文档生成正确（website/docs/ 下生成对应的文档文件）（已执行）
