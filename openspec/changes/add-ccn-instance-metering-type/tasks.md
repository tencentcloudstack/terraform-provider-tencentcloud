# 实现任务清单

## 1. Schema 定义
- [x] 1.1 在 `resource_tc_ccn.go` 的 Schema map 中添加 `instance_metering_type` 字段
- [x] 1.2 设置字段类型为 `schema.TypeString`
- [x] 1.3 设置字段为 `Optional: true`（可选参数）
- [x] 1.4 设置 `ForceNew: true`（不支持修改，需重建）
- [x] 1.5 添加清晰的 Description 说明参数用途和可选值
- [x] 1.6 将字段放置在 `bandwidth_limit_type` 之后，`route_ecmp_flag` 之前（保持逻辑分组）

## 2. 服务层实现
- [x] 2.1 更新 `service_tencentcloud_ccn.go` 中的 `CcnBasicInfo` 结构体，添加 `instanceMeteringType string` 字段
- [x] 2.2 更新 `CreateCcn` 方法签名，添加 `instanceMeteringType` 参数
- [x] 2.3 在 `CreateCcn` 方法中将参数传递给 `request.InstanceMeteringType`
- [x] 2.4 更新 `DescribeCcns` 方法中的解析逻辑，从 SDK 响应中读取 `InstanceMeteringType` 字段
- [x] 2.5 处理 `InstanceMeteringType` 字段为 nil 的情况（老实例可能不返回此字段）

## 3. 资源 CRUD 实现
- [x] 3.1 在 `resourceTencentCloudCcnCreate` 中读取 `instance_metering_type` 参数
- [x] 3.2 将 `instanceMeteringType` 传递给 `service.CreateCcn` 方法
- [x] 3.3 在 `resourceTencentCloudCcnRead` 中设置 `instance_metering_type` 到状态（使用 `d.Set`）
- [x] 3.4 确认 Update 函数不处理该字段（因为 ForceNew: true）

## 4. 测试实现
- [x] 4.1 更新 `resource_tc_ccn_test.go` 中的测试用例
- [x] 4.2 在基础测试配置中添加 `instance_metering_type` 参数
- [x] 4.3 添加测试检查项验证 `instance_metering_type` 值正确设置
- [x] 4.4 确保现有测试用例仍然通过（向后兼容性）
- [ ] 4.5 运行 `make test` 验证单元测试
- [ ] 4.6 运行验收测试验证端到端功能

## 5. 代码质量检查
- [x] 5.1 运行 `make fmt` 格式化代码
- [x] 5.2 运行 `make lint` 确保无 lint 错误
- [x] 5.3 检查错误处理逻辑（nil 检查，日志记录）
- [x] 5.4 确保所有新增代码都有适当的注释
- [x] 5.5 验证日志记录包含新参数信息

## 6. 文档更新（可选）
- [ ] 6.1 更新资源文档中的参数列表（如果存在单独的 .md 文件）
- [ ] 6.2 在示例配置中添加 `instance_metering_type` 参数
- [ ] 6.3 说明参数的可选值和使用场景
- [ ] 6.4 运行 `make doc` 生成文档（如果需要）

## 7. 最终验证
- [ ] 7.1 手动测试创建带 `instance_metering_type` 参数的 CCN 实例
- [ ] 7.2 验证参数值正确传递到云平台
- [ ] 7.3 测试不指定参数时使用默认值的情况
- [ ] 7.4 测试资源导入功能（确保能正确读取 `instance_metering_type`）
- [ ] 7.5 测试修改 `instance_metering_type` 触发资源重建
- [ ] 7.6 确认所有测试通过，无回归问题
