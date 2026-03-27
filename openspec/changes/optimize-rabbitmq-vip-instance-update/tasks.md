# 优化 RabbitMQ 实例 update 逻辑 - 实施任务

## 代码修改任务

### 任务 1: 优化 RabbitMQ VIP 实例的 Update 函数

**文件**: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`

**修改内容**:
1. 更新 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数的实现
2. 调整不可变字段列表，仅保留真正不可变的字段
3. 添加对不同 API 的调用逻辑：
   - `ModifyRabbitMQVipInstance` - 修改集群名称和标签
   - `ModifyRabbitMQVipInstanceSpec` - 修改节点规格、数量、存储
   - `ModifyRabbitMQVipInstancePublicAccess` - 修改公网访问和带宽
4. 为规格变更添加状态等待逻辑

**具体改动**:
- 将不可变字段列表从 12 个缩减到 7 个
- 移除 `node_spec`, `node_num`, `storage_size`, `enable_public_access`, `band_width` 的不可变限制
- 按字段变更类型分别调用对应的 API
- 添加 `ModifyRabbitMQVipInstanceSpec` 调用后的状态轮询逻辑

---

### 任务 2: 更新 Schema 字段的 ForceNew 标记

**文件**: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`

**修改内容**:
1. 检查并移除以下字段的 `ForceNew: true` 标记（如果存在）：
   - `node_spec`
   - `node_num`
   - `storage_size`
   - `enable_public_access`
   - `band_width`

**具体改动**:
- 在 Schema 定义中，从上述字段的配置中删除 `ForceNew: true`
- 确保这些字段现在是可更新的

---

### 任务 3: 更新资源示例文档

**文件**: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.md`

**修改内容**:
1. 在文档中明确说明哪些字段可以更新
2. 添加关于 Update 操作的说明
3. 添加注意事项和最佳实践

**具体改动**:
- 在参数说明中标记哪些字段支持更新
- 添加 Update 操作的示例
- 添加规格变更等待时间的说明

---

### 任务 4: 更新验收测试

**文件**: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`

**修改内容**:
1. 添加 Update 相关的测试用例
2. 验证可更新字段的正确性
3. 验证不可变字段的限制

**具体改动**:
- 添加测试用例 `TestAccTencentCloudTdmqRabbitmqVipInstance_updateNodeNum`
- 添加测试用例 `TestAccTencentCloudTdmqRabbitmqVipInstance_updateStorageSize`
- 添加测试用例 `TestAccTencentCloudTdmqRabbitmqVipInstance_updateNodeSpec`
- 添加测试用例 `TestAccTencentCloudTdmqRabbitmqVipInstance_updatePublicAccess`
- 添加测试用例 `TestAccTencentCloudTdmqRabbitmqVipInstance_updateBandwidth`
- 添加测试用例 `TestAccTencentCloudTdmqRabbitmqVipInstance_updateImmutableField`（验证不可变字段拒绝更新）

---

### 任务 5: 更新网站文档

**文件**: `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown`

**修改内容**:
1. 运行 `make doc` 命令自动生成文档
2. 验证生成的文档内容正确

**具体改动**:
- 执行 `make doc` 命令
- 检查生成的文档是否正确反映新的 Update 能力
- 确认文档格式正确

---

## 验证任务

### 任务 6: 编译检查

**命令**:
```bash
cd /repo
go build -o terraform-provider-tencentcloud
```

**验证内容**:
- 代码可以成功编译
- 没有语法错误
- 没有类型错误

---

### 任务 7: 静态代码检查

**命令**:
```bash
cd /repo
golangci-lint run
```

**验证内容**:
- 通过 golangci-lint 检查
- 没有代码风格问题
- 没有潜在的错误

---

### 任务 8: 运行单元测试

**命令**:
```bash
cd /repo
go test ./tencentcloud/services/trabbit -v -run TestAccTencentCloudTdmqRabbitmqVipInstance
```

**验证内容**:
- 所有测试用例通过
- 新增的 Update 测试用例通过
- 现有的测试用例不受影响

---

### 任务 9: 运行验收测试（需要真实环境）

**前置条件**:
- 设置 `TF_ACC=1` 环境变量
- 设置 `TENCENTCLOUD_SECRET_ID` 和 `TENCENTCLOUD_SECRET_KEY` 环境变量
- 有可用的腾讯云账户和 RabbitMQ 实例配额

**命令**:
```bash
cd /repo
TF_ACC=1 TENCENTCLOUD_SECRET_ID=xxx TENCENTCLOUD_SECRET_KEY=xxx \
  go test ./tencentcloud/services/trabbit -v \
  -run TestAccTencentCloudTdmqRabbitmqVipInstance_update -timeout 30m
```

**验证内容**:
- Update 测试用例在真实环境中通过
- API 调用正确
- 状态等待逻辑正确
- 不触发资源重建

---

### 任务 10: Terraform 计划验证（手动测试）

**前置条件**:
- 编译好的 provider
- 一个测试用的 Terraform 配置

**测试步骤**:
1. 创建一个 RabbitMQ 实例
2. 运行 `terraform apply`
3. 修改 `node_num` 字段
4. 运行 `terraform plan`，确认显示 "update" 而不是 "destroy + create"
5. 运行 `terraform apply`，确认更新成功
6. 验证实例状态和配置正确

**验证内容**:
- Terraform plan 显示正确的变更类型
- Update 操作成功执行
- 资源不会被销毁重建

---

### 任务 11: 文档生成验证

**命令**:
```bash
cd /repo
make doc
```

**验证内容**:
- 文档成功生成
- 生成的文档位于 `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown`
- 文档内容正确反映 Update 能力

---

## 任务执行顺序

1. 任务 1: 优化 Update 函数（核心改动）
2. 任务 2: 更新 Schema ForceNew 标记
3. 任务 3: 更新资源示例文档
4. 任务 4: 更新验收测试
5. 任务 6: 编译检查
6. 任务 7: 静态代码检查
7. 任务 8: 运行单元测试
8. 任务 9: 运行验收测试（可选，需要真实环境）
9. 任务 10: Terraform 计划验证（可选，需要手动测试）
10. 任务 5: 更新网站文档（通过 make doc）
11. 任务 11: 文档生成验证

---

## 注意事项

1. **API 兼容性**: 确保 `ModifyRabbitMQVipInstanceSpec` 和 `ModifyRabbitMQVipInstancePublicAccess` API 在当前 SDK 版本中可用
2. **状态等待**: 规格变更后的状态等待时间可能较长，确保超时时间设置合理
3. **错误处理**: 确保所有 API 调用都有适当的错误处理和重试逻辑
4. **向后兼容**: 确保修改不会破坏现有用户的配置
5. **测试覆盖**: 新增的 Update 功能需要有充分的测试覆盖
6. **文档同步**: 确保代码修改和文档修改同步完成
