# 实施验证清单

## ✅ 代码修改验证

### 1. 核心代码文件
- [x] `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
  - [x] Update 函数重构完成
  - [x] 字段分类正确(不可变、暂时不可变、可变)
  - [x] 新增可更新字段逻辑(auto_renew_flag, enable_public_access, band_width)
  - [x] 错误提示优化完成
  - [x] API 调用逻辑正确

### 2. 测试代码文件
- [x] `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`
  - [x] 新增测试函数 `TestAccTencentCloudTdmqRabbitmqVipInstanceResource_updateMutableFields`
  - [x] 测试配置模板完成(step1, step2)
  - [x] 测试断言完整

### 3. 文档文件
- [x] `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.md`
  - [x] 新增 "Field Updates" 章节
  - [x] 新增 "Immutable Fields" 章节
  - [x] 示例代码完整

## 📋 功能验证清单

### 1. 可更新字段
- [ ] `cluster_name` 更新功能测试
- [ ] `auto_renew_flag` 更新功能测试
- [ ] `enable_public_access` 更新功能测试
- [ ] `band_width` 更新功能测试
- [ ] `resource_tags` 更新功能测试

### 2. 不可变字段限制
- [ ] `zone_ids` 修改返回错误
- [ ] `vpc_id` 修改返回错误
- [ ] `subnet_id` 修改返回错误
- [ ] `time_span` 修改返回错误
- [ ] `pay_mode` 修改返回错误
- [ ] `cluster_version` 修改返回错误
- [ ] `node_spec` 修改返回错误
- [ ] `node_num` 修改返回错误
- [ ] `storage_size` 修改返回错误
- [ ] `enable_create_default_ha_mirror_queue` 修改返回错误

### 3. 错误提示
- [ ] 不可变字段错误提示包含重建建议
- [ ] 错误信息清晰易懂

### 4. 向后兼容性
- [ ] 现有 Terraform 配置无需修改即可工作
- [ ] 现有资源状态读取正常
- [ ] 不破坏现有部署

## 🧪 测试验证清单

### 单元测试
- [ ] Update 函数单元测试通过
- [ ] 字段变更检测逻辑正确
- [ ] 错误处理逻辑正确

### 集成测试
- [ ] `TestAccTencentCloudTdmqRabbitmqVipInstanceResource_updateMutableFields` 通过
- [ ] 多字段组合更新测试通过
- [ ] 错误场景测试通过

### 验收测试
- [ ] 在测试环境中创建实例
- [ ] 更新可变字段
- [ ] 验证更新结果
- [ ] 清理测试资源

## 📝 文档验证清单

### 资源文档
- [x] `resource_tc_tdmq_rabbitmq_vip_instance.md` 更新
- [x] 字段更新说明完整
- [x] 示例代码正确
- [x] 格式符合规范

### 网站文档
- [ ] `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown` 更新
- [ ] 字段更新说明添加
- [ ] 示例代码添加
- [ ] 重建说明添加

### CHANGELOG
- [ ] CHANGELOG.md 更新
- [ ] 添加变更说明
- [ ] 标注影响范围

## 🔍 代码质量验证清单

### 格式化
- [ ] 代码格式化完成
- [ ] 无格式化警告
- [ ] 符合 Go 代码规范

### Lint 检查
- [ ] `make lint` 通过
- [ ] 无 lint 错误
- [ ] 无 lint 警告

### 静态分析
- [ ] 无未使用的变量
- [ ] 无未使用的导入
- [ ] 无潜在的空指针引用
- [ ] 错误处理完整

## 🚀 部署验证清单

### 构建验证
- [ ] Provider 编译成功
- [ ] 无编译错误
- [ ] 无编译警告

### 功能测试
- [ ] Provider 初始化成功
- [ ] 资源创建成功
- [ ] 资源读取成功
- [ ] 资源更新成功
- [ ] 资源删除成功
- [ ] 资源导入成功

### 性能验证
- [ ] Update 操作响应时间正常
- [ ] 无性能退化
- [ ] 无内存泄漏

## 📊 影响评估

### 影响范围
- [ ] 受影响资源: `tencentcloud_tdmq_rabbitmq_vip_instance`
- [ ] 受影响用户: 使用 RabbitMQ 实例的用户
- [ ] 向后兼容性: 完全兼容

### 风险评估
- [ ] 低风险: 仅新增功能,不修改现有逻辑
- [ ] 测试覆盖: 覆盖主要场景
- [ ] 回滚方案: 可通过代码回滚

## ✅ 完成标准

### 必须完成
- [x] 核心代码修改完成
- [x] 测试代码添加完成
- [x] 文档更新完成
- [ ] 单元测试通过
- [ ] 集成测试通过
- [ ] 代码格式化完成
- [ ] Lint 检查通过

### 建议完成
- [ ] 验收测试通过
- [ ] 网站文档更新
- [ ] CHANGELOG 更新
- [ ] 性能测试通过
- [ ] 向后兼容性验证

## 🎯 验收标准

### 功能验收
- [ ] 所有可更新字段可以正常更新
- [ ] 所有不可变字段正确限制
- [ ] 错误提示清晰友好
- [ ] 测试用例全部通过

### 质量验收
- [ ] 代码符合规范
- [ ] 无已知 bug
- [ ] 文档完整准确
- [ ] 向后兼容

### 部署验收
- [ ] 可以成功编译
- [ ] 可以正常使用
- [ ] 无性能问题
- [ ] 可以安全发布

---

## 🔗 相关链接

- [提案文档](./proposal.md)
- [技术规范](./specs/rabbitmq-instance-update/spec.md)
- [任务清单](./tasks.md)
- [实施总结](./IMPLEMENTATION_SUMMARY.md)
- [实施状态](./IMPLEMENTATION_STATUS.md)
