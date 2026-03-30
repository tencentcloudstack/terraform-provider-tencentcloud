## 1. 代码修改

- [x] 1.1 在 resource_tc_teo_origin_group.go 的 resourceTencentCloudTeoOriginGroupCreate 函数中添加 HostHeader 参数读取和设置逻辑
- [x] 1.2 确认 HostHeader 参数处理逻辑与 Update 函数中的实现保持一致（参考第 449-451 行）

## 2. 代码验证

- [ ] 2.1 执行项目构建验证（make build）确保代码编译通过 (需要Go环境)
- [ ] 2.2 运行 tencentcloud_teo_origin_group 资源的单元测试（TF_ACC=1 go test ./tencentcloud/services/teo -v -run TestAccTencentCloudTeoOriginGroup）(需要Go环境和测试环境)
- [ ] 2.3 验证测试包含创建带 HostHeader 参数和不带 HostHeader 参数的场景

## 3. 文档更新

- [x] 3.1 检查 resource_tc_teo_origin_group.md 示例文件，确认是否需要添加 HostHeader 参数的使用示例
- [x] 3.2 如需要，更新示例文件以包含 HostHeader 参数的使用场景
- [ ] 3.3 执行 make doc 生成 website/docs/ 下的文档，确保文档正确生成 (需要make命令)

## 4. 最终验证

- [ ] 4.1 执行完整的 make 流程（make build && make test）确保所有检查通过 (需要make命令)
- [x] 4.2 确认代码变更符合 Terraform Provider 的编码规范和最佳实践
- [x] 4.3 验证向后兼容性，确保现有配置不受影响
