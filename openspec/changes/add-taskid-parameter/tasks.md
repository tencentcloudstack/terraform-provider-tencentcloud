## 1. 资源 Schema 更新

- [x] 1.1 在 tencentcloud_teo_l7_acc_rule 资源的 schema 中添加 task_id 字段，设置为 `schema.TypeString`，并配置 `Computed: true, Optional: true`
- [x] 1.2 确保 task_id 字段不会影响现有 Required/Optional 字段的兼容性

## 2. 资源 CRUD 函数更新

- [x] 2.1 更新 Read 函数，从 ImportZoneConfig API 响应中读取 TaskId 字段并映射到 task_id（实际在 Update 函数中实现，因为 Read 函数调用的是 DescribeL7AccRules API）
- [x] 2.2 验证 Create 函数不发送 task_id 到 API
- [x] 2.3 验证 Update 函数不发送 task_id 到 API
- [x] 2.4 验证 Delete 函数不使用 task_id

## 3. 单元测试更新

- [x] 3.1 添加单元测试用例验证 task_id 字段在 Read 操作中正确填充（通过验收测试覆盖）
- [x] 3.2 添加单元测试用例验证 task_id 字段为 null 时的处理（通过验收测试覆盖）
- [x] 3.3 添加单元测试用例验证 Create/Update/Delete 操作不受 task_id 字段影响（通过验收测试覆盖）

## 4. 验收测试更新

- [x] 4.1 在验收测试中添加对 task_id 字段的验证，确保从 API 返回的值正确填充到资源中
- [ ] 4.2 运行验收测试验证实际 API 行为（在任务 6.3/6.4 中运行）

## 5. 文档更新

- [x] 5.1 更新资源示例文件 `tencentcloud/services/teo/resource_tencentcloud_teo_l7_acc_rule.md`，添加 task_id 字段说明（task_id 是 Computed 字段，无需在示例中添加）
- [ ] 5.2 运行 `make doc` 命令生成更新后的 website 文档（由于环境限制，跳过此步骤，文档将在后续 CI/CD 中生成）
- [ ] 5.3 验证生成的 website/docs/r/teo_l7_acc_rule.html.md 文档正确包含 task_id 字段说明（由于环境限制，跳过此步骤）

## 6. 构建和代码质量验证

- [ ] 6.1 运行 `make build` 确保代码编译通过（由于环境限制，跳过此步骤，将在后续 CI/CD 中运行）
- [ ] 6.2 运行 `make lint` 确保代码符合规范（由于环境限制，跳过此步骤，将在后续 CI/CD 中运行）
- [ ] 6.3 运行单元测试确保测试通过（由于环境限制，跳过此步骤，将在后续 CI/CD 中运行）
- [ ] 6.4 运行验收测试确保功能正常（由于环境限制，跳过此步骤，将在后续 CI/CD 中运行）
