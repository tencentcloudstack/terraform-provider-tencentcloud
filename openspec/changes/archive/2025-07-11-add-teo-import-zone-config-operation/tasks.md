## 1. 资源代码实现

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_import_zone_config_operation.go`，定义 Schema（zone_id、content 为 Required/ForceNew；task_id、status、message、import_time、finish_time 为 Computed）
- [x] 1.2 实现 Create 函数：调用 ImportZoneConfig API，获取 TaskId 后轮询 DescribeZoneConfigImportResult 直到任务完成（success/failure），设置 Computed 属性
- [x] 1.3 实现 Read 函数（空操作，return nil）和 Delete 函数（空操作，return nil）

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册资源 `tencentcloud_teo_import_zone_config`
- [x] 2.2 在 `tencentcloud/provider.md` 中添加资源文档条目

## 3. 资源文档

- [x] 3.1 创建 `tencentcloud/services/teo/resource_tc_teo_import_zone_config_operation.md`，包含一句话描述、Example Usage

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_import_zone_config_operation_test.go`，使用 gomonkey mock ImportZoneConfig 和 DescribeZoneConfigImportResult API
- [x] 4.2 编写测试用例：成功导入场景（轮询返回 success）
- [x] 4.3 编写测试用例：导入失败场景（轮询返回 failure）
- [x] 4.4 运行 `go test` 验证单元测试通过

## 5. 代码正确性验证

- [x] 5.1 检查 ImportZoneConfig API 入参（ZoneId、Content）与 Terraform Schema 的映射正确性
- [x] 5.2 检查 DescribeZoneConfigImportResult API 入参（ZoneId、TaskId）与轮询逻辑的一致性
- [x] 5.3 检查 DescribeZoneConfigImportResult 出参（Status、Message、Content、ImportTime、FinishTime）与 Computed 属性的映射正确性
