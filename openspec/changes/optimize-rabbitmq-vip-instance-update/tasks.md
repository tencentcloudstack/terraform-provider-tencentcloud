## 1. Schema 定义修改

- [x] 1.1 在 resource_tc_tdmq_rabbitmq_vip_instance.go 的 Schema 中添加 remark 字段（Optional: true, Computed: true, Type: String）
- [x] 1.2 在 resource_tc_tdmq_rabbitmq_vip_instance.go 的 Schema 中添加 enable_deletion_protection 字段（Optional: true, Computed: true, Type: Bool）
- [x] 1.3 在 resource_tc_tdmq_rabbitmq_vip_instance.go 的 Schema 中添加 enable_risk_warning 字段（Optional: true, Computed: true, Type: Bool）

## 2. Read 函数修改

- [x] 2.1 在 resourceTencentCloudTdmqRabbitmqVipInstanceRead 函数中从 rabbitmqVipInstance.ClusterInfo 读取 Remark 并设置到 state
- [x] 2.2 在 resourceTencentCloudTdmqRabbitmqVipInstanceRead 函数中从 rabbitmqVipInstance.ClusterInfo 读取 EnableDeletionProtection 并设置到 state
- [x] 2.3 在 resourceTencentCloudTdmqRabbitmqVipInstanceRead 函数中从 rabbitmqVipInstance.ClusterInfo 读取 EnableRiskWarning 并设置到 state
- [x] 2.4 在 resourceTencentCloudTdmqRabbitmqVipInstanceRead 函数的 DescribeTdmqRabbitmqVipInstanceByFilter 部分从 result[0] 读取 Remark 并设置到 state
- [x] 2.5 在 resourceTencentCloudTdmqRabbitmqVipInstanceRead 函数的 DescribeTdmqRabbitmqVipInstanceByFilter 部分从 result[0] 读取 EnableDeletionProtection 并设置到 state

## 3. Update 函数修改

- [x] 3.1 在 resourceTencentCloudTdmqRabbitmqVipInstanceUpdate 函数中添加 remark 字段的更新逻辑（检测 HasChange 并调用 Modify API）
- [x] 3.2 在 resourceTencentCloudTdmqRabbitmqVipInstanceUpdate 函数中添加 enable_deletion_protection 字段的更新逻辑（检测 HasChange 并调用 Modify API）
- [x] 3.3 在 resourceTencentCloudTdmqRabbitmqVipInstanceUpdate 函数中添加 enable_risk_warning 字段的更新逻辑（检测 HasChange 并调用 Modify API）
- [x] 3.4 优化 resource_tags 更新逻辑：当 resource_tags 为空时设置 RemoveAllTags 为 true，当 resource_tags 非空时设置 Tags 参数

## 4. 代码验证和测试

- [x] 4.1 编译代码验证：go build -o /dev/null ./tencentcloud/services/trabbit/
- [ ] 4.2 代码格式化：gofmt -w ./tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go（在收尾阶段执行）
- [ ] 4.3 运行单元测试：go test -v ./tencentcloud/services/trabbit/ -run TestAccTencentCloudTdmqRabbitmqVipInstance（在收尾阶段执行）

## 5. 文档更新

- [ ] 5.1 更新 resource_tc_tdmq_rabbitmq_vip_instance.go 中的字段描述（为新增字段添加 Description）（在收尾阶段执行）
- [ ] 5.2 运行 make doc 生成 website/docs/r/tdmq_rabbitmq_vip_instance.html.md（在收尾阶段执行）
- [ ] 5.3 验证生成的文档包含新增字段的说明（在收尾阶段执行）
