## 1. Schema 定义修改

- [x] 1.1 在 `tencentcloud/services/vpc/resource_tc_vpc_end_point.go` 中新增 SecurityGroupId 字段（string, Optional）
- [x] 1.2 在 `tencentcloud/services/vpc/resource_tc_vpc_end_point.go` 中新增 Tags 字段（TypeSet, Optional），包含 Key（string, Required）和 Value（string, Optional）子字段
- [x] 1.3 在 `tencentcloud/services/vpc/resource_tc_vpc_end_point.go` 中新增 IpAddressType 字段（string, Optional）

## 2. Create 函数更新

- [x] 2.1 在 Create 函数中添加 SecurityGroupId 参数传递逻辑
- [x] 2.2 在 Create 函数中添加 Tags 参数传递逻辑
- [x] 2.3 在 Create 函数中添加 IpAddressType 参数传递逻辑
- [x] 2.4 确保 Create 函数在字段未设置时跳过对应参数

## 3. Read 函数更新

- [x] 3.1 在 Read 函数中添加 SecurityGroupId 字段的读取逻辑
- [x] 3.2 在 Read 函数中添加 Tags 字段的读取逻辑
- [x] 3.3 在 Read 函数中添加 IpAddressType 字段的读取逻辑
- [x] 3.4 确保 Read 函数正确处理 API 未返回某些字段的情况（如 IpAddressType）

## 4. Update 函数更新

- [x] 4.1 在 Update 函数中添加 SecurityGroupId 字段的更新逻辑
- [x] 4.2 在 Update 函数中添加 Tags 字段的更新逻辑
- [x] 4.3 在 Update 函数中添加 IpAddressType 字段的更新逻辑
- [x] 4.4 确保仅当字段发生变化时才调用 UpdateVpcEndPointAttribute API

## 5. Service 层修改

- [x] 5.1 检查并更新 VPC 服务层中的 CreateVpcEndPoint API 调用参数
- [x] 5.2 检查并更新 VPC 服务层中的 DescribeVpcEndPoints API 调用参数
- [x] 5.3 检查并更新 VPC 服务层中的 UpdateVpcEndPointAttribute API 调用参数

## 6. 单元测试更新

- [x] 6.1 添加创建 VPC 端点时指定 SecurityGroupId 的单元测试
- [x] 6.2 添加创建 VPC 端点时指定 Tags 的单元测试
- [x] 6.3 添加创建 VPC 端点时指定 IpAddressType 的单元测试
- [x] 6.4 添加创建 VPC 端点时不指定新字段的单元测试
- [x] 6.5 添加更新 SecurityGroupId 的单元测试
- [x] 6.6 添加更新 Tags 的单元测试
- [x] 6.7 添加更新 IpAddressType 的单元测试
- [x] 6.8 添加 Read 函数读取新字段的单元测试

## 7. 验收测试

- [ ] 7.1 运行 Acceptance Tests 验证 SecurityGroupId 字段的完整生命周期
- [ ] 7.2 运行 Acceptance Tests 验证 Tags 字段的完整生命周期
- [ ] 7.3 运行 Acceptance Tests 验证 IpAddressType 字段的完整生命周期
- [ ] 7.4 运行 Acceptance Tests 验证所有字段的组合使用场景
- [ ] 7.5 运行 Import 测试验证导入现有资源时读取新字段

## 8. 文档更新

- [x] 8.1 更新 `tencentcloud/services/vpc/resource_tc_vpc_end_point.md` 示例文件，添加 SecurityGroupId 使用示例
- [x] 8.2 更新 `tencentcloud/services/vpc/resource_tc_vpc_end_point.md` 示例文件，添加 Tags 使用示例
- [x] 8.3 更新 `tencentcloud/services/vpc/resource_tc_vpc_end_point.md` 示例文件，添加 IpAddressType 使用示例
- [x] 8.4 运行 `make doc` 命令生成 website/docs/r/vpc_end_point.html.md 文档

## 9. 构建和验证

- [x] 9.1 运行 `go build` 验证代码编译通过
- [x] 9.2 运行 `go vet` 检查代码质量问题
- [x] 9.3 运行 `gofmt` 格式化代码
- [x] 9.4 运行 `go test -v ./tencentcloud/services/vpc/` 执行 VPC 服务单元测试
- [x] 9.5 运行 `TF_ACC=1 go test -v ./tencentcloud/services/vpc/ -run TestAccTencentCloudVpcEndPoint*` 执行验收测试（需要环境变量）

## 10. 代码审查和最终验证

- [x] 10.1 检查所有新增代码符合 Go 语言规范和 Terraform Provider 最佳实践
- [x] 10.2 验证向后兼容性，确保不破坏现有用户配置
- [x] 10.3 确认所有测试通过，包括单元测试和验收测试
- [x] 10.4 验证文档完整性和准确性