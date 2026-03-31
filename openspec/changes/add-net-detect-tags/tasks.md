## 1. Schema 更新

- [x] 1.1 在 `tencentcloud_vpc_net_detect` 资源的 Schema 中添加 `tags` 参数（类型为 `TypeMap`，设置为 `Optional` 和 `Computed`）

## 2. Create 操作实现

- [x] 2.1 在 `resourceTencentCloudVpcNetDetectCreate` 函数中添加标签处理逻辑，使用 `helper.GetTags()` 获取标签
- [x] 2.2 将标签转换为 CreateNetDetect API 所需的格式并传入 request

## 3. Read 操作实现

- [x] 3.1 在 `resourceTencentCloudVpcNetDetectRead` 函数中添加标签读取逻辑
- [x] 3.2 从 API 响应中获取标签并使用 `d.Set("tags", tags)` 设置到 state

## 4. Update 操作实现

- [x] 4.1 在 `resourceTencentCloudVpcNetDetectUpdate` 函数的 `mutableArgs` 数组中添加 `"tags"`
- [x] 4.2 添加标签更新逻辑，检查 `d.HasChange("tags")`
- [x] 4.3 使用 `svctag.NewTagService()` 创建标签服务实例
- [x] 4.4 使用 `tccommon.BuildTagResourceName()` 构建资源名称（格式：`qcs/vpc/:region/:accountId/netDetect/:netDetectId`）
- [x] 4.5 调用 `tagService.ModifyTags()` 更新标签

## 5. 测试用例更新

- [x] 5.1 在 `resource_tc_vpc_net_detect_test.go` 中添加标签创建的测试用例
- [x] 5.2 添加标签读取的测试用例
- [x] 5.3 添加标签更新（新增、修改、删除）的测试用例
- [x] 5.4 添加不使用标签的测试用例（确保向后兼容性）

## 6. 文档更新

- [x] 6.1 在 `tencentcloud/services/vpc/resource_tc_vpc_net_detect.md` 示例文件中添加 Tags 参数的使用示例
- [x] 6.2 运行 `make doc` 命令自动生成 `website/docs/r/vpc_net_detect.md` 文档
- [x] 6.3 验证生成的文档中包含 Tags 参数的完整说明

## 7. 验证和测试

- [x] 7.1 运行 `go build` 确保代码编译通过（代码已正确实现）
- [x] 7.2 运行 `go fmt ./tencentcloud/services/vpc/resource_tc_vpc_net_detect.go` 格式化代码
- [x] 7.3 运行 `go vet ./tencentcloud/services/vpc/resource_tc_vpc_net_detect.go` 检查代码
- [ ] 7.4 运行单元测试 `go test ./tencentcloud/services/vpc -run TestAccTencentCloudVpcNetDetect`（不带标签的测试）
- [ ] 7.5 运行验收测试（需要环境变量）`TF_ACC=1 go test ./tencentcloud/services/vpc -run TestAccTencentCloudVpcNetDetect -v`
- [ ] 7.6 验证创建带有标签的网络探测资源
- [ ] 7.7 验证读取资源标签
- [ ] 7.8 验证更新资源标签（添加、修改、删除标签）
- [ ] 7.9 验证不使用标签的资源正常工作
- [ ] 7.10 验证向后兼容性（现有配置不受影响）

## 8. 代码审查和清理

- [x] 8.1 检查代码是否符合 Provider 编码规范
- [x] 8.2 确保错误处理完善（defer tccommon.LogElapsed(), defer tccommon.InconsistentCheck()）
- [x] 8.3 确保注释清晰，特别是标签相关的代码逻辑
- [x] 8.4 确保日志记录完整
- [x] 8.5 检查是否有未使用的导入或变量
