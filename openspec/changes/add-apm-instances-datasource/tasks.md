# Tasks: Add APM Instances Data Source

## Preparation
- [ ] 确认 APM service 层已支持 DescribeApmInstances（已存在 DescribeApmInstanceById 方法，需确认是否需要新增列表查询方法）
- [ ] 确认 SDK 版本支持 DescribeApmInstances 接口（tencentcloud-sdk-go/tencentcloud/apm/v20210622）

## Implementation
- [ ] 创建 data source 文件 `data_source_tc_apm_instances.go`
  - 定义 schema，包含过滤条件和输出字段
  - 实现 Read 方法调用 service 层
  - 处理响应数据并设置到 ResourceData
- [ ] 创建或更新 service 层方法
  - 如需要，在 `service_tencentcloud_apm.go` 中添加 `DescribeApmInstances` 方法
  - 支持多种过滤条件（instance_ids, instance_id, instance_name, tags 等）
- [ ] 在 provider.go 中注册 data source
  - 添加 `"tencentcloud_apm_instances": apm.DataSourceTencentCloudApmInstances()`
- [ ] 创建文档文件 `data_source_tc_apm_instances.md`
  - 添加一句话描述
  - 添加使用示例（基本查询、按 ID 过滤、按名称过滤、按标签过滤）

## Testing
- [ ] 创建测试文件 `data_source_tc_apm_instances_test.go`
  - 实现基本的验收测试
  - 测试各种过滤条件
  - 验证返回的字段完整性
- [ ] 本地运行 `make fmt` 格式化代码
- [ ] 本地运行 `make lint` 检查代码质量
- [ ] 本地运行测试验证功能

## Documentation
- [ ] 更新 `tencentcloud/provider.md`
  - 在 Application Performance Management(APM) 部分的 Data Source 节添加 `tencentcloud_apm_instances`
- [ ] 运行 `make doc` 生成文档
- [ ] 验证生成的 `website/docs/d/apm_instances.html.markdown` 文件

## Validation
- [ ] 运行 `openspec validate add-apm-instances-datasource --strict` 验证提案
- [ ] 确认所有文件已创建且格式正确
- [ ] 确认测试通过
